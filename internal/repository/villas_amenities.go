package repository

import (
	"context"
	"database/sql"
	"strings"
)

type VillasAmenitiesRepository struct {
	db *sql.DB
}

type VillasAmenities struct {
	VillaId   int `json:"villa_id"`
	AmenityId int `json:"amenity_id"`
}

func (va *VillasAmenitiesRepository) Create(ctx context.Context, VillaAmenities *VillasAmenities) error {
	query := `INSERT INTO villas_amenities(villa_id,amenity_id) VALUES(?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := va.db.ExecContext(ctx, query, VillaAmenities.VillaId, VillaAmenities.AmenityId)
	if err != nil {

		duplicateKey := "Error 1062"

		switch {
		case strings.Contains(err.Error(), duplicateKey):
			return ErrDuplicateVillaAmenity
		default:
			return err
		}
	}

	return nil
}

func (va *VillasAmenitiesRepository) Delete(ctx context.Context, villaId, amenityId int) error {
	query := `DELETE FROM villas_amenities WHERE villa_id = ? AND amenity_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := va.db.ExecContext(ctx, query, villaId, amenityId)
	if err != nil {
		return err
	}

	return nil
}
