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
	"fmt"
	"os"

	"github.com/thesyscoder/kylon/internal/domain/models"
	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"github.com/thesyscoder/kylon/internal/infrastructure/database"
	"github.com/thesyscoder/kylon/internal/infrastructure/kubernetes"
	"github.com/thesyscoder/kylon/pkg/logger"
	"gorm.io/gorm"
	k8sClient "k8s.io/client-go/kubernetes" // Added k8s.io/client-go/kubernetes alias
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
	// Ensure logger.SetLogger exists and correctly configures logrus.
	logger.SetLogger(cfg.Log.Level)
	log := logger.GetLogger() // This 'log' is your configured *logrus.Logger
	log.Info("Kylon Backend Server: Starting up...")

	var db *gorm.DB
	// --- Step 3: Connect to Database ---
	// Database connection is attempted. The server will start even if it fails,
	// allowing the health endpoint to report the DB status.
	db, err = database.ConnectPostgres(cfg) // Use 'err' from the outer scope
	if err != nil {
		log.WithError(err).Warn("Failed to connect to database. The health endpoint will report this.")
		db = nil // Ensure db is explicitly nil if connection fails
	} else {
		log.Info("Database connection established successfully.")

		// --- Auto-Migrate Database schemas (ONLY if DB connection was successful) ---
		// If DB connection failed, migrations cannot run.
		migrateErr := db.AutoMigrate(&models.Cluster{})
		if migrateErr != nil {
			log.WithError(migrateErr).Fatal("Failed to auto-migrate database. Application cannot function without migrations.")
			// os.Exit(1) is handled by Fatal, which calls os.Exit(1) by default
		}
		log.Info("Database migrations completed successfully.")
	}

	// --- Step 4: Initialize Kubernetes Client ---
	// Initialize the Kubernetes client. This function handles both in-cluster and kubeconfig setups.
	log.Info("Initializing Kubernetes client...")
	initKubeErr := kubernetes.InitKubernetesClient(*cfg)
	if initKubeErr != nil {
		// Log the error using the configured logger before exiting.
		log.WithError(initKubeErr).Fatal("Failed to initialize Kubernetes client. Exiting.")
	}
	log.Info("Kubernetes client initialized successfully.")

	// --- Step 5: Retrieve Kubernetes Client Instance ---
	// Get the initialized Kubernetes client.
	var kubeClient *k8sClient.Clientset // Declare with the correct type alias
	kubeClient, err = kubernetes.GetKubernetesClient()
	if err != nil {
		// This error should ideally not happen if InitKubernetesClient succeeded,
		// but it's good practice to check.
		log.WithError(err).Fatal("Failed to retrieve Kubernetes client instance after successful initialization. Exiting.")
	}
	// --- Step 6: Application Server Startup ---
	log.Infof("Starting Kylon backend in %s mode on port %s", cfg.App.Env, cfg.App.Port) // Use App.Env for clarity

	// Instantiate our custom server, injecting dependencies.
	// Pass the 'db' variable, which will be nil if the connection failed,
	// or a valid *gorm.DB if successful.
	appServer := NewServer(cfg, db, kubeClient, log) // Pass 'log' (logrus) instance

	// --- Step 7: Start Server ---
	// Start the HTTP server and block until a shutdown signal is received.
	// The server will now start even if the database connection failed,
	// allowing the health endpoint to respond appropriately.
	appServer.Start()
}
