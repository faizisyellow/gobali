package repository

import (
	"context"
	"database/sql"
	"strings"
)

type AmenitiesRepository struct {
	db *sql.DB
}

type SelectedType struct {
	Name string `json:"name"`
}

type Amenity struct {
	Id        int          `json:"id"`
	Name      string       `json:"name"`
	TypeId    int          `json:"type_id"`
	Type      SelectedType `json:"type"`
	CreatedAt string       `json:"created_at"`
	UpdateAt  string       `json:"updated_at"`
}

func (a *AmenitiesRepository) Create(ctx context.Context, name string, typeId int) error {
	query := `INSERT INTO amenities(name,type_id) VALUE(?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := a.db.ExecContext(ctx, query, &name, &typeId)
	if err != nil {
		duplicateKey := "Error 1062"
		emptyKey := "Error 1452"
		switch {
		case strings.Contains(err.Error(), duplicateKey):
			return ErrDuplicateAmenities
		case strings.Contains(err.Error(), emptyKey):
			return ErrTypeNotExist
		default:
			return err
		}
	}

	return nil
}

func (a *AmenitiesRepository) GetByID(ctx context.Context, id int) (*Amenity, error) {
	query := `SELECT  a.id, a.name, a.type_id, t.name, a.created_at, a.updated_at
	 FROM amenities a LEFT JOIN types t ON a.type_id = t.id WHERE a.id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	am := &Amenity{}

	err := a.db.QueryRowContext(ctx, query, id).Scan(&am.Id, &am.Name, &am.TypeId, &am.Type.Name, &am.CreatedAt, &am.UpdateAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return am, nil
}

func (a *AmenitiesRepository) GetAmenities(ctx context.Context) ([]*Amenity, error) {
	query := `SELECT a.id, a.name, a.type_id, t.name, a.created_at FROM amenities a LEFT JOIN types t ON a.type_id = t.id `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	amenities := []*Amenity{}

	for rows.Next() {
		am := &Amenity{}
		err := rows.Scan(&am.Id, &am.Name, &am.TypeId, &am.Type.Name, &am.CreatedAt)
		if err != nil {
			return nil, err
		}

		amenities = append(amenities, am)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return amenities, nil
}

func (a *AmenitiesRepository) Update(ctx context.Context, amentity *Amenity) error {
	query := `UPDATE amenities SET name = ?, type_id = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := a.db.ExecContext(ctx, query, &amentity.Name, &amentity.TypeId, &amentity.Id)
	if err != nil {
		return err
	}

	return nil
}

func (a *AmenitiesRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM amenities WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
