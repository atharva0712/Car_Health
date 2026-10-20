package handler

import (
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/service"
	"database/sql"
	"encoding/json"
	"net/http"
)

type CarHandler struct {
	Service *service.CarService
	Log     *logger.Logger
}

// NewCarHandler creates a new CarHandler instance
func NewCarHandler(service *service.CarService, log *logger.Logger) *CarHandler {
	return &CarHandler{
		Service: service,
		Log:     log,
	}
}

// POST /cars

func (h *CarHandler) RegisterCar(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context(), h.Log)
	var req struct {
		VIN string `json:"vin"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("Failed to decode request: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	log.Infof("Registering car with VIN: %s", req.VIN)
	err = h.Service.RegisterCar(r.Context(), req.VIN)
	if err != nil {
		log.Errorf("Failed to register car: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Infof("Car registered successfully: %s", req.VIN)
	w.Write([]byte("Car Registered"))
}

// GET /cars

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context(), h.Log)
	log.Debug("Fetching all cars")
	cars, err := h.Service.GetCars(r.Context())
	if err != nil {
		log.Errorf("Error fetching cars: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("Successfully fetched %d cars", len(cars))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

// GET /cars/{vin}
func (h *CarHandler) GetCarByVIN(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context(), h.Log)
	vin := r.PathValue("vin")
	log.Infof("Fetching car with VIN: %s", vin)

	car, err := h.Service.GetCar(r.Context(), vin)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("Car not found: %s", vin)
			http.Error(w, "Car not found", http.StatusNotFound)
		} else {
			log.Errorf("Database error fetching car %s: %v", vin, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

// DELETE /cars/{vin}
func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context(), h.Log)
	vin := r.PathValue("vin")
	log.Infof("Deleting car with VIN: %s", vin)

	err := h.Service.DeleteCar(r.Context(), vin)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Car not found", http.StatusNotFound)
		} else {
			log.Errorf("Failed to delete car: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	log.Infof("Car deleted successfully: %s", vin)
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Car Deleted. VIN: " + vin))
}

// PATCH /cars/{vin}/simulation

func (h *CarHandler) UpdateSimulationStatus(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context(), h.Log)
	vin := r.PathValue("vin")
	var req struct {
		SimulationEnabled bool `json:"simulation_enabled"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("Failed to decode simulation update request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Infof("Updating simulation status for VIN %s to %v", vin, req.SimulationEnabled)
	err = h.Service.UpdateSimulationStatus(r.Context(), vin, req.SimulationEnabled)

	if err == sql.ErrNoRows {
		log.Warnf("Car not found for update: %s", vin)
		http.Error(w, "Car Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Errorf("Failed to update simulation status: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("Simulation status updated successfully for VIN: %s", vin)
	w.Write([]byte("Simulation status updated successfully"))
}
