package repository

import (
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/tenant"
	"context"
	"database/sql"
	"errors"
)

type TelemetryRepository struct {
	DB *sql.DB
}

func NewTelemetryRepository(db *sql.DB) *TelemetryRepository {
	return &TelemetryRepository{DB: db}
}

// GetCurrentStatus - get current status of a specific attribute

func (r *TelemetryRepository) GetCurrentStatus(ctx context.Context, vin string, attribute string) (models.Status, error) {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return "", err
	}
	// Whitelist valid attributes to prevent SQL injection
	columnMap := map[string]string{
		"tire_pressure":   "tire_pressure_status",
		"engine_temp":     "engine_temp_status",
		"battery_voltage": "battery_voltage_status",
		"oil_level":       "oil_level_status",
		"brake_condition": "brake_condition_status",
	}

	column, ok := columnMap[attribute]
	if !ok {
		return "", errors.New("invalid attribute name")
	}

	query := "SELECT " + column + " FROM attribute_status WHERE tenant_id = ? AND vin = ?"

	var status models.Status
	err = r.DB.QueryRowContext(ctx, query, tenantID, vin).Scan(&status)

	if err == sql.ErrNoRows {
		// If no row exists yet, assume NORMAL
		return models.NORMAL, nil
	}

	return status, err
}

func (r *TelemetryRepository) GetFullStatus(ctx context.Context, vin string) (*models.AttributeStatus, error) {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
	SELECT
		vin,
		tire_pressure_value, tire_pressure_status,
		engine_temp_value, engine_temp_status,
		battery_voltage_value, battery_voltage_status,
		oil_level_value, oil_level_status,
		brake_condition_value, brake_condition_status
	FROM attribute_status
	WHERE tenant_id = ? AND vin = ?
`

	var status models.AttributeStatus

	err = r.DB.QueryRowContext(ctx, query, tenantID, vin).Scan(
		&status.VIN,

		&status.TirePressureValue,
		&status.TirePressureStatus,

		&status.EngineTempValue,
		&status.EngineTempStatus,

		&status.BatteryVoltageValue,
		&status.BatteryVoltageStatus,

		&status.OilLevelValue,
		&status.OilLevelStatus,

		&status.BrakeConditionValue,
		&status.BrakeConditionStatus,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &status, err
}

// InsertEvent - create new WARNING or CRITICAL event
func (r *TelemetryRepository) InsertEvent(ctx context.Context, vin string, attribute string, status models.Status) error {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	if status == models.NORMAL {
		return errors.New("cannot insert NORMAL events")
	}

	query := `INSERT INTO health_events (tenant_id, vin, attribute_name, event_type, start_time)
	          VALUES (?, ?, ?, ?, NOW())`

	_, err = r.DB.ExecContext(ctx, query, tenantID, vin, attribute, status)
	return err
}

// CloseEvent - close the active event for this attribute
func (r *TelemetryRepository) CloseEvent(ctx context.Context, vin string, attribute string) error {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	query := `UPDATE health_events 
	          SET end_time = NOW(),
	              duration_seconds = TIMESTAMPDIFF(SECOND, start_time, NOW())
	          WHERE tenant_id = ?
	            AND vin = ?
	            AND attribute_name = ? 
	            AND end_time IS NULL`

	_, err = r.DB.ExecContext(ctx, query, tenantID, vin, attribute)
	return err
}

// GetEvents - retrieve all events for a car, optionally filtered by date
func (r *TelemetryRepository) GetEvents(ctx context.Context, vin string, date string) ([]models.HealthEvent, error) {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	query := `
        SELECT
            vin,
            attribute_name,
            event_type,
            start_time,
            end_time,
            duration_seconds
        FROM health_events
        WHERE tenant_id = ? AND vin = ?
    `

	args := []any{tenantID, vin}

	if date != "" {
		query += `
            AND start_time >= ?
            AND start_time < DATE_ADD(?, INTERVAL 1 DAY)
        `
		args = append(args, date, date)
	}

	query += ` ORDER BY start_time DESC`

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.HealthEvent

	for rows.Next() {
		var event models.HealthEvent

		err := rows.Scan(
			&event.VIN,
			&event.AttributeName,
			&event.EventType,
			&event.StartTime,
			&event.EndTime,
			&event.DurationSeconds,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// UpsertAttributeStatus - insert or update all attribute values and statuses
func (r *TelemetryRepository) UpsertAttributeStatus(ctx context.Context, vin string, metrics map[string]float64, status map[string]models.Status) error {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	query := `INSERT INTO attribute_status (
		tenant_id, vin,
		tire_pressure_value, tire_pressure_status,
		engine_temp_value, engine_temp_status,
		battery_voltage_value, battery_voltage_status,
		oil_level_value, oil_level_status,
		brake_condition_value, brake_condition_status
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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

	_, err = r.DB.ExecContext(ctx, query,
		tenantID,
		vin,
		metrics["tire_pressure"], status["tire_pressure"],
		metrics["engine_temp"], status["engine_temp"],
		metrics["battery_voltage"], status["battery_voltage"],
		metrics["oil_level"], status["oil_level"],
		metrics["brake_condition"], status["brake_condition"],
	)

	return err
}

// Checking if Simulation is Enabled or Not

func (r *TelemetryRepository) IsSimulationEnabled(ctx context.Context, vin string) (bool, error) {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return false, err
	}
	query := `SELECT simulation_enabled FROM cars WHERE tenant_id = ? AND vin = ?`

	var enabled bool

	err = r.DB.QueryRowContext(ctx, query, tenantID, vin).Scan(&enabled)

	return enabled, err
}
