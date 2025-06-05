package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "ping")
}
