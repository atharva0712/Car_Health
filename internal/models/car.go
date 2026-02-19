package models

import "time"

type Car struct {
	VIN          string    `json:"vin"`
	RegisteredAt time.Time `json:"registered_at"`
}