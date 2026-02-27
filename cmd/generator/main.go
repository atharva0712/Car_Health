package main

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/generator"
	"car-monitoring-system/internal/logger"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting telemetry generator...")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Errorf("Error loading config: %v", err)
	}

	ticker := generator.NewTicker(cfg, log)
	ticker.Start()
}