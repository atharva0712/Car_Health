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

// GET /cars/{vin}
func (h *CarHandler) GetCarByVIN(w http.ResponseWriter, r *http.Request) {
	vin := r.PathValue("vin")

	car, err := h.Service.GetCar(vin) // Uses GetCar from service
	if err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

// DELETE /cars/{vin}
func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	vin := r.PathValue("vin")

	err := h.Service.DeleteCar(vin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
