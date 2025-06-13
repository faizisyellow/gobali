package main

import (
	"context"
	"net/http"

	"github.com/faizisyellow/gobali/internal/uploader"
)

type filenamectxKey string

var (
	filenameKey filenamectxKey = "filenames"
)

func (app *application) UploadImagesMiddleware(next http.HandlerFunc, dst string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowContent := []string{"image/png", "image/jpeg"}

		var maxMem int64 = 3 * 1024 * 1024 // 3 mb

		filenames, err := app.upload.Upload(r, dst, maxMem, allowContent)

		if err != nil {
			switch err {
			case uploader.ErrExtNotAllowed:
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
