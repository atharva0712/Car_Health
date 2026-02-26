package handler

import (
	"encoding/json"
	"net/http"
	"car-monitoring-system/internal/service"
)

type StatusHandler struct {
	TelemetryService *service.TelemetryService
}

func NewStatusHandler(s *service.TelemetryService) *StatusHandler {
	return &StatusHandler{
		TelemetryService: s,
	}
}

// GET /cars/{vin}/status
func (h *StatusHandler) GetCarStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	vin := r.PathValue("vin")
	attribute := r.URL.Query().Get("attribute")

	result, err := h.TelemetryService.GetCarStatus(vin, attribute)
	if err != nil {
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
		http.Error(w, "vin is required", http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")

	result, err := h.TelemetryService.GetEvents(vin, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
	