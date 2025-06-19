package repository

import (
	"context"
	"database/sql"
	"strings"
)

type TypesRepository struct {
	db *sql.DB
}

type Type struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"updated_at"`
}

type SelectedType struct {
	Name string `json:"name"`
}

func (t *TypesRepository) Create(ctx context.Context, name string) error {
	query := `INSERT INTO types(name) VALUE(?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := t.db.ExecContext(ctx, query, &name)
	if err != nil {
		duplicateKey := "Error 1062"
		switch {
		case strings.Contains(err.Error(), duplicateKey):
			return ErrDuplicateTypes
		default:
			return err
		}
	}

	return nil
}

func (t *TypesRepository) GetByID(ctx context.Context, id int) (*Type, error) {
	query := `SELECT  id, name, created_at,updated_at FROM types WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	ty := &Type{}

	err := t.db.QueryRowContext(ctx, query, id).Scan(&ty.Id, &ty.Name, &ty.CreatedAt, &ty.UpdateAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return ty, nil
}

func (t *TypesRepository) GetTypes(ctx context.Context) ([]*Type, error) {
	query := `SELECT id, name, created_at FROM types`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	types := []*Type{}

	for rows.Next() {
		ty := &Type{}
		err := rows.Scan(&ty.Id, &ty.Name, &ty.CreatedAt)
		if err != nil {
			return nil, err
		}

		types = append(types, ty)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return types, nil
}

func (t *TypesRepository) Update(ctx context.Context, ty *Type) error {
	query := `UPDATE types SET name = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := t.db.ExecContext(ctx, query, &ty.Name, &ty.Id)
	if err != nil {
		duplicateKey := "Error 1062"
		switch {
		case strings.Contains(err.Error(), duplicateKey):
			return ErrDuplicateTypes
		default:
			return err
		}
	}

	return nil
}

func (t *TypesRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM types WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := t.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
