package models

import "time"

type Car struct {
	VIN          string    `json:"vin"`
	SimulationEnabled bool `json:"simulation_enabled"`
	LastTelemetryAt *time.Time `json:"last_telemetry_at"`
	RegisteredAt time.Time `json:"registered_at"`	
}