// Equivalent to a Controller or Minimal API endpoint in .NET
// e.g., app.MapGet("/health", () => Results.Ok("Healthy"));
package main

import (
	"net/http"
)

// The (app *application) receiver makes this method part of the application struct,
// allowing it to access dependencies like a Controller instance method.
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    app.config.env,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
