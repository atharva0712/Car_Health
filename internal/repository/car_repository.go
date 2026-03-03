package repository

import (
	"car-monitoring-system/internal/models"
	"database/sql"
)

type CarRepository struct {
	DB *sql.DB // pointer to the database connection pool, inject the DB dependency into repo
}

func NewCarRepository(db *sql.DB) *CarRepository { //returns a pointer to a new CarRepository instance initialized with the provided database connection
	return &CarRepository{DB: db}
}

// Create Car

func (r *CarRepository) Create(car models.Car) error {
	query := `INSERT INTO cars (vin,simulation_enabled) VALUES (?,?)`

	_, err := r.DB.Exec(query, car.VIN, car.SimulationEnabled) // Exec used for Insert,Update,Delete

	return err
}

// Get All Cars

func (r *CarRepository) GetAll() ([]models.Car, error) {

	query := `SELECT vin, simulation_enabled, registered_at FROM cars`

	rows, err := r.DB.Query(query) // Query used when multiple rows are expected

	if err != nil {
		return nil, err
	}

	defer rows.Close() //release db connection back to pool

	cars := []models.Car{}

	for rows.Next() {
		var car models.Car
		err := rows.Scan(&car.VIN, &car.SimulationEnabled, &car.RegisteredAt) // copies column values into struct fields
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

	query := `SELECT vin, simulation_enabled, registered_at FROM cars WHERE vin=?`

	err := r.DB.QueryRow(query, vin).Scan(&car.VIN, &car.SimulationEnabled, &car.RegisteredAt) //QueryRow : when expecting exactly one row

	return car, err
}

// DELETE CAR

func (r *CarRepository) Delete(vin string) error {
	query := "DELETE FROM cars WHERE vin=?"

	result, err := r.DB.Exec(query, vin)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected() // To prevent deletion of  non existent VIN, if there is no VIN it will return no error
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Toggling Simulation Enabled

func (r *CarRepository) UpdateSimulationStatus(vin string, enabled bool) error {
	query := `UPDATE cars SET simulation_enabled = ? WHERE vin = ?`
	result, err := r.DB.Exec(query, enabled, vin)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
