/**
 * @File: server.go
 * @Title: Server Initialization and Management
 * @Description: Handles the setup and lifecycle of the HTTP server, including graceful shutdown.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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
}

// NewServer creates and returns a new Server instance.
func NewServer(cfg *config.Config, db *gorm.DB, kubeClient *kubernetes.Clientset) *Server {
	gin.SetMode(gin.ReleaseMode)

	// SetupRouter returns a configured *gin.Engine
	router := routes.InitializeRoutes(cfg, db, kubeClient)

	return &Server{
		Gin: router,
		Cfg: cfg,
		DB:  db,
	}
}

// Start runs the HTTP server and handles graceful shutdown.
func (s *Server) Start() {
	addr := net.JoinHostPort(s.Cfg.App.Host, s.Cfg.App.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: s.Gin,
	}

	// Run server in goroutine
	go func() {
		log.Printf("[Server]: Starting HTTP server at %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server]: ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Server]: Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[Server]: Graceful shutdown failed: %v", err)
	}

	log.Println("[Server]: Server exited cleanly.")
}
