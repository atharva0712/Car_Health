package repository

import (
	"car-monitoring-system/internal/models"
	"database/sql"
)

type CarRepository struct {
	DB *sql.DB // pointer to the database connection pool, inject the DB dependency into repo
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{DB: db}
}

// Create Car

func (r *CarRepository) Create(car models.Car) error {
	query := "INSERT INTO cars (vin) VALUES (?)"

	_, err := r.DB.Exec(query, car.VIN) // Exec used for Insert,Update,Delete

	return err
}

// Get All Cars

func (r *CarRepository) GetAll() ([]models.Car, error) {

	rows, err := r.DB.Query("SELECT vin, registered_at FROM cars") // Query used when multiple rows are expected

	if err != nil {
		return nil, err
	}

	defer rows.Close() //release db connection back to pool

	cars := []models.Car{}

	for rows.Next() {
		var car models.Car
		err := rows.Scan(&car.VIN, &car.RegisteredAt) // copies column values into struct fields
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

// GET CAR BY VIN

func (r *CarRepository) GetByVin(vin string) (models.Car, error) {
	var car models.Car

	query := "SELECT vin, registered_at FROM cars WHERE vin=?"

	err := r.DB.QueryRow(query, vin).Scan(&car.VIN, &car.RegisteredAt) //QueryRow : when expecting exactly one row

	return car, err
}

// DELETE CAR

func (r *CarRepository) Delete(vin string) error {
	query := "DELETE FROM cars WHERE vin=?"

	_, err := r.DB.Exec(query, vin)

	return err
}
