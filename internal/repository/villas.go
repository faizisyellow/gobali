package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
)

type VillasRepository struct {
	db *sql.DB
}

type Villa struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CategoryId  int              `json:"category_id"`
	LocationId  int              `json:"location_id"`
	Category    SelectedCategory `json:"category"`
	Location    SelectedLocation `json:"location"`
	Amenity     SelectedAmenity  `json:"amentiy"`
	MinGuest    int              `json:"min_guest"`
	Bedrooms    int              `json:"bedrooms"`
	Price       float64          `json:"price"`
	Baths       int              `json:"baths"`
	ImageUrls   []string         `json:"image_urls"`
	CreatedAt   string           `json:"created_at"`
	UpdateAt    string           `json:"updated_at"`
}

func (v *VillasRepository) Create(ctx context.Context, tx *sql.Tx, villa *Villa) (int64, error) {
	query := `INSERT INTO villas(image_urls,name,description,category_id,location_id,min_guest,bedrooms,price,baths)
	VALUES(?,?,?,?,?,?,?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	images, err := json.Marshal(villa.ImageUrls)
	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx, query,
		images,
		villa.Name,
		villa.Description,
		villa.CategoryId,
		villa.LocationId,
		villa.MinGuest,
		villa.Bedrooms,
		villa.Price,
		villa.Baths,
	)

	if err != nil {
		nexist := "Error 1452"

		switch {
		case strings.Contains(err.Error(), nexist):
			return 0, ErrCatOrLocNotExist
		default:
			return 0, err
		}
	}

	return res.LastInsertId()
}

func (v *VillasRepository) CreateVillasAmenities(ctx context.Context, tx *sql.Tx, villaId, amenityId int) error {
	query := `INSERT INTO villas_amenities(villa_id,amenity_id) VALUES(?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, villaId, amenityId)
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

func (v *VillasRepository) CreateVillaWithAmenity(ctx context.Context, payload *Villa) error {
	return withTx(v.db, ctx, func(tx *sql.Tx) error {

		villaId, err := v.Create(ctx, tx, payload)
		if err != nil {
			return err
		}

		if err := v.CreateVillasAmenities(ctx, tx, int(villaId), payload.Amenity.Id); err != nil {
			return err
		}

		return nil
	})
}

func (v *VillasRepository) GetById(ctx context.Context, id int) (*Villa, error) {
	query := `
	SELECT
		v.id,
		v.name,
		v.description,
		v.category_id,
		v.location_id,
		v.min_guest,
		v.bedrooms,
		v.baths,
		v.price,
		v.image_urls,
		c.id,
		c.name,
		l.id,
		l.area,
		v.created_at,
		v.updated_at
	FROM
    	villas v LEFT JOIN categories c ON v.category_id = c.id LEFT JOIN locations l ON v.location_id = l.id
		WHERE v.id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	villa := &Villa{}

	rowUrls := []uint8{}

	err := v.db.QueryRowContext(ctx, query, id).Scan(
		&villa.Id,
		&villa.Name,
		&villa.Description,
		&villa.CategoryId,
		&villa.LocationId,
		&villa.MinGuest,
		&villa.Bedrooms,
		&villa.Baths,
		&villa.Price,
		&rowUrls,
		&villa.Category.Id,
		&villa.Category.Name,
		&villa.Location.Id,
		&villa.Location.Area,
		&villa.CreatedAt,
		&villa.UpdateAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	err = json.Unmarshal(rowUrls, &villa.ImageUrls)
	if err != nil {
		return nil, err
	}

	return villa, nil
}

// TODO: SHOULD RETURN AMANITIES
func (v *VillasRepository) GetVillas(ctx context.Context) ([]*Villa, error) {
	query := `
	SELECT
		v.id,
		v.name,
		v.description,
		v.category_id,
		v.location_id,
		v.min_guest,
		v.bedrooms,
		v.baths,
		v.price,
		v.image_urls,
		c.id,
		c.name,
		l.id,
		l.area,
		a.id,
		a.name,
		v.created_at,
		v.updated_at
	FROM
    	villas v LEFT JOIN categories c ON v.category_id = c.id LEFT JOIN locations l ON v.location_id = l.id
		LEFT JOIN villas_amenities va ON v.id = va.villa_id LEFT JOIN amenities a ON va.amenity_id = a.id
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := v.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	villas := []*Villa{}

	for rows.Next() {
		villa := &Villa{}

		rowUrls := []uint8{}

		err := rows.Scan(
			&villa.Id,
			&villa.Name,
			&villa.Description,
			&villa.CategoryId,
			&villa.LocationId,
			&villa.MinGuest,
			&villa.Bedrooms,
			&villa.Baths,
			&villa.Price,
			&rowUrls,
			&villa.Category.Id,
			&villa.Category.Name,
			&villa.Location.Id,
			&villa.Location.Area,
			&villa.Amenity.Id,
			&villa.Amenity.Name,
			&villa.CreatedAt,
			&villa.UpdateAt,
		)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(rowUrls, &villa.ImageUrls)
		if err != nil {
			return nil, err
		}

		villas = append(villas, villa)
	}

	return villas, nil
}

func (v *VillasRepository) Update(ctx context.Context, villa *Villa) error {

	query := `UPDATE villas SET image_urls=?, name=?, description=?, min_guest=?, bedrooms=?, price=?, baths=?,location_id=?,category_id=?
	WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	images, err := json.Marshal(villa.ImageUrls)
	if err != nil {
		return err
	}

	_, err = v.db.ExecContext(ctx, query,
		images,
		&villa.Name,
		&villa.Description,
		&villa.MinGuest,
		&villa.Bedrooms,
		&villa.Price,
		&villa.Baths,
		&villa.LocationId,
		&villa.CategoryId,
		&villa.Id,
	)

	if err != nil {
		return err
	}

	return nil

}

func (v *VillasRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM villas WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := v.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
