// Equivalent to a Controller or Minimal API endpoint in .NET
// e.g., app.MapGet("/health", () => Results.Ok("Healthy"));
package main

import (
	"net/http"
)

// The (app *application) receiver makes this method part of the application struct,
// allowing it to access dependencies like a Controller instance method.
// healthCheckHandler godoc
//	@Summary		Health Check
//	@Description	Endpoint to verify the API is running
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    app.config.env,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
