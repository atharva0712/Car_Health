package main

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/db"
	"car-monitoring-system/internal/handler"
	"car-monitoring-system/internal/repository"
	"car-monitoring-system/internal/service"
	"log"
	"net/http"
	"os"
	
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()
	
	// Load config.yaml
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	
	// Database connection
	database := db.Connect()
	
	// Car layer
	carRepo := repository.NewCarRepository(database)
	carService := service.NewCarService(carRepo)
	carHandler := handler.CarHandler{Service: carService}
	
	// Telemetry layer
	telemetryRepo := repository.NewTelemetryRepository(database)
	telemetryService := service.NewTelemetryService(telemetryRepo, cfg)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService)
	
	// Status layer
	statusHandler := handler.NewStatusHandler(telemetryService)


	// Routes
	mux := http.NewServeMux()
	
	// Car endpoints
	mux.HandleFunc("POST /cars", carHandler.RegisterCar)
	mux.HandleFunc("GET /cars", carHandler.GetCars)
	mux.HandleFunc("GET /cars/{vin}", carHandler.GetCarByVIN)
	mux.HandleFunc("DELETE /cars/{vin}", carHandler.DeleteCar)
	
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
	
	log.Printf("Backend server running on port %s", port)
	
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}