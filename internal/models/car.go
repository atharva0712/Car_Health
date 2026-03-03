package models

import "time"

type Car struct {
	VIN               string    `json:"vin"`
	SimulationEnabled bool      `json:"simulation_enabled"`
	RegisteredAt      time.Time `json:"registered_at"`
}
