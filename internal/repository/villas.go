package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/charmbracelet/log"
)

type VillasRepository struct {
	db *sql.DB
}

type Villa struct {
	Id          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CategoryId  int               `json:"category_id"`
	LocationId  int               `json:"location_id"`
	Category    SelectedCategory  `json:"category"`
	Location    SelectedLocation  `json:"location"`
	Amenity     []SelectedAmenity `json:"amentiy"`
	MinGuest    int               `json:"min_guest"`
	Bedrooms    int               `json:"bedrooms"`
	Price       float64           `json:"price"`
	Baths       int               `json:"baths"`
	ImageUrls   []string          `json:"image_urls"`
	CreatedAt   string            `json:"created_at"`
	UpdateAt    string            `json:"updated_at"`
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
		duplicate := "Error 1062"

		switch {
		case strings.Contains(err.Error(), nexist):
			return 0, ErrCatOrLocNotExist
		case strings.Contains(err.Error(), duplicate):
			return 0, ErrDuplicateVilla
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
		notExist := "Error 1452"

		switch {
		case strings.Contains(err.Error(), duplicateKey):
			log.Error(err)
			return ErrDuplicateVillaAmenity
		case strings.Contains(err.Error(), notExist):
			log.Error(err)
			return ErrAmenitiesNotExist
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

		for _, amenity := range payload.Amenity {
			if err := v.CreateVillasAmenities(ctx, tx, int(villaId), amenity.Id); err != nil {
				return err
			}
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
		a.id,
		a.name as amenity_name,
		t.name as type_amenity,
		v.created_at,
		v.updated_at
	FROM
    	villas v LEFT JOIN categories c ON v.category_id = c.id LEFT JOIN locations l ON v.location_id = l.id
		LEFT JOIN villas_amenities va ON va.villa_id = v.id LEFT JOIN amenities a ON va.amenity_id = a.id
		LEFT JOIN types t ON t.id = a.type_id
		WHERE v.id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	villa := &Villa{}

	rowUrls := []uint8{}

	rows, err := v.db.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		amenity := &SelectedAmenity{}

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
			&amenity.Id,
			&amenity.Name,
			&amenity.Type,
			&villa.CreatedAt,
			&villa.UpdateAt,
		)

		if err != nil {
			return nil, err
		}

		// if the row exist, just append the amenity
		if villa.Id != 0 {
			villa.Amenity = append(villa.Amenity, *amenity)
		}
	}

	// TODO: FIX UNEXPECTED UNMARSHAL
	log.Info(villa, "sup", "asdsad")

	err = json.Unmarshal(rowUrls, &villa.ImageUrls)
	if err != nil {
		return nil, err
	}

	return villa, nil
}

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
		a.name as amenity_name,
		t.name as type_amenity,
		v.created_at,
		v.updated_at
	FROM
    	villas v LEFT JOIN categories c ON v.category_id = c.id LEFT JOIN locations l ON v.location_id = l.id
		LEFT JOIN villas_amenities va ON v.id = va.villa_id LEFT JOIN amenities a ON va.amenity_id = a.id
		LEFT JOIN types t ON t.id = a.type_id
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := v.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	villas := []*Villa{}

	villaMap := make(map[int]*Villa)

	for rows.Next() {
		villa := &Villa{}

		amenity := &SelectedAmenity{}

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
			&amenity.Id,
			&amenity.Name,
			&amenity.Type,
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

		// if the row not exist add not the map
		if _, ok := villaMap[villa.Id]; !ok {
			villaMap[villa.Id] = villa
		}

		// update the amenity with the result of the amenity row
		villaMap[villa.Id].Amenity = append(villaMap[villa.Id].Amenity, *amenity)
	}

	// convert the map to slice
	for _, val := range villaMap {
		villas = append(villas, val)
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
