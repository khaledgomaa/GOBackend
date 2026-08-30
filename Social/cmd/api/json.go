package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Global validator instance.
// Equivalent to injecting IValidator<T> from FluentValidation in .NET.
var validate *validator.Validate

// init() is called automatically by the Go runtime before main() executes.
// It acts exactly like a Static Constructor in C# (e.g., static ClassName() { ... }).
// We use it here to ensure our global validate instance is ready before the app starts.
func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// writeJSON acts like returning Results.Json() or Ok(data) in a .NET Controller.
// It sets the status code and serializes the object to the response stream.
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// json.Encoder is similar to System.Text.Json.JsonSerializer.SerializeAsync
	return json.NewEncoder(w).Encode(data) // Serialize
}

// readJSON is equivalent to reading from [FromBody] in .NET or using JsonSerializer.Deserialize.
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	// Similar to setting JsonSerializerOptions.UnmappedMemberHandling = JsonUnmappedMemberHandling.Disallow
	decoder.DisallowUnknownFields() // Restrict unknown fields

	return decoder.Decode(data) // Deserialize
}

// writeJSONError standardizes error responses.
// Equivalent to returning ProblemDetails or a custom error envelope in .NET.
func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}
