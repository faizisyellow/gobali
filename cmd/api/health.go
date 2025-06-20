package main

import (
	"net/http"
)

//	@Summary		Health
//	@Description	Check Response
//	@Tags			Health
//
//	@Security		JWT
//	@Success		200	{object}	main.jsonResponse.envelope{data=string}
//	@Failure		500	{object}	main.WriteJSONError.envelope
//	@Router			/health [GET]
func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.jsonResponse(w, http.StatusOK, "ping"); err != nil {
		app.internalServerError(w, r, err)
	}
}
