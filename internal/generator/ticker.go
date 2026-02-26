package generator

import (
	"bytes"
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Ticker struct {
	generator  *TelemetryGenerator
	backendURL string
	interval   time.Duration
}

// NewTicker creates a new Ticker instance with the provided configuration.
func NewTicker(cfg *config.Config) *Ticker {
	return &Ticker{
		generator:  NewTelemetryGenerator(cfg),
		backendURL: cfg.Backend.URL,
		interval:   time.Duration(cfg.Generation.IntervalSeconds) * time.Second,
	}
}

// Start begins the ticker that generates and sends telemetry data at regular intervals.

func (t *Ticker) Start() {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	log.Printf("Starting telemetry ticker with interval: %v\n", t.interval)

	t.tick() // Generate immediately on start

	// Then run on schedule
	for range ticker.C {
		t.tick()
	}
}

func (t *Ticker) tick() {
	//Fetch all cars from backend

	vins, err := t.fetchVINs()
	if err != nil {
		log.Println("Error fetching VINs:", err)
		return
	}

	log.Printf("Generating telemetry for %d cars", len(vins))

	// Generate and send telemetry for each car

	for _, vin := range vins {
		payload := t.generator.GenerateTelemetry(vin)
		err := t.sendTelemetry(payload)
		if err != nil {
			log.Println("Error sending telemetry for VIN", vin, ":", err)
		}
	}
}

func (t *Ticker) fetchVINs() ([]string, error) {
	resp, err := http.Get(t.backendURL + "/cars")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cars []models.Car
	err = json.NewDecoder(resp.Body).Decode(&cars)
	if err != nil {
		return nil, err
	}

	// Filter only simulation_enabled = true

	var vins []string
	for _, car := range cars {
		if car.SimulationEnabled {
			vins = append(vins, car.VIN)
		}
	}
	return vins, nil

	// vins := make([]string, len(cars))
	// for i, car := range cars {
	// 	vins[i] = car.VIN
	// }
	// return vins, nil
}

func (t *Ticker) sendTelemetry(payload models.TelemetryPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		t.backendURL+"/telemetry",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend returned status: %d", resp.StatusCode)
	}
	return nil
}
