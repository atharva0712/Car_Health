package main

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/db"
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/server"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize logger
	log := logger.NewLogger()
	log.Infof("Starting backend server...")

	// Load config.yaml
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Errorf("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Database connection
	database := db.Connect()
	log.Info("Database connected successfully")

	// Get port from env
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	log.Infof("Backend server running on port %s", port)

	err = http.ListenAndServe(":"+port, server.NewHandler(database, cfg, log))
	if err != nil {
		log.Errorf("Server error: %v", err)
	}
}
