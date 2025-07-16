package main

import (
	"fmt"
	"os"

	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"github.com/thesyscoder/kylon/pkg/logger"
)

func main() {
	// Initialize the logger
	logger.SetLogger("debug")
	log := logger.GetLogger()
	log.Info("[Main]: Kylon Backend Server: Starting up...")

	// Load the configuration
	var err error
	cfg, err := config.LoadConfig()
	if err != nil {
		// Log to stderr directly before logger is fully configured
		fmt.Fprintf(os.Stderr, "Fatal: Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	log.Debug(cfg)
}
