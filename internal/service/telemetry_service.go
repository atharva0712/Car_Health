package service

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/repository"
	"errors"
)

type TelemetryService struct {
	Repo   *repository.TelemetryRepository
	Config *config.Config
}

func NewTelemetryService(r *repository.TelemetryRepository, cfg *config.Config) *TelemetryService {
	return &TelemetryService{
		Repo:   r,
		Config: cfg,
	}
}

// ProcessTelemetry - main telemetry processor
func (s *TelemetryService) ProcessTelemetry(t models.TelemetryPayload) error {

	// Validating if car exists and simulation is enabled
	enabled, err := s.Repo.IsSimulationEnabled(t.VIN)
	if !enabled {
		return errors.New("Simulation disabled for this car")
	}
	if err != nil {
		return err
	}

	// Map of attribute name → value
	metrics := map[string]float64{
		"tire_pressure":   t.TirePressure,
		"engine_temp":     t.EngineTemp,
		"battery_voltage": t.BatteryVoltage,
		"oil_level":       t.OilLevel,
		"brake_condition": t.BrakeCondition,
	}

	// Classify new status for each metric
	newStatus := make(map[string]models.Status)
	for attribute, value := range metrics {
		newStatus[attribute] = s.classifyStatus(attribute, value)
	}

	// Process each metric
	for attribute, newStat := range newStatus {
		// Get current status from DB
		oldStat, err := s.Repo.GetCurrentStatus(t.VIN, attribute)
		if err != nil {
			return err
		}

		// If status unchanged, skip event logic
		if oldStat == newStat {
			continue
		}

		// Handle state transition
		err = s.handleTransition(t.VIN, attribute, oldStat, newStat)
		if err != nil {
			return err
		}
	}

	// Update all attribute statuses
	err = s.Repo.UpsertAttributeStatus(t.VIN, metrics, newStatus)
	if err != nil {
		return err
	}

	// Update last telemetry time
	return s.Repo.UpdateLastTelemetry(t.VIN)
}

// classifyStatus - uses config.yaml thresholds
func (s *TelemetryService) classifyStatus(attribute string, value float64) models.Status {
	metric := s.Config.Metrics[attribute]

	// Check CRITICAL
	if s.inRange(value, metric.Critical) {
		return models.CRITICAL
	}

	// Check WARNING
	if s.inRange(value, metric.Warning) {
		return models.WARNING
	}

	// Default NORMAL
	return models.NORMAL
}

// inRange - checks if value falls in range (handles tire_pressure dual range)
func (s *TelemetryService) inRange(value float64, r config.Range) bool {
	// Single range
	if value >= r.Min && value <= r.Max {
		return true
	}

	// Dual range (for tire_pressure warning: 28-29 OR 36-37)
	if r.Max2 > 0 && value >= r.Min2 && value <= r.Max2 {
		return true
	}

	return false
}

// handleTransition - implements the transition matrix
func (s *TelemetryService) handleTransition(vin, attribute string, oldStatus, newStatus models.Status) error {
	// Close old event if it was WARNING or CRITICAL
	if oldStatus == models.WARNING || oldStatus == models.CRITICAL {
		err := s.Repo.CloseEvent(vin, attribute)
		if err != nil {
			return err
		}
	}

	// Open new event if transitioning to WARNING or CRITICAL
	if newStatus == models.WARNING || newStatus == models.CRITICAL {
		err := s.Repo.InsertEvent(vin, attribute, newStatus)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TelemetryService) GetCarStatus(vin string, attribute string) (interface{}, error) {

	// Attribute-specific request
	if attribute != "" {
		return s.Repo.GetCurrentStatus(vin, attribute)
	}

	// Full status request
	return s.Repo.GetFullStatus(vin)
}

func (s *TelemetryService) GetEvents(
	vin string,
	date string,
) ([]models.HealthEvent, error) {

	if vin == "" {
		return nil, errors.New("vin is required")
	}

	return s.Repo.GetEvents(vin, date)
}
