package handler

import (
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/service"
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
	var telemetry models.TelemetryPayload

	err := json.NewDecoder(r.Body).Decode(&telemetry)
	if err != nil {
		h.Log.Errorf("Failed to decode telemetry payload: %v", err)
		http.Error(w, "Invalid telemetry payload", http.StatusBadRequest)
		return
	}

	h.Log.Debugf("Received telemetry for VIN: %s", telemetry.VIN)

	err = h.Service.ProcessTelemetry(telemetry)
	if err != nil {
		h.Log.Errorf("Error processing telemetry: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Log.Infof("Telemetry processed successfully for VIN: %s", telemetry.VIN)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Telemetry processed successfully"))
}
