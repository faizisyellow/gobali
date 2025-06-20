package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
