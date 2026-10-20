package service

import (
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/repository"
	"context"
)

type CarService struct {
	Repo *repository.CarRepository
	Log  *logger.Logger
}

func NewCarService(repo *repository.CarRepository, log *logger.Logger) *CarService {
	return &CarService{
		Repo: repo,
		Log:  log,
	}
}

func (s *CarService) RegisterCar(ctx context.Context, vin string) error {
	logger.FromContext(ctx, s.Log).Debugf("Registering car with VIN: %s", vin)
	car := models.Car{
		VIN:               vin,
		SimulationEnabled: true,
	}
	return s.Repo.Create(ctx, car)
}

func (s *CarService) UpdateSimulationStatus(ctx context.Context, vin string, enabled bool) error {
	logger.FromContext(ctx, s.Log).Debugf("Updating simulation status for VIN: %s to %v", vin, enabled)
	return s.Repo.UpdateSimulationStatus(ctx, vin, enabled)
}

func (s *CarService) GetCars(ctx context.Context) ([]models.Car, error) {
	logger.FromContext(ctx, s.Log).Debug("Fetching all cars")
	return s.Repo.GetAll(ctx)
}

func (s *CarService) GetCar(ctx context.Context, vin string) (models.Car, error) {
	logger.FromContext(ctx, s.Log).Debugf("Fetching car with VIN: %s", vin)
	return s.Repo.GetByVin(ctx, vin)
}

func (s *CarService) DeleteCar(ctx context.Context, vin string) error {
	logger.FromContext(ctx, s.Log).Debugf("Deleting car with VIN: %s", vin)
	return s.Repo.Delete(ctx, vin)
}
