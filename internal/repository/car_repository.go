package repository

import "database/sql"

type CarRepository struct {
	DB *sql.DB
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{DB : db}
}