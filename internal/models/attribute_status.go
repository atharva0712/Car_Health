package models

type AttributeStatus struct {
	VIN string `db:"vin" json:"vin"`

	TirePressureValue  float64 `db:"tire_pressure_value" json:"tire_pressure_value"`
	TirePressureStatus Status  `db:"tire_pressure_status" json:"tire_pressure_status"`

	EngineTempValue  float64 `db:"engine_temp_value" json:"engine_temp_value"`
	EngineTempStatus Status  `db:"engine_temp_status" json:"engine_temp_status"`

	BatteryVoltageValue  float64 `db:"battery_voltage_value" json:"battery_voltage_value"`
	BatteryVoltageStatus Status  `db:"battery_voltage_status" json:"battery_voltage_status"`

	OilLevelValue  float64 `db:"oil_level_value" json:"oil_level_value"`
	OilLevelStatus Status  `db:"oil_level_status" json:"oil_level_status"`

	BrakeConditionValue  float64 `db:"brake_condition_value" json:"brake_condition_value"`
	BrakeConditionStatus Status  `db:"brake_condition_status" json:"brake_condition_status"`
}
