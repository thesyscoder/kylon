// main.go (or server.go)
package main

import (
	"context" // Added for net.JoinHostPort
	// Standard log for Fatalf in main, distinct from logrus
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thesyscoder/kylon/internal/app/routes"
	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
)

// Server encapsulates the Gin engine and application configuration.
type Server struct {
	Gin *gin.Engine
	Cfg *config.Config
	DB  *gorm.DB
	Log *logrus.Logger // Renamed from 'log' to 'Log' to avoid shadowing
}

// NewServer creates and returns a new Server instance.
func NewServer(cfg *config.Config, db *gorm.DB, kubeClient *kubernetes.Clientset, appLogger *logrus.Logger) *Server { // Renamed param to appLogger
	// Set Gin mode (ReleaseMode is good for production)
	gin.SetMode(gin.ReleaseMode)

	// Initialize routes with all necessary dependencies
	// Pass kubeClient to InitializeRoutes
	router := routes.InitializeRoutes(cfg, db, appLogger, kubeClient)

	return &Server{
		Gin: router,
		Cfg: cfg,
		DB:  db,
		Log: appLogger, // Use the passed appLogger
	}
}

// Start runs the HTTP server and handles graceful shutdown.
func (s *Server) Start() {
	// Construct the address for the server
	addr := net.JoinHostPort(s.Cfg.App.Host, s.Cfg.App.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: s.Gin,
	}

	// Run server in a goroutine
	go func() {
		s.Log.Infof("Starting HTTP server at %s", addr) // Use s.Log for structured logging
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Use Fatalf from s.Log for unrecoverable errors after server starts
			s.Log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	// Listen for interrupt (Ctrl+C) and terminate signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received

	s.Log.Info("Shutdown signal received, shutting down gracefully...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5 seconds timeout
	defer cancel()                                                          // Ensure cancel is called to release context resources

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		s.Log.Fatalf("Graceful shutdown failed: %v", err) // Log fatal if shutdown fails
	}

	s.Log.Info("Server exited cleanly.")
}
