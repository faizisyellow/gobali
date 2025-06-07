package repository

import (
	"context"
	"database/sql"
	"strings"
)

type LocationsRepository struct {
	db *sql.DB
}

type Location struct {
	Id        int    `json:"id"`
	Area      string `json:"area"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"updated_at"`
}

func (l *LocationsRepository) Create(ctx context.Context, area string) error {
	query := `INSERT INTO locations(area) VALUE(?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := l.db.ExecContext(ctx, query, area)
	if err != nil {
		duplicateKey := "Error 1062"
		switch {
		case strings.Contains(err.Error(), duplicateKey):
			return ErrDuplicateLocation
		default:
			return err
		}
	}

	return nil
}

func (l *LocationsRepository) GetByID(ctx context.Context, id int) (*Location, error) {
	query := `SELECT id, area, created_at, updated_at FROM locations WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	location := &Location{}
	err := l.db.QueryRowContext(ctx, query, id).Scan(&location.Id, &location.Area, &location.CreatedAt, &location.UpdateAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return location, nil
}

func (l *LocationsRepository) GetLocations(ctx context.Context) ([]*Location, error) {
	query := `SELECT id, area, created_at, updated_at FROM locations`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := l.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	locations := []*Location{}

	for rows.Next() {
		location := &Location{}

		err := rows.Scan(&location.Id, &location.Area, &location.CreatedAt, &location.UpdateAt)
		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil

}

func (l *LocationsRepository) Update(ctx context.Context, location *Location) error {
	query := `UPDATE locations SET area = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := l.db.ExecContext(ctx, query, &location.Area, &location.Id)
	if err != nil {
		return err
	}

	return nil
}

func (l *LocationsRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM locations WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := l.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
