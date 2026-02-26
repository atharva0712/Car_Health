package handler

import (
	"car-monitoring-system/internal/service"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type CarHandler struct {
	Service *service.CarService
}

// POST /cars

func (h *CarHandler) RegisterCar(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VIN string `json:"vin"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = h.Service.RegisterCar(req.VIN)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("Car Registered"))
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
	w.Write([]byte("Car Deleted. VIN: " + vin))
}


// PATCH /cars/{vin}/simulation

func (h *CarHandler) UpdateSimulationStatus(w http.ResponseWriter, r *http.Request) {
	vin := r.PathValue("vin")
	var req struct {
		SimulationEnabled bool `json:"simulation_enabled"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w,err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateSimulationStatus(vin, req.SimulationEnabled)

	if err == sql.ErrNoRows{
		http.Error(w, "Car Not Found", http.StatusNotFound)
		return
	}

	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Simulation status updated successfully"))
}