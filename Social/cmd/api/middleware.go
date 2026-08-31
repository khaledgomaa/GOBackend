package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type userKey string

const userContextKey = userKey("user")

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		tokenString := parts[1]
		claims, err := app.authenticator.VerifyToken(tokenString)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("email claim not found"))
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByEmail(ctx, email)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
