package models

type HealthEvent struct {
	VIN             string  `db:"vin" json:"vin"`
	AttributeName   string  `db:"attribute_name" json:"attribute_name"`
	EventType       string  `db:"event_type" json:"event_type"` // WARNING or CRITICAL
	StartTime       string  `db:"start_time" json:"start_time"`
	EndTime         *string `db:"end_time" json:"end_time,omitempty"` // pointer to allow null
	DurationSeconds *int    `db:"duration_seconds" json:"duration_seconds,omitempty"`
}