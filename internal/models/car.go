package models

import "time"

type Car struct {
	VIN          string    `json:"vin"`
	RegisteredAt time.Time `json:"registered_at"`
}

type TelemetryPayload struct {
	VIN            string  `json:"vin"`
	TirePressure   float64 `json:"tire_pressure"`
	EngineTemp     float64 `json:"engine_temp"`
	BatteryVoltage float64 `json:"battery_voltage"`
	OilLevel       float64 `json:"oil_level"`
	BrakeCondition float64 `json:"brake_condition"`
}