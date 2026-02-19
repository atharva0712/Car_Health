package repository

import (
	"car-monitoring-system/internal/models"
	"database/sql"
)

type TelemetryRepository struct {
	DB *sql.DB
}

func NewTelemetryRepository(db *sql.DB) *TelemetryRepository {
	return &TelemetryRepository{DB: db}
}

// Event Insert
func (r *TelemetryRepository) InsertEvent(vin string, attribute string, status models.Status) error {
	query := `INSERT INTO health_events (vin, attribute_name, event_type) VALUES (?,?,?)`

	_, err := r.DB.Exec(query, vin, attribute, status)
	return err
}

// Close Events

func (r *TelemetryRepository) CloseEvents(vin string, attribute string) error {
	query := `UPDATE health_events 
				SET end_time = NOW()
				duration_second = TIMESTAMPDIFF(SECOND, start_time, NOW())
				WHERE vin=? AND attribute_name=? AND end_time IS NULL`
	_, err := r.DB.Exec(query, vin, attribute)
	return err
}

// Attribute Status Upsert

func (r *TelemetryRepository) UpdateStatus(vin string, telemetry models.Telemetry, statuses map[string]models.Status) error {
	query := `INSERT INTO attribute_status(vin, tire_pressure_value, tire_pressure_status, engine_temp_value, engine_temp_status, battery_voltage_value, battery_voltage_status, oil_level_value, oil_level_status, brake_condition_value, brake_condition_status) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	tire_pressure_value=?,
	tire_pressure_status=?,
	engine_temp_value=?,
	engine_temp_status=?,
	battery_voltage_value=?,
	battery_voltage_status=?
	oil_level_value=?,
	oil_level_status=?,
	brake_condition_value=?,
	brake_condition_status=?`

	_,err := r.DB.Exec(query,vin,telemetry.TirePressure, statuses["tire"],
						telemetry.EngineTemp, statuses["engine"],
						telemetry.BatteryVoltage, statuses["battery"],
						telemetry.OilLevel, statuses["oil"],
						telemetry.BrakeCondition, statuses["brake"],
						// update

						telemetry.TirePressure, statuses["tire"],
						telemetry.EngineTemp, statuses["engine"],
						telemetry.BatteryVoltage, statuses["battery"],
						telemetry.OilLevel, statuses["oil"],
						telemetry.BrakeCondition, statuses["brake"],
					)

	return err
}

