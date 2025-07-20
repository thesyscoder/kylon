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
	"github.com/thesyscoder/kylon/internal/infrastructure/database"
	"github.com/thesyscoder/kylon/internal/infrastructure/kubernetes"
	"github.com/thesyscoder/kylon/pkg/logger"
	"gorm.io/gorm"
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

	var db *gorm.DB            // Declare db variable
	dbConnectErr := error(nil) // Declare an error variable for DB connection

	db, dbConnectErr = database.ConnectPostgres(cfg)
	if dbConnectErr != nil {
		log.Printf("[Main]: Failed to connect to database: %v. The health endpoint will report this.", dbConnectErr)
		// db will be nil here, which is handled by the HealthHandler.
	} else {
		log.Println("[Main]: Database connection established successfully.")

		// --- 3. Auto-Migrate Database schemas (ONLY if DB connection was successful) ---
		// If DB connection failed, we cannot run migrations, and this is typically a fatal issue
		// for the application's full functionality.
		migrateErr := db.AutoMigrate()
		if migrateErr != nil {
			log.Fatalf("[Database]: Failed to auto-migrate database: %v. Application cannot function without migrations.", migrateErr)
		}
		log.Println("[Database]: Database migrations completed successfully.")
	}

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

	// --- Step 4:  Kubernetes API Interaction (Optional) ---
	// This block demonstrates how to use the initialized Kubernetes client
	// to list nodes. It also serves as a basic connectivity test.
	kubeClient, err := kubernetes.GetKubernetesClient()
	if err != nil {
		// This error should ideally not happen if InitKubernetesClient succeeded,
		// but it's good practice to check.
		log.WithError(err).Error("Failed to retrieve Kubernetes client instance after successful initialization.")
		os.Exit(1) // Treat as fatal if client unexpectedly unavailable.
	}
	// --- Step 5: Application Startup (Placeholder) ---
	// In a real application, you would start your HTTP server, message queues,
	// background workers, etc., here.
	log.Infof("Starting Kylon backend in %s mode on port %s", cfg.App.Host, cfg.App.Port)
	// Instantiate our custom server, injecting dependencies (config, db).
	// 'db' might be nil if the connection failed, which the server and handlers must handle.
	appServer := NewServer(cfg, nil, kubeClient)

	// --- 5. Start Server ---
	// Start the HTTP server and block until a shutdown signal is received.
	// The server will now start even if the database connection failed, allowing the health endpoint to respond.
	appServer.Start()

}
