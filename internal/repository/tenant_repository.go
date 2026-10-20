package repository

import (
	"context"
	"database/sql"
)

type TenantRepository struct {
	DB *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{DB: db}
}

func (r *TenantRepository) GetTenantIDByAPIKeyHash(ctx context.Context, hash string) (int64, error) {
	var id int64
	err := r.DB.QueryRowContext(ctx, `SELECT id FROM tenants WHERE api_key_hash = ?`, hash).Scan(&id)
	return id, err
}

func (r *TenantRepository) Create(ctx context.Context, name, keyHash string) (int64, error) {
	result, err := r.DB.ExecContext(ctx,
		`INSERT INTO tenants (name, api_key_hash) VALUES (?, ?)`, name, keyHash)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
