package main

import (
	"log"

	"github.com/thesyscoder/kylon/internal/infrastructure/config"
)

func main() {

	log.Println("[Main]: Starting the kylon backend service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		// If config loading fails, the application cannot proceed. This is still a fatal error.
		log.Fatalf("[Main]: Failed to load configuration: %v", err)
	}
	log.Printf("[Main]: Configuration loaded successfully. App Environment: %s, Server will run on %s:%s.",
		cfg.App.Env, cfg.App.Host, cfg.App.Port)

	appServer := NewServer(cfg, nil)

	// --- 5. Start Server ---
	// Start the HTTP server and block until a shutdown signal is received.
	// The server will now start even if the database connection failed, allowing the health endpoint to respond.
	appServer.Start()

	log.Println("[Main]: Application stopped.")
}
