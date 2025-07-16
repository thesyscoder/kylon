/**
 * @File: kubernetes.go
 * @Title: Kubernetes Client Management
 * @Description: Handles the initialization and retrieval of the Kubernetes Clientset,
 * @Description: supporting both in-cluster and kubeconfig-based configurations.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package kubernetes

import (
	"fmt"
	"sync" // For sync.Once to ensure single initialization

	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"
	"github.com/thesyscoder/kylon/pkg/logger" // Import the custom logger
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// log is the logger instance for this package.
var log = logger.GetLogger().WithField("component", "kubernetes_client")

// clientSet holds the singleton Kubernetes Clientset instance.
var clientSet *kubernetes.Clientset

// initOnce ensures that InitKubernetesClient is called only once.
var initOnce sync.Once

// InitKubernetesClient initializes the Kubernetes Clientset.
// It first attempts an in-cluster configuration. If that fails (e.g., not running inside a cluster),
// it falls back to building the configuration from a kubeconfig file.
// The kubeconfig path can be specified in the application config; otherwise, it defaults to the recommended home file.
// This function is designed to be called once during application startup.
func InitKubernetesClient(cfg config.Config) error {
	var initializationError error // Variable to capture error from the once.Do block

	initOnce.Do(func() {
		log.Info("Attempting to initialize Kubernetes client...")
		var restConfig *rest.Config
		var err error

		// Try to create in-cluster config first.
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			log.Warnf("Failed to create in-cluster Kubernetes config: %v. Attempting kubeconfig...", err)

			// Fallback to kubeconfig.
			kubeConfigPath := cfg.Kubernetes.KubeconfigPath
			if kubeConfigPath == "" {
				kubeConfigPath = clientcmd.RecommendedHomeFile
				log.Infof("Kubeconfig path not specified in config. Defaulting to '%s'.", kubeConfigPath)
			}

			restConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
			if err != nil {
				initializationError = customerrors.NewCustomError(
					customerrors.ErrCodeK8sClientInitFailed,
					fmt.Sprintf("Failed to create Kubernetes config from kubeconfig '%s'", kubeConfigPath),
					err,
				)
				log.Errorf("Kubernetes client initialization failed: %v", initializationError)
				return // Exit Do func
			}
		}

		// Create the Clientset from the determined rest.Config.
		clientSet, err = kubernetes.NewForConfig(restConfig)
		if err != nil {
			initializationError = customerrors.NewCustomError(
				customerrors.ErrCodeK8sClientInitFailed,
				"Failed to create Kubernetes clientset",
				err,
			)
			log.Errorf("Kubernetes client initialization failed: %v", initializationError)
			return // Exit Do func
		}

		log.Info("Kubernetes client initialized successfully.")
	})

	return initializationError // Return the error captured during initialization
}

// GetKubernetesClient returns the initialized Kubernetes Clientset instance.
// It returns an error if the client has not been successfully initialized yet.
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	if clientSet == nil {
		return nil, customerrors.NewCustomError(
			customerrors.ErrCodeK8sClientNotInitialized,
			"Kubernetes client has not been initialized. Call InitKubernetesClient first.",
			nil,
		)
	}
	return clientSet, nil
}
