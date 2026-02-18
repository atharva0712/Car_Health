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

	repo := repository.NewCarRepository(database)

	service := service.NewCarService(repo)

	handler := handler.CarHandler{Service: service}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /cars",handler.RegisterCar)

	mux.HandleFunc("GET /cars",handler.GetCars)

	log.Println("Server running on 8000")

	http.ListenAndServe(":8000",mux)

}