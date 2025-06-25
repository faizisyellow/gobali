package main

import (
	"context"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/faizisyellow/gobali/internal/helpers"
	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

type villaCtxKey string

const villaKey villaCtxKey = "villa"

type CreateVillaProp struct {
	Name        string  `json:"name" validate:"required,min=4"`
	Description string  `json:"description" validate:"required,min=8"`
	MinGuest    int     `json:"min_guest" validate:"required,min=1"`
	Bedrooms    int     `json:"bedrooms" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,min=1"`
	Baths       int     `json:"baths" validate:"required,min=1"`
	AmenityId   []int   `json:"amenity_id"`
	LocationId  int     `json:"location_id"`
	CategoryId  int     `json:"category_id"`
}

type UpdateVillaPayload struct {
	Name        *string  `json:"name" `
	Description *string  `json:"description"`
	MinGuest    *int     `json:"min_guest" `
	Bedrooms    *int     `json:"bedrooms"`
	Price       *float64 `json:"price" `
	Baths       *int     `json:"baths"`
	LocationId  *int     `json:"location_id"`
	CategoryId  *int     `json:"category_id"`
}

func (u *UpdateVillaPayload) Apply(villa *repository.Villa) {
	if u.Name != nil {
		villa.Name = *u.Name
	}

	if u.Baths != nil {
		villa.Baths = *u.Baths
	}

	if u.Bedrooms != nil {
		villa.Bedrooms = *u.Bedrooms
	}

	if u.CategoryId != nil {
		villa.CategoryId = *u.CategoryId
	}

	if u.Description != nil {
		villa.Description = *u.Description
	}

	if u.LocationId != nil {
		villa.LocationId = *u.LocationId
	}

	if u.MinGuest != nil {
		villa.MinGuest = *u.MinGuest
	}

	if u.Price != nil {
		villa.Price = *u.Price
	}
}

// @Summary		Create Villa
// @Description	Create Villa
// @Tags			Villas
// @Produce		json
// @Accept			mpfd
// @Param			thumbnail	formData	file	true	"Image file"
// @Param			others		formData	file	false	"Image file"
// @Param			properties	formData	string	true	"CreateVillaProp JSON string"	example({"name":"villa name","description":"villa description","min_guest":1,"bedrooms":1,"price":25,"location_id":3,"category_id":2,"baths":1,"amenity_id":[4]})
// @Security		JWT
// @Success		201	{object}	main.jsonResponse.envelope{data=string}
// @Success		400	{object}	main.WriteJSONError.envelope
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/villas [post]
func (app *application) CreateVillaHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	images := ctx.Value(filenameKey).([]string)
	payload := &CreateVillaProp{}

	if err := readJsonMultiPartForm(r, "properties", payload); err != nil {
		for _, image := range images {
			helpers.RemoveFile(filepath.Join(app.configs.upload.baseDir, "villas", image))
		}

		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		for _, image := range images {
			helpers.RemoveFile(filepath.Join(app.configs.upload.baseDir, "villas", image))
		}

		app.badRequestResponse(w, r, err)
		return
	}

	var amenity = []repository.SelectedAmenity{}

	for _, id := range payload.AmenityId {
		amenity = append(amenity, repository.SelectedAmenity{Id: id})
	}

	newVilla := &repository.Villa{
		Name:        payload.Name,
		Description: payload.Description,
		MinGuest:    payload.MinGuest,
		Bedrooms:    payload.Bedrooms,
		Baths:       payload.Baths,
		Price:       payload.Price,
		ImageUrls:   images,
		CategoryId:  payload.CategoryId,
		LocationId:  payload.LocationId,
		Amenity:     amenity,
	}

	err := app.repository.Villas.CreateVillaWithAmenity(ctx, newVilla)
	if err != nil {
		for _, image := range images {
			helpers.RemoveFile(filepath.Join(app.configs.upload.baseDir, "villas", image))
		}

		switch err {
		case repository.ErrCatOrLocNotExist, repository.ErrDuplicateAmenities, repository.ErrAmenitiesNotExist, repository.ErrDuplicateVilla:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}

		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "villa created successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Update Villa
// @Description	Update Villa by ID
// @Tags			Villas
// @Produce		json
// @Param			ID	path	int	true	"Villa ID"
// @Accept			mpfd
// @Param			thumbnail	formData	file	false	"Image file"
// @Param			others		formData	file	false	"Image file"
// @Param			properties	formData	string	false	"Update Villa Props JSON string"	example({"name":"villa name","description":"villa description","min_guest":1,"bedrooms":1,"price":25,"location_id":3,"category_id":2,"baths":1})
// @Security		JWT
// @Success		201	{object}	main.jsonResponse.envelope{data=string}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/villas/{ID} [PUT]
func (app *application) UpdateVillaHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UpdateVillaPayload{}

	ctx := r.Context()

	images := ctx.Value(filenameKey)

	var imagesUpdated []string

	if images != nil {
		v, _ := images.([]string)
		imagesUpdated = v
	}

	if err := readJsonMultiPartForm(r, "properties", payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	villa := GetVillaFromContext(r)

	payload.Apply(villa)

	if images != nil {
		villa.ImageUrls = imagesUpdated
	}

	if err := app.repository.Villas.Update(ctx, villa); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "updated villa successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// @Summary		Get Villa
// @Description	Get Villa By ID
// @Tags			Villas
// @Produce		json
// @Param			villaID	path		int	true	"Villa ID"
// @Success		200		{object}	main.jsonResponse.envelope{data=repository.Villa}
// @Failure		404		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/villas/{villaID} [get]
func (app *application) GetVillaByIdHandler(w http.ResponseWriter, r *http.Request) {

	villa := GetVillaFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, villa); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Get Villa
// @Description	Get All Villas
// @Tags			Villas
// @Produce		json
// @Param			limit		query		string	false	"limit each page"
// @Param			offset		query		string	false	"skip rows"
// @Param			sort		query		string	false	"sort villa latest(desc), older(asc)"
// @Param			location	query		string	false	"location villa"
// @Param			category	query		string	false	"category villa"
// @Param			bedrooms	query		string	false	"bedrooms villa"
// @Param			min_guest	query		string	false	"min guest villa"
// @Success		200			{object}	main.jsonResponse.envelope{data=[]repository.Villa}
// @Failure		500			{object}	main.WriteJSONError.envelope
// @Router			/villas [get]
func (app *application) GetVillasHandler(w http.ResponseWriter, r *http.Request) {

	vq, err := repository.PaginatedVillaQuery{
		Limit:    5,
		Offset:   0,
		Sort:     "asc",
		Location: "",
		Category: "",
		MinGuest: "",
		Bedrooms: "",
	}.Parse(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(vq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	villas, err := app.repository.Villas.GetVillas(r.Context(), vq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, villas); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Delete Villa
// @Description	Delete Villa By ID
// @Tags			Villas
// @Produce		json
// @Param			villaID	path	int	true	"Villa ID"
// @Security		JWT
// @Success		204
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/villas/{villaID} [delete]
func (app *application) DeleteVillaByIdHandler(w http.ResponseWriter, r *http.Request) {

	villa := GetVillaFromContext(r)

	ctx := r.Context()
	err := app.repository.Villas.Delete(ctx, villa.Id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.responseNoContent(w)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) VillaContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		villaId := chi.URLParam(r, "villaID")

		id, err := strconv.Atoi(villaId)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		villa, err := app.repository.Villas.GetById(ctx, id)
		if err != nil {
			switch err {
			case repository.ErrNoRows:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, villaKey, villa)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetVillaFromContext(r *http.Request) *repository.Villa {
	villa := r.Context().Value(villaKey).(*repository.Villa)

	return villa
}
