package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/faizisyellow/gobali/internal/uploader"
	"github.com/golang-jwt/jwt/v5"
)

type filenamectxKey string
type userKey string

var (
	filenameKey filenamectxKey = "filenames"
	userCtx     userKey        = "user"
)

func (app *application) UploadImagesMiddleware(next http.HandlerFunc, dst string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowMime := []string{"image/png", "image/jpeg"}

		var maxMem int64 = 3 * 1024 * 1024 // 3 mb

		r.ParseMultipartForm(maxMem)

		fileFields := r.MultipartForm.File

		if len(fileFields) == 0 {
			app.badRequestResponse(w, r, fmt.Errorf("image file required"))
			return
		}

		filenames, err := app.upload.Upload(r, dst, maxMem, allowMime)
		if err != nil {
			switch err {
			case uploader.ErrExtNotAllowed, uploader.ErrSizeTooLarger:
				app.badRequestResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}

			return
		}

		ctx := context.WithValue(r.Context(), filenameKey, filenames)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unAuthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		// signed token
		token := parts[1]

		// validate token and decode the token
		jwtToken, err := app.authentication.VerifyToken(token)
		if err != nil {
			app.unAuthorizedErrorResponse(w, r, err)
			return
		}

		// claim's token
		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["id"]), 10, 64)
		if err != nil {
			app.unAuthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.repository.Users.GetByID(ctx, int(userID))
		if err != nil {
			app.unAuthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) AuthBasicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unAuthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			app.unAuthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			app.unAuthorizedBasicErrorResponse(w, r, err)
			return
		}

		username := app.configs.auth.basic.username
		password := app.configs.auth.basic.password

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			app.unAuthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) OfficerOnlyAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := getUserFromContext(r)

		granted, err := app.CheckRolePresedence(ctx, user, "admin")
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !granted {
			app.forbiddenErrorResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) UserAction(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)

		if user.Role.Name != "user" {
			app.forbiddenErrorResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) BookingAccess(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		booking := GetBookingFromContext(r)
		user := getUserFromContext(r)

		// if the booking's is the user then access it, but can not remove it
		if booking != nil && r.Method != "DELETE" && user.Id == booking.UserId {
			next.ServeHTTP(w, r)
			return
		}

		granted, err := app.CheckRolePresedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !granted {
			app.forbiddenErrorResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) CheckRolePresedence(ctx context.Context, user *repository.User, rolename string) (bool, error) {
	role, err := app.repository.Roles.GetByName(ctx, rolename)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
