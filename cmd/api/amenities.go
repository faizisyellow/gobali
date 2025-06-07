package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

type CreateAmenityPayload struct {
	Name   string `json:"name" validate:"min=3"`
	TypeId int    `json:"type_id"`
}

type AmenityResponse struct {
	Id        int                     `json:"id"`
	Name      string                  `json:"name"`
	TypeId    int                     `json:"type_id"`
	CreatedAt string                  `json:"created_at"`
	Type      repository.SelectedType `json:"type"`
}

type UpdateAmenityPayload struct {
	CreateAmenityPayload
}

// @Summary		Create Amenity
// @Description	Create Amenity
// @Tags			Amenities
// @Accept			json
// @Produce		json
// @Param			payload	body		CreateAmenityPayload	true	"json format payload"
//
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		400		{object}	main.WriteJSONError.envelope
// @Failure		409		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/amenities [POST]
func (app *application) CreateAmenityHandler(w http.ResponseWriter, r *http.Request) {
	payload := &CreateAmenityPayload{}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.repository.Amenities.Create(r.Context(), payload.Name, payload.TypeId); err != nil {
		switch err {
		case repository.ErrDuplicateAmenities:
			app.conflictErrorResponse(w, r, err)
		case repository.ErrTypeNotExist:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "amenity created successfull"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// @Summary		Get Amenity
// @Description	Get Amenity by ID
// @Tags			Amenities
// @Produce		json
// @Param			ID	path		int	true	"amenity id"
//
// @Success		200	{object}	main.jsonResponse.envelope{data=repository.Amenity}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/amenities/{ID} [GET]
func (app *application) GetAmenityByIDHandler(w http.ResponseWriter, r *http.Request) {
	amenityId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(amenityId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	cat, err := app.repository.Amenities.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusOK, cat); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Get Amenities
// @Description	Get all Amenities
// @Tags			Amenities
// @Produce		json
//
// @Success		200	{object}	main.jsonResponse.envelope{data=[]AmenityResponse}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/amenities [GET]
func (app *application) GetAmenitiesHandler(w http.ResponseWriter, r *http.Request) {

	amenities, err := app.repository.Amenities.GetAmenities(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	amenitiesRes := []*AmenityResponse{}

	for _, amenity := range amenities {
		newAmenity := &AmenityResponse{}

		newAmenity.Id = amenity.Id
		newAmenity.Name = amenity.Name
		newAmenity.TypeId = amenity.TypeId
		newAmenity.Type.Name = amenity.Type.Name
		newAmenity.CreatedAt = amenity.CreatedAt

		amenitiesRes = append(amenitiesRes, newAmenity)
	}

	if err := app.jsonResponse(w, http.StatusOK, amenitiesRes); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// @Summary		Update Amenity
// @Description	Update Amenity by ID
// @Tags			Amenities
// @Accept			json
// @Produce		json
// @Param			ID		path		int						true	"Amenity ID"
// @Param			Payload	body		UpdateAmenityPayload	true	"Payload update Amenity"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		404		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/amenities/{ID} [PUT]
func (app *application) UpdateAmenityHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UpdateAmenityPayload{}
	catId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(catId)
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

	amenity, err := app.repository.Amenities.GetByID(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if amenity.Name == payload.Name {
		alreadyExist := errors.New("can not update name to previous name")

		app.badRequestResponse(w, r, alreadyExist)
		return
	}

	amenity.Name = payload.Name

	// if change the type update new one
	if payload.TypeId != 0 {
		amenity.TypeId = payload.TypeId
	}

	if err := app.repository.Amenities.Update(ctx, amenity); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "update amenity successfull"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Delete Amenity
// @Description	Delete Amenity by ID
// @Tags			Amenities
// @Param			ID	path	int	true	"Amenity ID"
// @Success		204
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/amenities/{ID} [delete]
func (app *application) DeleteAmenityHandler(w http.ResponseWriter, r *http.Request) {
	amenityId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(amenityId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	amenity, err := app.repository.Amenities.GetByID(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.repository.Amenities.Delete(ctx, amenity.Id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.responseNoContent(w); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
