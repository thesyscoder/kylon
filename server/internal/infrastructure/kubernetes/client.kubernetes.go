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
	"net/http" // Import for HTTP status codes
	"sync"     // For sync.Once to ensure single initialization

	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"
	"github.com/thesyscoder/kylon/pkg/logger" // Import the custom logger
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// log is the logger instance for this package, providing contextual logging for Kubernetes client operations.
var log = logger.GetLogger().WithField("component", "kubernetes_client")

// clientSet holds the singleton Kubernetes Clientset instance after successful initialization.
var clientSet *kubernetes.Clientset

// initOnce ensures that InitKubernetesClient is called only once across the application's lifetime.
var initOnce sync.Once

// InitKubernetesClient initializes the Kubernetes Clientset.
// It prioritizes an in-cluster configuration (for deployments within a Kubernetes cluster).
// If in-cluster configuration fails (e.g., when running locally outside a cluster),
// it falls back to building the configuration from a kubeconfig file.
// The kubeconfig path can be specified in the application's configuration; otherwise,
// it defaults to the standard kubeconfig location (e.g., ~/.kube/config).
// This function is designed to be called exactly once during application startup.
// It returns a `customerrors.CustomError` if initialization fails at any step.
func InitKubernetesClient(cfg config.Config) error {
	var initializationError error // Variable to capture any error that occurs during the `initOnce.Do` block.

	initOnce.Do(func() {
		log.Info("Attempting to initialize Kubernetes client...")
		var restConfig *rest.Config
		var err error

		// First, attempt to create an in-cluster configuration. This is standard for pods running inside K8s.
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			log.Warnf("Failed to create in-cluster Kubernetes config: %v. Attempting kubeconfig fallback...", err)

			// If in-cluster config fails, fall back to building from a kubeconfig file.
			kubeConfigPath := cfg.Kubernetes.KubeconfigPath
			if kubeConfigPath == "" {
				// If no kubeconfig path is provided in the application config, use the recommended default.
				kubeConfigPath = clientcmd.RecommendedHomeFile
				log.Infof("Kubeconfig path not specified in config. Defaulting to '%s'.", kubeConfigPath)
			}

			// Build config from flags (empty master URL, specified kubeconfig path).
			restConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
			if err != nil {
				// Capture and log the specific error if kubeconfig-based config creation fails.
				initializationError = customerrors.NewCustomError(
					customerrors.ErrCodeK8sClientInitFailed,
					fmt.Sprintf("Failed to create Kubernetes config from kubeconfig '%s'", kubeConfigPath),
					err,
					http.StatusInternalServerError, // Appropriate HTTP status for this internal setup error
					nil,                            // No additional data for this error
				)
				log.WithError(initializationError).Error("Kubernetes client initialization failed.")
				return // Exit the `once.Do` function, preventing further execution within this block.
			}
		}

		// With a valid rest.Config, create the Kubernetes Clientset.
		clientSet, err = kubernetes.NewForConfig(restConfig)
		if err != nil {
			// Capture and log the error if Clientset creation fails.
			initializationError = customerrors.NewCustomError(
				customerrors.ErrCodeK8sClientInitFailed,
				"Failed to create Kubernetes clientset",
				err,
				http.StatusInternalServerError, // Appropriate HTTP status for this internal setup error
				nil,                            // No additional data
			)
			log.WithError(initializationError).Error("Kubernetes client initialization failed.")
			return // Exit the `once.Do` function.
		}

		log.Info("Kubernetes client initialized successfully.")
	})

	// Return the error captured during the initialization process (if any).
	return initializationError
}

// GetKubernetesClient returns the initialized Kubernetes Clientset instance.
// It should be called after `InitKubernetesClient` has completed successfully.
// Returns a `customerrors.CustomError` if the client has not been initialized.
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	if clientSet == nil {
		// Return a specific error if the client was not initialized.
		return nil, customerrors.NewCustomError(
			customerrors.ErrCodeK8sClientNotInitialized,
			"Kubernetes client has not been initialized. Ensure InitKubernetesClient was called successfully during startup.",
			nil,                            // No underlying error
			http.StatusInternalServerError, // Appropriate HTTP status for this state
			nil,                            // No additional data
		)
	}
	return clientSet, nil
}
