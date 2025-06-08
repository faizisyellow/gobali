package main

import (
	"net/http"
	"net/textproto"
)

func (app *application) UploadImagesMiddleware(next http.HandlerFunc, dst string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var allow = make(textproto.MIMEHeader, 0)
		allow.Add("images", "jpg")
		allow.Add("images", "jpeg")
		allow.Add("images", "png")

		var maxMem int64 = 3 * 1024 * 1024 // 3 mb

		app.upload.Upload(r, dst, &maxMem, allow)

		next.ServeHTTP(w, r)
	})
}
