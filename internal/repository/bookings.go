package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/faizisyellow/gobali/internal/helpers"
)

type BookingsRepository struct {
	db *sql.DB
}

type Booking struct {
	Id            int     `json:"id"`
	UserId        int     `json:"user_id"`
	VillaId       int     `json:"villa_id"`
	VillaName     string  `json:"villa_name"`
	VillaLocation string  `json:"villa_location"`
	VillaPrice    int     `json:"villa_price"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Status        string  `json:"status"`
	StartAt       string  `json:"start_at"`
	EndAt         string  `json:"end_at"`
	TotalPrice    int     `json:"total_price"`
	ExpireAt      string  `json:"expire_at"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
}

func (b *Booking) Include(fields ...string) *Booking {

	newBooking := helpers.CreateNewStructByReflect(b, fields...)

	return newBooking
}

func (b *BookingsRepository) Create(ctx context.Context, newBooking *Booking, exp time.Duration) error {
	query := `INSERT INTO bookings(user_id,villa_id,villa_name,villa_location,villa_price,first_name,last_name,start_at,end_at,total_price,expire_at) 
	VALUES(?,?,?,?,?,?,?,?,?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := b.db.ExecContext(ctx, query,
		newBooking.UserId,
		newBooking.VillaId,
		newBooking.VillaName,
		newBooking.VillaLocation,
		newBooking.VillaPrice,
		newBooking.FirstName,
		newBooking.LastName,
		newBooking.StartAt,
		newBooking.EndAt,
		newBooking.TotalPrice,
		time.Now().Add(exp),
	)

	if err != nil {
		return err
	}

	return nil
}

func (b *BookingsRepository) GetById(ctx context.Context, id int) (*Booking, error) {
	query := `SELECT id,first_name,last_name,status,villa_id,villa_name,villa_price,villa_location,total_price,
	start_at,end_at,created_at,updated_at,user_id FROM bookings WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	booking := Booking{}
	err := b.db.QueryRowContext(ctx, query, id).Scan(
		&booking.Id,
		&booking.FirstName,
		&booking.LastName,
		&booking.Status,
		&booking.VillaId,
		&booking.VillaName,
		&booking.VillaPrice,
		&booking.VillaLocation,
		&booking.TotalPrice,
		&booking.StartAt,
		&booking.EndAt,
		&booking.CreatedAt,
		&booking.UpdatedAt,
		&booking.UserId,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return &booking, nil
}

func (b *BookingsRepository) GetBookings(ctx context.Context) ([]*Booking, error) {
	query := `SELECT id,first_name,last_name,status,villa_name,villa_price,villa_location,total_price,
	start_at,end_at,expire_at,villa_id,user_id,created_at FROM bookings`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := b.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bookings := []*Booking{}

	for rows.Next() {
		booking := &Booking{}

		err := rows.Scan(
			&booking.Id,
			&booking.FirstName,
			&booking.LastName,
			&booking.Status,
			&booking.VillaName,
			&booking.VillaPrice,
			&booking.VillaLocation,
			&booking.TotalPrice,
			&booking.StartAt,
			&booking.EndAt,
			&booking.ExpireAt,
			&booking.VillaId,
			&booking.UserId,
			&booking.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (b *BookingsRepository) GetBookingVillaByDate(ctx context.Context, startAt, endAt string, villaId int) (*Booking, error) {
	query := `SELECT  id,first_name,last_name,status,villa_id,villa_name,villa_price,villa_location,total_price,
	start_at,end_at FROM bookings WHERE ? >= start_at AND ? <= end_at AND villa_id = ?`

	booking := &Booking{}

	err := b.db.QueryRowContext(ctx, query, startAt, endAt, villaId).Scan(
		&booking.Id,
		&booking.FirstName,
		&booking.LastName,
		&booking.Status,
		&booking.VillaId,
		&booking.VillaName,
		&booking.VillaPrice,
		&booking.VillaLocation,
		&booking.TotalPrice,
		&booking.StartAt,
		&booking.EndAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return booking, nil
}

func (b *BookingsRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM bookings WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := b.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
