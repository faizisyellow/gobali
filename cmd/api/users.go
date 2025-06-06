package main

import (
	"net/http"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
)

type CreateUserPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,max=21,email"`
	Password string `json:"password" validate:"required,max=12,withspace,validpassword"`
}

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userPayload := &CreateUserPayload{}

	err := readJSON(w, r, userPayload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(userPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &repository.User{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	}

	user.Password.Set(userPayload.Password)

	err = app.repository.Users.Create(r.Context(), user)
	if err != nil {

		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "user created successfuly"); err != nil {

		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	inviteToken := chi.URLParam(r, "token")

	err := app.repository.Users.Activate(r.Context(), inviteToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "user activated successfuly"); err != nil {

		app.internalServerError(w, r, err)
		return
	}
}
