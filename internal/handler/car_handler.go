package handler

import (
	"car-monitoring-system/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type CarHandler struct {
	Service *service.CarService
}

// POST /cars
func (h *CarHandler) RegisterCar(w http.ResponseWriter, r *http.Request) {
	var car struct {
		VIN string `json:"vin"`
	}
	json.NewDecoder(r.Body).Decode(&car)

	err := h.Service.RegisterCar(car.VIN)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("Car Registered. VIN: " + car.VIN))
}

// GET /cars

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	cars, err := h.Service.GetCars()
	if err != nil {
		log.Println("Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Cars:", cars)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}
