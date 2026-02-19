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
