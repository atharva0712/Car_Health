package main

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/db"
	"car-monitoring-system/internal/handler"
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/repository"
	"car-monitoring-system/internal/service"
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
	}

	// Database connection
	database := db.Connect()
	log.Info("Database connected successfully")

	// Car layer
	carRepo := repository.NewCarRepository(database)
	carService := service.NewCarService(carRepo, log)
	carHandler := handler.NewCarHandler(carService, log)

	// Telemetry layer
	telemetryRepo := repository.NewTelemetryRepository(database)
	telemetryService := service.NewTelemetryService(telemetryRepo, cfg, log)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService, log)

	// Status layer
	statusHandler := handler.NewStatusHandler(telemetryService, log)

	// Routes
	mux := http.NewServeMux()

	// Car endpoints
	mux.HandleFunc("POST /cars", carHandler.RegisterCar)
	mux.HandleFunc("GET /cars", carHandler.GetCars)
	mux.HandleFunc("GET /cars/{vin}", carHandler.GetCarByVIN)
	mux.HandleFunc("DELETE /cars/{vin}", carHandler.DeleteCar)
	mux.HandleFunc("PATCH /cars/{vin}/simulation", carHandler.UpdateSimulationStatus)

	// Telemetry endpoint
	mux.HandleFunc("POST /telemetry", telemetryHandler.ReceiveTelemetry)

	// Status endpoint
	mux.HandleFunc("GET /cars/{vin}/status", statusHandler.GetCarStatus)
	mux.HandleFunc("GET /cars/{vin}/events", statusHandler.GetEvents)

	// Get port from env
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8000"
	}

	log.Infof("Backend server running on port %s", port)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Errorf("Server error: %v", err)
	}
}
