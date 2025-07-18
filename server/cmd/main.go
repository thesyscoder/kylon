/**
 * @File: main.go
 * @Title: Kylon Backend Server Entry Point
 * @Description: This is the main entry point for the Kylon backend application.
 * @Description: It handles application initialization, including configuration loading,
 * @Description: logger setup, and Kubernetes client initialization.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package main

import (
	// For context.Background() when making API calls
	"fmt"
	"os"

	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"github.com/thesyscoder/kylon/internal/infrastructure/kubernetes"
	"github.com/thesyscoder/kylon/pkg/logger"
	// For metav1.ListOptions{}
)

func main() {
	// --- Step 1: Load Configuration ---
	// Configuration is loaded first. If this fails, we cannot proceed,
	// so log to stderr directly before the custom logger is fully set up.
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// --- Step 2: Initialize Logger based on Config ---
	// The custom logger's level is set based on the loaded configuration.
	logger.SetLogger(cfg.Log.Level)
	log := logger.GetLogger()
	log.Info("Kylon Backend Server: Starting up...")

	// --- Step 3: Initialize Kubernetes Client ---
	// Initialize the Kubernetes client. This function handles both in-cluster and kubeconfig setups.
	log.Info("Initializing Kubernetes client...")
	err = kubernetes.InitKubernetesClient(*cfg)
	if err != nil {
		// Log the error using the configured logger before exiting.
		log.WithError(err).Error("Failed to initialize Kubernetes client. Exiting.")
		os.Exit(1)
	}
	log.Info("Kubernetes client initialized successfully.")

	// --- Step 4: Example Kubernetes API Interaction (Optional) ---
	// This block demonstrates how to use the initialized Kubernetes client
	// to list nodes. It also serves as a basic connectivity test.

	// --- Step 5: Application Startup (Placeholder) ---
	// In a real application, you would start your HTTP server, message queues,
	// background workers, etc., here.
	log.Info("Kylon Backend Server: Application is ready. (Placeholder for server startup)")

	// Keep the main goroutine alive, e.g., by starting an HTTP server
	// For demonstration, we'll just print a message and exit.
	// In a real application, this would be replaced by server.ListenAndServe() or similar.
	// select {} // Uncomment to keep the application running indefinitely
}
