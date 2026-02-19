package repository

import (
	"car-monitoring-system/internal/models"
	"database/sql"
	"errors"
)

type TelemetryRepository struct {
	DB *sql.DB
}

func NewTelemetryRepository(db *sql.DB) *TelemetryRepository {
	return &TelemetryRepository{DB: db}
}

// Get Current Status of Attribute
func (r *TelemetryRepository) GetCurrentStatus(vin string, attribute string) (models.Status, error) {
	query := `
		SELECT ` + attribute + `_status		
		FROM attribute_status
		WHERE vin = ?
	`

	var status models.Status

	err := r.DB.QueryRow(query, vin).Scan(&status)

	if err == sql.ErrNoRows {

		return models.NORMAL, nil
	}

	return status, err
}

// Event Insert
func (r *TelemetryRepository) InsertEvent(vin string, attribute string, status models.Status) error {
	if status == models.NORMAL {
		return errors.New("Cannot insert NORMAL events")
	}
	query := `INSERT INTO health_events (vin, attribute_name, event_type) VALUES (?,?,?)`

	_, err := r.DB.Exec(query, vin, attribute, status)
	return err
}

// Close Events

func (r *TelemetryRepository) CloseEvents(vin string, attribute string) error {
	query := `UPDATE health_events 
				SET end_time = NOW(),
				duration_seconds = TIMESTAMPDIFF(SECOND, start_time, NOW())
				WHERE vin=? 
				AND attribute_name=? 
				AND end_time IS NULL`
	_, err := r.DB.Exec(query, vin, attribute)
	return err
}

// Attribute Status Upsert

func (r *TelemetryRepository) UpdateAttributeStatus(vin string, telemetry models.Telemetry, status map[string]models.Status) error {
	query := `INSERT INTO attribute_status(vin, tire_pressure_value, tire_pressure_status, engine_temp_value, engine_temp_status, battery_voltage_value, battery_voltage_status, oil_level_value, oil_level_status, brake_condition_value, brake_condition_status) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	tire_pressure_value = VALUES(tire_pressure_value),
	tire_pressure_status = VALUES(tire_pressure_status),

	engine_temp_value = VALUES(engine_temp_value),
	engine_temp_status = VALUES(engine_temp_status),

	battery_voltage_value = VALUES(battery_voltage_value),
	battery_voltage_status = VALUES(battery_voltage_status),

	oil_level_value = VALUES(oil_level_value),
	oil_level_status = VALUES(oil_level_status),

	brake_condition_value = VALUES(brake_condition_value),
	brake_condition_status = VALUES(brake_condition_status)`

	_, err := r.DB.Exec(query, vin, telemetry.TirePressure, status["tire_pressure"],
		telemetry.EngineTemp, status["engine_temp"],
		telemetry.BatteryVoltage, status["battery_voltage"],
		telemetry.OilLevel, status["oil_level"],
		telemetry.BrakeCondition, status["brake_condition"],
	)

	return err
}
