package models

type Status string

const (
	NORMAL Status = "NORMAL"
	WARNING Status = "WARNING"
	CRITICAL Status = "CRITICAL"
)