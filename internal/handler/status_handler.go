package handler

import (
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/service"
	"encoding/json"
	"net/http"
)

type StatusHandler struct {
	TelemetryService *service.TelemetryService
	Log              *logger.Logger
}

func NewStatusHandler(s *service.TelemetryService, log *logger.Logger) *StatusHandler {
	return &StatusHandler{
		TelemetryService: s,
		Log:              log,
	}
}

// GET /cars/{vin}/status
func (h *StatusHandler) GetCarStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	vin := r.PathValue("vin")
	attribute := r.URL.Query().Get("attribute")
	h.Log.Infof("Fetching status for VIN: %s, attribute: %s", vin, attribute)

	result, err := h.TelemetryService.GetCarStatus(vin, attribute)
	if err != nil {
		h.Log.Errorf("Error fetching car status: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET /cars/{vin}/events?date=YYYY-MM-DD
func (h *StatusHandler) GetEvents(
	w http.ResponseWriter,
	r *http.Request,
) {
	vin := r.PathValue("vin")
	if vin == "" {
		h.Log.Warn("GetEvents called without VIN")
		http.Error(w, "vin is required", http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")
	h.Log.Infof("Fetching events for VIN: %s, date: %s", vin, date)

	result, err := h.TelemetryService.GetEvents(vin, date)
	if err != nil {
		h.Log.Errorf("Error fetching events: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.Log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
