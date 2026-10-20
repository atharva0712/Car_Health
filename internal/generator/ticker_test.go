package generator

import (
	"car-monitoring-system/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTickerSendsAPIKeyOnBothRequests(t *testing.T) {
	const apiKey = "test-tenant-key"
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-API-Key"); got != apiKey {
			t.Errorf("X-API-Key = %q, want %q", got, apiKey)
		}
		requests++
		switch r.URL.Path {
		case "/cars":
			_ = json.NewEncoder(w).Encode([]models.Car{{VIN: "VIN-1", SimulationEnabled: true}})
		case "/telemetry":
			if r.Method != http.MethodPost {
				t.Errorf("telemetry method = %s, want POST", r.Method)
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ticker := &Ticker{
		backendURL: server.URL,
		apiKey:     apiKey,
		httpClient: server.Client(),
	}
	vins, err := ticker.fetchVINs()
	if err != nil {
		t.Fatal(err)
	}
	if len(vins) != 1 || vins[0] != "VIN-1" {
		t.Fatalf("fetchVINs returned %v", vins)
	}
	if err := ticker.sendTelemetry(models.TelemetryPayload{VIN: "VIN-1"}); err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("server received %d requests, want 2", requests)
	}
}
