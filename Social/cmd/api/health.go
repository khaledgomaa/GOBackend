// Equivalent to a Controller or Minimal API endpoint in .NET
// e.g., app.MapGet("/health", () => Results.Ok("Healthy"));
package main

import "net/http"

// The (app *application) receiver makes this method part of the application struct,
// allowing it to access dependencies like a Controller instance method.
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy"))
	w.WriteHeader(http.StatusOK)
}
