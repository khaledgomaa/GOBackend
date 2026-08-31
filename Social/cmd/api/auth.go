package main

import (
	"net/http"
)

type GenerateTokenPayload struct {
	Email string `json:"email" validate:"required,email"`
}

// generateTokenHandler godoc
//
//	@Summary		Generate JWT token
//	@Description	Generate a JWT token for testing authentication
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		GenerateTokenPayload	true	"Token payload"
//	@Success		200		{object}	map[string]string
//	@Router			/auth/token [post]
func (app *application) generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload GenerateTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate user email already exists
	_, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	// Just for testing, generating a token directly without password validation
	token, err := app.authenticator.GenerateToken(payload.Email)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, map[string]string{"token": token}); err != nil {
		app.internalServerError(w, r, err)
	}
}
