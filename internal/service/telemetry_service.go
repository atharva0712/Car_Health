package service

import (
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/repository"
)

type TelemetryService struct {
	Repo *repository.TelemetryRepository
}

func NewTelemetryService(r *repository.TelemetryRepository) * TelemetryService{
	return &TelemetryService{Repo: r}
}

// Main Telemetry Processor

func (s *TelemetryService) ProcessTelemetry(t models.Telemetry) error {
	// Classifying New Status

	newStatus := map[string]models.Status{
		"tire_pressure" : ClassifyTirePressure(t.TirePressure),
		"engine_temp" : ClassifyEngineTemp(t.EngineTemp),
		"battery_voltage" : ClassifyBattery(t.BatteryVoltage),
		"oil_level" : ClassifyOil(t.OilLevel),
		"brake_condition" : ClassifyBrake(t.BrakeCondition),
	}

	// Checking Previous Status

	for attribute,newStat := range newStatus{
		
		oldStat, err := s.Repo.GetCurrentStatus(t.VIN,attribute)
		if err!=nil{
			return err
		}

		// Handle Transition
		if oldStat == newStat{
			continue
		}

		// Closing old event if exists
		if oldStat == models.WARNING || oldStat == models.CRITICAL{
			err = s.Repo.CloseEvents(t.VIN,attribute)
			if err != nil{
				return  err
			}
		}
		// Opening New Event if needed
		if newStat == models.WARNING || newStat == models.CRITICAL{
			err = s.Repo.InsertEvent(t.VIN,attribute,newStat)
			if err != nil{
				return err
			}
		}
	}

	// Upserting Attribute status
	return s.Repo.UpdateAttributeStatus(t.VIN,t,newStatus)
}