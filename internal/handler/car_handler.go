package handler

import (
	"car-monitoring-system/internal/service"
	"encoding/json"
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

	if err != nil{
		http.Error(w,err.Error(), 500)
		return
	}

	w.Write([]byte("Car Registered"))
}

// GET /cars

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request){
	cars, _ := h.Service.GetCars()
	json.NewEncoder(w).Encode(cars)
}



