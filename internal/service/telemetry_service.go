package service

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/repository"
	"errors"
)

type TelemetryService struct {
	Repo   *repository.TelemetryRepository
	Config *config.Config
	Log    *logger.Logger
}

func NewTelemetryService(r *repository.TelemetryRepository, cfg *config.Config, log *logger.Logger) *TelemetryService {
	return &TelemetryService{
		Repo:   r,
		Config: cfg,
		Log:    log,
	}
}

// ProcessTelemetry - main telemetry processor
func (s *TelemetryService) ProcessTelemetry(t models.TelemetryPayload) error {
	s.Log.Debugf("Processing telemetry for VIN: %s", t.VIN)

	// Validating if car exists and simulation is enabled
	enabled, err := s.Repo.IsSimulationEnabled(t.VIN)
	if err != nil {
		s.Log.Errorf("Error checking if simulation is enabled: %v", err)
		return err
	}
	if !enabled {
		s.Log.Warnf("Simulation disabled for car: %s", t.VIN)
		return errors.New("Simulation disabled for this car")
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
			s.Log.Errorf("Error getting current status: %v", err)
			return err
		}

		// If status unchanged, skip event logic
		if oldStat == newStat {
			continue
		}

		s.Log.Infof("Status transition for VIN %s, attribute %s: %v -> %v", t.VIN, attribute, oldStat, newStat)

		// Handle state transition
		err = s.handleTransition(t.VIN, attribute, oldStat, newStat)
		if err != nil {
			s.Log.Errorf("Error handling transition: %v", err)
			return err
		}
	}

	// Update all attribute statuses
	err = s.Repo.UpsertAttributeStatus(t.VIN, metrics, newStatus)
	if err != nil {
		s.Log.Errorf("Error updating attribute status: %v", err)
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
			s.Log.Errorf("Error closing event: %v", err)
			return err
		}
	}

	// Open new event if transitioning to WARNING or CRITICAL
	if newStatus == models.WARNING || newStatus == models.CRITICAL {
		err := s.Repo.InsertEvent(vin, attribute, newStatus)
		if err != nil {
			s.Log.Errorf("Error inserting event: %v", err)
			return err
		}
	}

	return nil
}

func (s *TelemetryService) GetCarStatus(vin string, attribute string) (interface{}, error) {
	s.Log.Debugf("Getting car status for VIN: %s, attribute: %s", vin, attribute)

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
	s.Log.Debugf("Getting events for VIN: %s, date: %s", vin, date)

	if vin == "" {
		s.Log.Warn("GetEvents called without VIN")
		return nil, errors.New("vin is required")
	}

	return s.Repo.GetEvents(vin, date)
}
