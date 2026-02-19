package main

import (
	"car-monitoring-system/internal/db"
	"car-monitoring-system/internal/handler"
	"car-monitoring-system/internal/repository"
	"car-monitoring-system/internal/service"
	"log"
	"net/http"
)

func main() {
	database := db.Connect()

	carRepo := repository.NewCarRepository(database)

	carService := service.NewCarService(carRepo)

	carHandler := handler.CarHandler{Service: carService}

	telemetryRepo := repository.NewTelemetryRepository(database)

	telemetryService := service.NewTelemetryService(telemetryRepo)

	telemetryHandler := handler.NewTelemetryHandler(telemetryService)

	// Routes

	mux := http.NewServeMux()

	mux.HandleFunc("POST /cars", carHandler.RegisterCar)

	mux.HandleFunc("GET /cars", carHandler.GetCars)

	mux.HandleFunc("POST /telemetry", telemetryHandler.ReceiverTelemetry)

	log.Println("Server running on 8000")

	err := http.ListenAndServe(":8000", mux)

	if err != nil {
		log.Fatal(err)
	}

}
