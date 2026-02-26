package handler

import (
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type TelemetryHandler struct {
	Service *service.TelemetryService
}

func NewTelemetryHandler(s *service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{Service: s}
}

// POST /telemetry
func (h *TelemetryHandler) ReceiveTelemetry(w http.ResponseWriter, r *http.Request) {
	var telemetry models.TelemetryPayload

	err := json.NewDecoder(r.Body).Decode(&telemetry)
	if err != nil {
		http.Error(w, "Invalid telemetry payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received telemetry for VIN: %s", telemetry.VIN)

	err = h.Service.ProcessTelemetry(telemetry)
	if err != nil {
		log.Printf("Error processing telemetry: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Telemetry processed successfully"))
}
