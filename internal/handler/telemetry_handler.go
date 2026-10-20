package handler

import (
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/service"
	"database/sql"
	"encoding/json"
	"net/http"
)

type TelemetryHandler struct {
	Service *service.TelemetryService
	Log     *logger.Logger
}

func NewTelemetryHandler(s *service.TelemetryService, log *logger.Logger) *TelemetryHandler {
	return &TelemetryHandler{
		Service: s,
		Log:     log,
	}
}

// POST /telemetry
func (h *TelemetryHandler) ReceiveTelemetry(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context(), h.Log)
	var telemetry models.TelemetryPayload

	err := json.NewDecoder(r.Body).Decode(&telemetry)
	if err != nil {
		log.Errorf("Failed to decode telemetry payload: %v", err)
		http.Error(w, "Invalid telemetry payload", http.StatusBadRequest)
		return
	}

	log.Debugf("Received telemetry for VIN: %s", telemetry.VIN)

	err = h.Service.ProcessTelemetry(r.Context(), telemetry)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}
		log.Errorf("Error processing telemetry: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Telemetry processed successfully for VIN: %s", telemetry.VIN)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Telemetry processed successfully"))
}
