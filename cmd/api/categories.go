package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

type CreateCategoryPayload struct {
	Name string `json:"name" validate:"required,min=5"`
}

type UpdateCategoryPayload struct {
	CreateCategoryPayload
}

type CategoryResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

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

func (app *application) GetCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	categoryId := chi.URLParam(r, "ID")

	id, err := strconv.Atoi(categoryId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	cat, err := app.repository.Categories.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			app.notFoundResponse(w, r, err)
		}

		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, cat); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	cats, err := app.repository.Categories.GetCategories(r.Context())
	if err != nil {
		switch err {
		case repository.ErrDuplicateCategory:
			if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
				app.internalServerError(w, r, err)
				return
			}
		default:
			app.internalServerError(w, r, err)
		}

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
		if errors.Is(err, repository.ErrNoRows) {
			app.notFoundResponse(w, r, err)
		}

		app.internalServerError(w, r, err)
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
		if errors.Is(err, repository.ErrNoRows) {
			app.notFoundResponse(w, r, err)
		}

		app.internalServerError(w, r, err)
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
