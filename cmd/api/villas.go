package main

import (
	"net/http"
	"strconv"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

// TODO: handle create villa
func (app *application) CreateVillaHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.jsonResponse(w, http.StatusCreated, "OK"); err != nil {
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
	villaId := chi.URLParam(r, "villaID")

	id, err := strconv.Atoi(villaId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	villa, err := app.repository.Villas.GetById(r.Context(), id)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusOK, villa); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Get Villa
// @Description	Get All Villas
// @Tags			Villas
// @Produce		json
// @Success		200	{object}	main.jsonResponse.envelope{data=[]repository.Villa}
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/villas [get]
func (app *application) GetVillasHandler(w http.ResponseWriter, r *http.Request) {

	villas, err := app.repository.Villas.GetVillas(r.Context())
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
// @Success		204
// @Failure		404	{object}	main.WriteJSONError.envelope
// @Failure		500	{object}	main.WriteJSONError.envelope
// @Router			/villas/{villaID} [delete]
func (app *application) DeleteVillaByIdHandler(w http.ResponseWriter, r *http.Request) {
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

	err = app.repository.Villas.Delete(ctx, villa.Id)
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
