package repository

import (
	"context"
	"database/sql"
)

type Role struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdateAt    string `json:"update_at"`
}

type RolesRepository struct {
	db *sql.DB
}

func (r *RolesRepository) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id, name, level, description FROM roles WHERE name = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var role Role
	err := r.db.QueryRowContext(ctx, query, name).Scan(&role.Id, &role.Name, &role.Level, &role.Description)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RolesRepository) Create(ctx context.Context, payload *Role) error {
	query := `INSERT INTO roles(name, level, description) VALUES(?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, payload.Name, payload.Level, payload.Description)
	if err != nil {
		return err
	}

	return nil
}
