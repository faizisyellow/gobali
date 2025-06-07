package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

type CreateTypePayload struct {
	Name string `json:"name" validate:"required,min=5"`
}

type UpdateTypePayload struct {
	CreateTypePayload
}

type TypeResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// @Summary		Create Type
// @Description	Create Type
// @Tags			Types
// @Accept			json
// @Produce		json
// @Param			payload	body		CreateTypePayload	true	"json format payload"
//
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		400		{object}	main.WriteJSONError.envelope
// @Failure		409		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/types [POST]
func (app *application) CreateTypeHandler(w http.ResponseWriter, r *http.Request) {
	payload := &CreateTypePayload{}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.repository.Types.Create(r.Context(), payload.Name); err != nil {
		switch err {
		case repository.ErrDuplicateTypes:
			app.conflictErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "type created successfull"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// @Summary		Get Type
// @Description	Get Type by ID
// @Tags			Types
// @Produce		json
// @Param			ID	path		int	true	"Type id"
//
// @Success		200	{object}	main.jsonResponse.envelope{data=repository.Type}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/types/{ID} [GET]
func (app *application) GetTypeByIDHandler(w http.ResponseWriter, r *http.Request) {
	categoryId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(categoryId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ty, err := app.repository.Types.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusOK, ty); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Get types
// @Description	Get all types
// @Tags			Types
// @Produce		json
//
// @Success		200	{object}	main.jsonResponse.envelope{data=[]TypeResponse}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/types [GET]
func (app *application) GetTypesHandler(w http.ResponseWriter, r *http.Request) {

	tys, err := app.repository.Types.GetTypes(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	typeRes := []*TypeResponse{}

	for _, ty := range tys {
		newType := &TypeResponse{}

		newType.Id = ty.Id
		newType.Name = ty.Name
		newType.CreatedAt = ty.CreatedAt

		typeRes = append(typeRes, newType)
	}

	if err := app.jsonResponse(w, http.StatusOK, typeRes); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// @Summary		Update type
// @Description	Update type by ID
// @Tags			Types
// @Accept			json
// @Produce		json
// @Param			ID		path		int					true	"type ID"
// @Param			Payload	body		UpdateTypePayload	true	"Payload update type"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		404		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/types/{ID} [PUT]
func (app *application) UpdateTypeHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UpdateTypePayload{}
	typeId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(typeId)
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

	ty, err := app.repository.Types.GetByID(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if ty.Name == payload.Name {
		alreadyExist := errors.New("can not update name to previous name")

		app.badRequestResponse(w, r, alreadyExist)
		return
	}

	newType := &repository.Type{}
	newType.Id = id
	newType.Name = payload.Name

	if err := app.repository.Types.Update(ctx, newType); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "update type successfull"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Delete Type
// @Description	Delete Type by ID
// @Tags			Types
// @Param			ID	path	int	true	"Type ID"
// @Success		204
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/types/{ID} [delete]
func (app *application) DeleteTypeHandler(w http.ResponseWriter, r *http.Request) {
	typeId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(typeId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	ty, err := app.repository.Types.GetByID(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.repository.Types.Delete(ctx, ty.Id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.responseNoContent(w); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
