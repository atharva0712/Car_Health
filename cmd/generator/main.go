package main

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/generator"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	ticker := generator.NewTicker(cfg)
	ticker.Start()
}