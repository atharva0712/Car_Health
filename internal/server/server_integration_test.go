package server

import (
	"bytes"
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/repository"
	"car-monitoring-system/internal/tenant"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestTenantIsolationHTTP(t *testing.T) {
	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		t.Skip("set MYSQL_TEST_DSN to run MySQL integration tests")
	}
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	database.SetConnMaxLifetime(time.Minute)
	if err := database.Ping(); err != nil {
		t.Fatalf("connect to MYSQL_TEST_DSN: %v", err)
	}

	ctx := context.Background()
	tenantRepo := repository.NewTenantRepository(database)
	keyA, err := tenant.GenerateAPIKey()
	if err != nil {
		t.Fatal(err)
	}
	keyB, err := tenant.GenerateAPIKey()
	if err != nil {
		t.Fatal(err)
	}
	tenantA, err := tenantRepo.Create(ctx, "integration-a-"+keyA[:8], tenant.HashAPIKey(keyA))
	if err != nil {
		t.Fatal(err)
	}
	tenantB, err := tenantRepo.Create(ctx, "integration-b-"+keyB[:8], tenant.HashAPIKey(keyB))
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadConfig(filepath.Join("..", "..", "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	api := httptest.NewServer(NewHandler(database, cfg, logger.NewLogger()))
	defer api.Close()

	doRequest := func(method, path, key string, body any) (*http.Response, []byte) {
		t.Helper()
		var requestBody io.Reader
		if body != nil {
			encoded, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			requestBody = bytes.NewReader(encoded)
		}
		request, err := http.NewRequest(method, api.URL+path, requestBody)
		if err != nil {
			t.Fatal(err)
		}
		if key != "" {
			request.Header.Set("X-API-Key", key)
		}
		if body != nil {
			request.Header.Set("Content-Type", "application/json")
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		return response, responseBody
	}

	if response, _ := doRequest(http.MethodGet, "/cars", "", nil); response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("missing key: expected 401, got %d", response.StatusCode)
	}
	if response, _ := doRequest(http.MethodGet, "/cars", "invalid-key", nil); response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid key: expected 401, got %d", response.StatusCode)
	}

	vin := "T" + strings.ToUpper(keyA[:16])
	otherVIN := "U" + strings.ToUpper(keyA[16:32])
	response, _ := doRequest(http.MethodPost, "/cars", keyA, map[string]string{"vin": otherVIN})
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant A create %s: expected 200, got %d", otherVIN, response.StatusCode)
	}
	response, _ = doRequest(http.MethodGet, "/cars/"+otherVIN, keyB, nil)
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("cross-tenant VIN lookup: expected 404, got %d", response.StatusCode)
	}

	for _, key := range []string{keyA, keyB} {
		response, _ := doRequest(http.MethodPost, "/cars", key, map[string]string{"vin": vin})
		if response.StatusCode != http.StatusOK {
			t.Fatalf("same VIN registration for tenant: expected 200, got %d", response.StatusCode)
		}
	}

	response, _ = doRequest(http.MethodPatch, "/cars/"+vin+"/simulation", keyA, map[string]bool{"simulation_enabled": false})
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant A simulation update: expected 200, got %d", response.StatusCode)
	}
	response, body := doRequest(http.MethodGet, "/cars/"+vin, keyB, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant B same-VIN lookup: expected 200, got %d", response.StatusCode)
	}
	var car struct {
		SimulationEnabled bool `json:"simulation_enabled"`
	}
	if err := json.Unmarshal(body, &car); err != nil {
		t.Fatal(err)
	}
	if !car.SimulationEnabled {
		t.Fatal("tenant A update changed tenant B's same-VIN car")
	}

	telemetry := map[string]any{
		"vin": vin, "tire_pressure": 25, "engine_temp": 120, "battery_voltage": 9,
		"oil_level": 20, "brake_condition": 40,
	}
	response, _ = doRequest(http.MethodPost, "/telemetry", keyB, telemetry)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant B telemetry: expected 200, got %d", response.StatusCode)
	}
	response, body = doRequest(http.MethodGet, "/cars/"+vin+"/status", keyA, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant A status: expected 200, got %d", response.StatusCode)
	}
	if got := strings.TrimSpace(string(body)); got != "{}" {
		t.Fatalf("tenant A should not see tenant B status, got %s", got)
	}
	response, body = doRequest(http.MethodGet, "/cars/"+vin+"/events", keyA, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant A events: expected 200, got %d", response.StatusCode)
	}
	if strings.TrimSpace(string(body)) != "null" && strings.TrimSpace(string(body)) != "[]" {
		t.Fatalf("tenant A should not see tenant B events, got %s", body)
	}
	response, body = doRequest(http.MethodGet, "/cars/"+vin+"/status", keyB, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant B status: expected 200, got %d", response.StatusCode)
	}
	var status map[string]any
	if err := json.Unmarshal(body, &status); err != nil {
		t.Fatalf("decode tenant B status: %v; body: %s", err, body)
	}
	statusFields := []string{
		"tire_pressure_status",
		"engine_temp_status",
		"battery_voltage_status",
		"oil_level_status",
		"brake_condition_status",
	}
	observedStatusFields := 0
	for _, field := range statusFields {
		value, ok := status[field]
		if !ok {
			t.Fatalf("tenant B status is missing %q: %s", field, body)
		}
		switch value {
		case "NORMAL", "WARNING", "CRITICAL":
			observedStatusFields++
		default:
			t.Fatalf("tenant B status field %q has unexpected value %v", field, value)
		}
	}
	if observedStatusFields != 5 {
		t.Fatalf("tenant B should see all five metric statuses, saw %d: %s", observedStatusFields, body)
	}
	response, body = doRequest(http.MethodGet, "/cars/"+vin+"/events", keyB, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant B events: expected 200, got %d", response.StatusCode)
	}
	var events []map[string]any
	if err := json.Unmarshal(body, &events); err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 {
		t.Fatal("tenant B should see its own critical health events")
	}

	_, err = database.ExecContext(ctx,
		`INSERT INTO attribute_status (tenant_id, vin, tire_pressure_value, tire_pressure_status,
		 engine_temp_value, engine_temp_status, battery_voltage_value, battery_voltage_status,
		 oil_level_value, oil_level_status, brake_condition_value, brake_condition_status)
		 VALUES (?, ?, 32, 'NORMAL', 80, 'NORMAL', 13, 'NORMAL', 80, 'NORMAL', 80, 'NORMAL')`,
		tenantB, otherVIN)
	if err == nil {
		t.Fatal("database accepted a status row with tenant B and tenant A's car")
	}
	if tenantA == tenantB {
		t.Fatal("provisioned test tenants unexpectedly share an ID")
	}

	response, body = doRequest(http.MethodGet, "/cars", keyB, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tenant B car list: expected 200, got %d", response.StatusCode)
	}
	var cars []struct {
		VIN string `json:"vin"`
	}
	if err := json.Unmarshal(body, &cars); err != nil {
		t.Fatal(err)
	}
	for _, listed := range cars {
		if listed.VIN == otherVIN {
			t.Fatalf("tenant B listing included tenant A's VIN %s", otherVIN)
		}
	}
}
