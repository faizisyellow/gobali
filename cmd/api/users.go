package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/faizisyellow/gobali/internal/repository"
)

type CreateUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. read body and decoded

	// limit of the body size for 1mb
	maxBytes := 1_048_578
	userBody := r.Body
	userBody = http.MaxBytesReader(w, userBody, int64(maxBytes))

	userPayload := &CreateUserPayload{}

	decoded := json.NewDecoder(userBody)
	decoded.Decode(userPayload)

	// 2. create user
	user := &repository.User{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	}

	user.Password.Set(userPayload.Password)

	err := app.repository.Users.Create(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprintf(w, "internal server error %v", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "user created ")
}
