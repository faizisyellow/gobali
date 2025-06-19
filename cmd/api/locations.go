package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

type CreateLocationPayload struct {
	Area string `json:"area" validate:"required,min=3"`
}

type LocationResponse struct {
	Id        int    `json:"id"`
	Area      string `json:"area"`
	CreatedAt string `json:"created_at"`
}

type UpdateLocationPayload struct {
	CreateLocationPayload
}

// @Summary		Create location
// @Description	Create new location
// @Tags			Locations
// @Accept			json
// @Produce		json
// @Param			Payload	body		CreateLocationPayload	true	"Payload create new location"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		400		{object}	main.WriteJSONError.envelope
// @Failure		409		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/locations [post]
func (app *application) CreateLocationHandler(w http.ResponseWriter, r *http.Request) {
	payload := &CreateLocationPayload{}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := app.repository.Location.Create(r.Context(), payload.Area)
	if err != nil {
		switch err {
		case repository.ErrDuplicateLocation:
			app.conflictErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "create location successfull"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Get location
// @Description	Get Location By ID
// @Tags			Locations
// @Produce		json
// @Param			ID	path		int	true	"Location ID"
// @Success		200	{object}	main.jsonResponse.envelope{data=repository.Location}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/locations/{ID} [get]
func (app *application) GetLocationByIdHandler(w http.ResponseWriter, r *http.Request) {
	locationId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(locationId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	location, err := app.repository.Location.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusOK, location); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Get locations
// @Description	Get All Locations
// @Tags			Locations
// @Produce		json
// @Success		200	{object}	main.jsonResponse.envelope{data=[]LocationResponse{}}
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/locations [get]
func (app *application) GetLocationsHandler(w http.ResponseWriter, r *http.Request) {

	location, err := app.repository.Location.GetLocations(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	locsRes := []*LocationResponse{}

	for _, l := range location {
		loc := &LocationResponse{}

		loc.Id = l.Id
		loc.Area = l.Area
		loc.CreatedAt = l.CreatedAt

		locsRes = append(locsRes, loc)
	}

	if err := app.jsonResponse(w, http.StatusOK, locsRes); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Update Location
// @Description	Update new Location
//
// @Tags			Locations
//
// @Accept			json
// @Produce		json
// @Param			ID		path		int						true	"Location Id"
// @Param			payload	body		UpdateLocationPayload	true	"body new location"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		404		{object}	main.WriteJSONError.envelope
// @Failure		400		{object}	main.WriteJSONError.envelope
// @Failure		409		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/locations/{ID} [put]
func (app *application) UpdateLocationHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UpdateLocationPayload{}

	locationId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(locationId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	location, err := app.repository.Location.GetByID(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if payload.Area == location.Area {
		alreadyExist := errors.New("can not update from previous area")
		app.conflictErrorResponse(w, r, alreadyExist)
		return
	}

	location.Area = payload.Area

	if err := app.repository.Location.Update(ctx, location); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "update location successfull"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Delete Location
// @Description	Delete Location by ID
// @Tags			Locations
// @Param			ID	path	int	true	"Location ID"
// @Success		204
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/locations/{ID} [delete]
func (app *application) DeleteLocationHandler(w http.ResponseWriter, r *http.Request) {

	locationId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(locationId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	location, err := app.repository.Location.GetByID(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.repository.Location.Delete(ctx, location.Id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.responseNoContent(w); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
