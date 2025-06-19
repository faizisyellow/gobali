package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

type bookingkey string

type StatusBooking int

type CreateBookingPayload struct {
	VillaId       int    `json:"villa_id" validate:"required"`
	VillaName     string `json:"villa_name" validate:"required"`
	VillaLocation string `json:"villa_location" validate:"required"`
	VillaPrice    int    `json:"villa_price" validate:"required"`
	StartAt       string `json:"start_at" validate:"required"`
	EndAt         string `json:"end_at" validate:"required"`
	TotalPrice    int    `json:"total_price" validate:"required,min=1"`
	UserId        int    `json:"user_id" validate:"required"`
	FirstName     string `json:"first_name" validate:"required,min=1"`
	LastName      string `json:"last_name" validate:"required,min=1"`
	Email         string `json:"email" validate:"required,min=1"`
	Guest         int    `json:"guest" validate:"required,min=1"`
}

type UpdateBookingStatus struct {
	Status string `json:"status" validate:"required"`
}

const (
	StatusOpen StatusBooking = iota
	StatusCheckIn
	StatusComplete
)

var (
	ErrAlreadyBooked error      = errors.New("this villa already booked between these days")
	bookingctx       bookingkey = "bookings"
)

var Status = map[StatusBooking]string{
	StatusOpen:     "open",
	StatusCheckIn:  "check_in",
	StatusComplete: "complete",
}

// @Summary		Check in Booking
// @Description	Check in Booking By ID
// @Tags			Bookings
// @Produce		json
// @Accept			json
// @Param			payload	body		UpdateBookingStatus	true	"status booking"
// @Param			Id		path		int					true	"booking id"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Success		400		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/bookings/{Id}/check-in [patch]
func (app *application) CheckInHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UpdateBookingStatus{}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Status != Status[StatusCheckIn] {
		app.badRequestResponse(w, r, fmt.Errorf("status is check_in"))
		return
	}

	booking := GetBookingFromContext(r)

	switch booking.Status {
	case Status[StatusCheckIn]:
		app.badRequestResponse(w, r, fmt.Errorf("already check in"))
		return
	case Status[StatusComplete]:
		app.badRequestResponse(w, r, fmt.Errorf("can not check in because, already check out"))
		return
	}

	if err := app.repository.Bookings.UpdateBookingStatus(r.Context(), booking.Id, payload.Status); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "check in successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Check out Booking
// @Description	Check out Booking By ID
// @Tags			Bookings
// @Produce		json
// @Accept			json
// @Param			payload	body		UpdateBookingStatus	true	"status booking"
// @Param			Id		path		int					true	"booking id"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Success		400		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/bookings/{Id}/check-out [patch]
func (app *application) CheckOutHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UpdateBookingStatus{}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Status != Status[StatusComplete] {
		app.badRequestResponse(w, r, fmt.Errorf("status is complete"))
		return
	}

	booking := GetBookingFromContext(r)

	switch booking.Status {
	case Status[StatusOpen]:
		app.badRequestResponse(w, r, fmt.Errorf("need to check in first"))
		return
	case Status[StatusComplete]:
		app.badRequestResponse(w, r, fmt.Errorf("already checkout"))
		return
	}

	if err := app.repository.Bookings.UpdateBookingStatus(r.Context(), booking.Id, payload.Status); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "check out successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Create Booking
// @Description	Create Booking
// @Tags			Bookings
// @Produce		json
// @Accept			json
// @Param			payload	body		CreateBookingPayload	true	"payload create booking"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Success		400		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/bookings [post]
func (app *application) CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	payload := &CreateBookingPayload{}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	bookingExist, err := app.repository.Bookings.GetBookingVillaByDate(ctx, payload.StartAt, payload.EndAt, payload.VillaId)
	if err != nil {
		if !errors.Is(err, repository.ErrNoRows) {
			app.internalServerError(w, r, err)
			return
		}
	}

	// if bookings not exist, can create new booking
	if bookingExist == nil {
		newBook := &repository.Booking{}
		newBook.VillaId = payload.VillaId
		newBook.VillaName = payload.VillaName
		newBook.VillaLocation = payload.VillaLocation
		newBook.VillaPrice = payload.VillaPrice
		newBook.StartAt = payload.StartAt
		newBook.EndAt = payload.EndAt
		newBook.TotalPrice = payload.TotalPrice
		newBook.UserId = payload.UserId
		newBook.FirstName = payload.FirstName
		newBook.LastName = payload.LastName
		newBook.Email = payload.Email
		newBook.Guest = payload.Guest

		if err := app.repository.Bookings.Create(ctx, newBook); err != nil {
			app.internalServerError(w, r, err)
			return
		}
	} else {
		app.badRequestResponse(w, r, ErrAlreadyBooked)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "villa booked successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// @Summary		Get Bookings
// @Description	Get All Bookings
// @Tags			Bookings
// @Produce		json
// @Accept			json
// @Success		200	{object}	main.jsonResponse.envelope{data=[]repository.Booking}
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/bookings [get]
func (app *application) GetBookingsHandler(w http.ResponseWriter, r *http.Request) {

	bookings, err := app.repository.Bookings.GetBookings(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, bookings); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Get Booking
// @Description	Get Booking By ID
// @Tags			Bookings
// @Produce		json
// @Accept			json
// @Param			Id	path		int	true	"Booking ID"
// @Success		200	{object}	main.jsonResponse.envelope{data=repository.Booking}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/bookings/{Id} [get]
func (app *application) GetBookingByIdHandler(w http.ResponseWriter, r *http.Request) {
	booking := GetBookingFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, booking); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Delete Booking
// @Description	Delete Booking By ID
// @Tags			Bookings
// @Produce		json
// @Accept			json
// @Param			Id	path	int	true	"Booking ID"
// @Success		204
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/bookings/{Id} [delete]
func (app *application) DeleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	booking := GetBookingFromContext(r)

	err := app.repository.Bookings.Delete(r.Context(), booking.Id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.responseNoContent(w); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) BookingContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookingId := chi.URLParam(r, "bookingID")

		id, err := strconv.Atoi(bookingId)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		booking, err := app.repository.Bookings.GetById(ctx, id)
		if err != nil {
			switch err {
			case repository.ErrNoRows:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, bookingctx, booking)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetBookingFromContext(r *http.Request) *repository.Booking {
	booking := r.Context().Value(bookingctx).(*repository.Booking)

	return booking
}
