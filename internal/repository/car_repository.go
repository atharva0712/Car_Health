package repository

import (
	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/tenant"
	"context"
	"database/sql"
)

type CarRepository struct {
	DB *sql.DB // pointer to the database connection pool, inject the DB dependency into repo
}

func NewCarRepository(db *sql.DB) *CarRepository { //returns a pointer to a new CarRepository instance initialized with the provided database connection
	return &CarRepository{DB: db}
}

// Create Car

func (r *CarRepository) Create(ctx context.Context, car models.Car) error {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	query := `INSERT INTO cars (tenant_id, vin, simulation_enabled) VALUES (?,?,?)`

	_, err = r.DB.ExecContext(ctx, query, tenantID, car.VIN, car.SimulationEnabled)

	return err
}

// Get All Cars

func (r *CarRepository) GetAll(ctx context.Context) ([]models.Car, error) {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT vin, simulation_enabled, registered_at FROM cars WHERE tenant_id = ?`

	rows, err := r.DB.QueryContext(ctx, query, tenantID)

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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

// GET CAR BY VIN

func (r *CarRepository) GetByVin(ctx context.Context, vin string) (models.Car, error) {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return models.Car{}, err
	}
	var car models.Car

	query := `SELECT vin, simulation_enabled, registered_at FROM cars WHERE tenant_id = ? AND vin = ?`

	err = r.DB.QueryRowContext(ctx, query, tenantID, vin).Scan(&car.VIN, &car.SimulationEnabled, &car.RegisteredAt)

	return car, err
}

// DELETE CAR

func (r *CarRepository) Delete(ctx context.Context, vin string) error {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	query := "DELETE FROM cars WHERE tenant_id = ? AND vin = ?"

	result, err := r.DB.ExecContext(ctx, query, tenantID, vin)
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

func (r *CarRepository) UpdateSimulationStatus(ctx context.Context, vin string, enabled bool) error {
	tenantID, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	query := `UPDATE cars SET simulation_enabled = ? WHERE tenant_id = ? AND vin = ?`
	result, err := r.DB.ExecContext(ctx, query, enabled, tenantID, vin)
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
