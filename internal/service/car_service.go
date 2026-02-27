package service

import (
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/repository"
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

func (s *CarService) RegisterCar(vin string) error {
	s.Log.Debugf("Registering car with VIN: %s", vin)
	car := models.Car{
		VIN: vin,
		SimulationEnabled: true,
	}
	err := s.Repo.Create(car)
	if err != nil {
		s.Log.Errorf("Failed to register car: %v", err)
	}
	return err
}

func (s *CarService) UpdateSimulationStatus(vin string, enabled bool) error {
	s.Log.Debugf("Updating simulation status for VIN: %s to %v", vin, enabled)
	return s.Repo.UpdateSimulationStatus(vin,enabled)
}

func (s *CarService) GetCars() ([]models.Car, error) {
	s.Log.Debug("Fetching all cars")
	return s.Repo.GetAll()
}

func (s *CarService) GetCar(vin string) (models.Car, error) {
	s.Log.Debugf("Fetching car with VIN: %s", vin)
	return s.Repo.GetByVin(vin)
}

func (s *CarService) DeleteCar(vin string) error {
	s.Log.Debugf("Deleting car with VIN: %s", vin)
	return s.Repo.Delete(vin)
}
