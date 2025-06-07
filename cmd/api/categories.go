package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

type CreateCategoryPayload struct {
	Name string `json:"name" validate:"required,min=3"`
}

type UpdateCategoryPayload struct {
	CreateCategoryPayload
}

type CategoryResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// @Summary		Create category
// @Description	Create category
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Param			payload	body		CreateCategoryPayload	true	"json format payload"
//
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		400		{object}	main.WriteJSONError.envelope
// @Failure		409		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/categories [POST]
func (app *application) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	payload := &CreateCategoryPayload{}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.repository.Categories.Create(r.Context(), payload.Name); err != nil {
		switch err {
		case repository.ErrDuplicateCategory:
			app.conflictErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "category created successfull"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// @Summary		Get category
// @Description	Get category by ID
// @Tags			Categories
// @Produce		json
// @Param			ID	path		int	true	"category id"
//
// @Success		200	{object}	main.jsonResponse.envelope{data=repository.Category}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/categories/{ID} [GET]
func (app *application) GetCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	categoryId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(categoryId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	cat, err := app.repository.Categories.GetByID(r.Context(), id)
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

// @Summary		Get categories
// @Description	Get all categories
// @Tags			Categories
// @Produce		json
//
// @Success		200	{object}	main.jsonResponse.envelope{data=[]CategoryResponse}
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/categories [GET]
func (app *application) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	cats, err := app.repository.Categories.GetCategories(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)

		return
	}

	catsRes := []*CategoryResponse{}

	for _, cat := range cats {
		newCat := &CategoryResponse{}

		newCat.Id = cat.Id
		newCat.Name = cat.Name
		newCat.CreatedAt = cat.CreatedAt

		catsRes = append(catsRes, newCat)
	}

	if err := app.jsonResponse(w, http.StatusOK, catsRes); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// @Summary		Update Category
// @Description	Update Category by ID
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Param			ID		path		int						true	"Category ID"
// @Param			Payload	body		UpdateCategoryPayload	true	"Payload update category"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		404		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/categories/{ID} [PUT]
func (app *application) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UpdateCategoryPayload{}
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

	cat, err := app.repository.Categories.GetByID(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if cat.Name == payload.Name {
		alreadyExist := errors.New("can not update name to previous name")

		app.badRequestResponse(w, r, alreadyExist)
		return
	}

	newCat := &repository.Category{}
	newCat.Id = id
	newCat.Name = payload.Name

	if err := app.repository.Categories.Update(ctx, newCat); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "update category successfull"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Delete Category
// @Description	Delete Category by ID
// @Tags			Categories
// @Param			ID	path	int	true	"Category ID"
// @Success		204
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/categories/{ID} [delete]
func (app *application) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	catId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(catId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	cat, err := app.repository.Categories.GetByID(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.repository.Categories.Delete(ctx, cat.Id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.responseNoContent(w); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
