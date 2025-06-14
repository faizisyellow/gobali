package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	Validate.RegisterValidation("withspace", func(fl validator.FieldLevel) bool {
		return !strings.Contains(fl.Field().String(), " ")
	})

	Validate.RegisterValidation("validpassword", func(fl validator.FieldLevel) bool {
		var hasUpper, hasDigit, hasSpecial bool

		for _, ch := range fl.Field().String() {
			switch {
			case unicode.IsUpper(ch):
				hasUpper = true
			case unicode.IsDigit(ch):
				hasDigit = true
			case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
				hasSpecial = true
			}
		}

		return hasUpper && hasDigit && hasSpecial
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	// limit of the body size for 1mb
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func readJsonMultiPartForm(r *http.Request, field string, data any) error {

	jsonField := r.MultipartForm.Value[field][0]

	decoder := json.NewDecoder(strings.NewReader(jsonField))
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, status int, message *[]string) error {
	type envelope struct {
		Errors []string `json:"errors"`
	}

	return writeJSON(w, status, &envelope{Errors: *message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}

func (app *application) responseNoContent(w http.ResponseWriter) error {

	w.WriteHeader(http.StatusNoContent)

	_, err := fmt.Fprint(w)
	if err != nil {
		return err
	}

	return nil
}
