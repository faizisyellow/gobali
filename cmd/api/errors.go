package main

import (
	"net/http"

	"github.com/charmbracelet/log"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Error("internal server errors:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusInternalServerError, &[]string{"the server encountered a problem"})
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Warn("bad request errors:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusBadRequest, &[]string{err.Error()})
}

func (app *application) conflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Error("conflict errors:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusConflict, &[]string{err.Error()})
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Warn("not found errors:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusNotFound, &[]string{err.Error()})

}

func (app *application) unAuthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Warn("unauthorize", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusUnauthorized, &[]string{"unauthorize"})
}

func (app *application) unAuthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Warn("unauthorize Basic error", "path", r.URL, "method", r.Method, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	WriteJSONError(w, http.StatusUnauthorized, &[]string{"unauthorize"})
}

func (app *application) forbiddenErrorResponse(w http.ResponseWriter, r *http.Request) {
	log.Warn("forbidden access", "path", r.URL, "method", r.Method)

	WriteJSONError(w, http.StatusForbidden, &[]string{"forbidden"})
}
