package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/faizisyellow/gobali/internal/mailer"
	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,min=5,email"`
	Password string `json:"password" validate:"required,min=5,withspace,validpassword"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// @Summary		Register user
// @Description	Register new user
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			Payload	body		RegisterPayload	true	"Payload register new user"
// @Success		201		{object}	main.jsonResponse.envelope{data=string}
// @Failure		404		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/authentication/register [POST]
func (app *application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &repository.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Token invitation
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	ctx := r.Context()

	err := app.repository.Users.CreateAndInvite(ctx, user, hashToken, app.configs.mail.exp)
	if err != nil {
		switch err {
		case repository.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	isDevEnv := app.configs.env == "Development"

	// the links is from the frontend router (http://localhost:5173/confirm/{plaintoken})
	activationUrl := fmt.Sprintf("%s/confirm/%s", app.configs.clientURL, plainToken)

	vars := struct {
		Username      string
		ActivationUrl string
	}{
		Username:      user.Username,
		ActivationUrl: activationUrl,
	}

	// send email
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, isDevEnv)
	if err != nil {
		log.Error("error sending welcome email", "error", err.Error())

		// rollback user creation if email fails (SAGA pattern)
		if err := app.repository.Users.Delete(ctx, user.Id); err != nil {
			log.Error("error deleting user while rollback", "error", err.Error())
		}

		app.internalServerError(w, r, err)
		return
	}

	log.Info("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusCreated, plainToken); err != nil {
		app.internalServerError(w, r, err)

		return
	}

}

// @Summary		Login user
// @Description	Login user
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			Payload	body		LoginPayload	true	"Payload credential user, password: Tester_1234"
// @Success		200		{object}	main.jsonResponse.envelope{data=string}
// @Failure		400		{object}	main.WriteJSONError.envelope
// @Failure		401		{object}	main.WriteJSONError.envelope
// @Failure		500		{object}	main.WriteJSONError.envelope
// @Router			/authentication/login [POST]
func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	payload := &LoginPayload{}

	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.repository.Users.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			app.unAuthorizedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = user.Password.Compare(payload.Password)
	if err != nil {
		app.unAuthorizedErrorResponse(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"iss": app.configs.auth.token.iss,
		"sub": app.configs.auth.token.sub,
		"exp": app.configs.auth.token.exp,
		"id":  user.Id,
	}

	signedToken, err := app.authentication.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, signedToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
