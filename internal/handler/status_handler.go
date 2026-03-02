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

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	encoder.Encode(data)

}

// GET /cars/{vin}/status
func (h *StatusHandler) GetCarStatus(w http.ResponseWriter, r *http.Request) {
	vin := r.PathValue("vin")
	attribute := r.URL.Query().Get("attribute")

	if attribute == "" {
		h.Log.Infof("Fetching full status for VIN: %s", vin)
	} else {
		h.Log.Infof("Fetching status for VIN: %s, attribute: %s", vin, attribute)
	}

	result, err := h.TelemetryService.GetCarStatus(vin, attribute)
	if err != nil {
		h.Log.Errorf("Error fetching car status: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// if no telemetry -> return empty json
	if result == nil {
		h.Log.Infof("No telemetry yet for VIN: %s", vin)
		writeJSON(w, http.StatusOK, map[string]interface{}{})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GET /cars/{vin}/events?date=YYYY-MM-DD
func (h *StatusHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	vin := r.PathValue("vin")
	if vin == "" {
		h.Log.Warn("GetEvents called without VIN")
		http.Error(w, "vin is required", http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		h.Log.Infof("Fetching events for VIN: %s", vin)
	} else {
		h.Log.Infof("Fetching events for VIN: %s, date: %s", vin, date)
	}

	result, err := h.TelemetryService.GetEvents(vin, date)
	if err != nil {
		h.Log.Errorf("Error fetching events: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
