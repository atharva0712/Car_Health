package main

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/generator"
	"car-monitoring-system/internal/logger"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting telemetry generator...")
	godotenv.Load()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Errorf("Error loading config: %v", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Error("API_KEY is required for the telemetry generator")
		os.Exit(1)
	}

	ticker := generator.NewTicker(cfg, log, apiKey)
	ticker.Start()
}
