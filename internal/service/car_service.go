package service

import (
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/repository"
)

type CarService struct {
	Repo *repository.CarRepository
}

func NewCarService(repo *repository.CarRepository) *CarService {
	return &CarService{Repo: repo}
}

func (s *CarService) RegisterCar(vin string) error {
	car := models.Car{
		VIN: vin,
	}
	return s.Repo.Create(car)
}

func (s *CarService) GetCars() ([]models.Car, error) {
	return s.Repo.GetAll()
}

func (s *CarService) GetCar(vin string) (models.Car, error) {
	return s.Repo.GetByVin(vin)
}

func (s *CarService) DeleteCar(vin string) error {
	return s.Repo.Delete(vin)
}
