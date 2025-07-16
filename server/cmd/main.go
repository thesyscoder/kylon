package main

import (
	"github.com/thesyscoder/kylon/pkg/logger"
)

func main() {
	// Initialize the logger
	logger.SetLogger("debug")
	log := logger.GetLogger()
	log.Info("[Main]: Kylon Backend Server: Starting up...")
}
