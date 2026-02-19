package handler

import (
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/service"
	"encoding/json"
	"net/http"
)

type TelemetryHandler struct {
	Service *service.TelemetryService
}

func NewTelemetryHandler(s *service.TelemetryService) *TelemetryHandler{
	return &TelemetryHandler{Service: s}
}

// Post /telemetry

func (h *TelemetryHandler) ReceiverTelemetry(w http.ResponseWriter, r *http.Request){
	var telemetry models.Telemetry
	err := json.NewDecoder(r.Body).Decode(&telemetry)
	if err != nil{
		http.Error(w,err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Service.ProcessTelemetry(telemetry)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Telemetry processed successfully"))
}