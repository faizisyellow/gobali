package repository

import (
	"context"
	"database/sql"
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
	Email         string  `json:"email"`
	Guest         int     `json:"guest"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
}

func (b *BookingsRepository) Create(ctx context.Context, newBooking *Booking) error {
	query := `INSERT INTO bookings(
	user_id,
	villa_id,
	villa_name,
	villa_location,
	villa_price,
	first_name,
	last_name,
	start_at,
	end_at,
	total_price,
	email,
	guest) 
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`

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
		newBooking.Email,
		newBooking.Guest,
	)

	if err != nil {
		return err
	}

	return nil
}

func (b *BookingsRepository) GetById(ctx context.Context, id int) (*Booking, error) {
	query := `SELECT id,first_name,last_name,status,villa_id,villa_name,villa_price,villa_location,total_price,
	start_at,end_at,created_at,updated_at,user_id,email,guest FROM bookings WHERE id = ?
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
		&booking.Email,
		&booking.Guest,
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

func (b *BookingsRepository) GetBookings(ctx context.Context, pq PaginatedBookingsQuery) ([]*Booking, error) {
	query := `SELECT id,first_name,last_name,status,villa_name,villa_price,villa_location,total_price,
	start_at,end_at,email,guest,villa_id,user_id,created_at,updated_at FROM bookings
	ORDER BY created_at ` + pq.Sort + ` LIMIT ? OFFSET ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := b.db.QueryContext(ctx, query, pq.Limit, pq.Offset)
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
			&booking.Email,
			&booking.Guest,
			&booking.VillaId,
			&booking.UserId,
			&booking.CreatedAt,
			&booking.UpdatedAt,
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
	start_at,end_at,email FROM bookings WHERE ? >= start_at AND ? <= end_at AND villa_id = ?`

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
		&booking.Email,
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

func (b *BookingsRepository) UpdateBookingStatus(ctx context.Context, bookId int, status string) error {
	query := `UPDATE bookings SET status = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := b.db.ExecContext(ctx, query, status, bookId)
	if err != nil {
		return err
	}

	return nil
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
