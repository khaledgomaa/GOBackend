package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gobackend/social/internal/auth"
	"github.com/gobackend/social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/gobackend/social/docs"
	"go.uber.org/zap"
)

// application struct acts as a container for dependencies.
// Methods on this struct act as HTTP handlers (similar to Controller actions).
type application struct {
	config        config
	store         store.Storage
	logger        *zap.Logger
	authenticator *auth.Authenticator
}

type config struct {
	addr string
	db   dbConfig
	env  string
	auth authConfig
}

type authConfig struct {
	secret string
	iss    string
	aud    string
}

type dbConfig struct {
	conn         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// mount() configures the HTTP pipeline and routing.
// Equivalent to app.Use...() and app.MapControllers() in .NET Program.cs
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // Request ID middleware
	r.Use(middleware.RealIP)    // Real IP middleware
	r.Use(middleware.Recoverer) // Exception middleware
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Get("/swagger/*", httpSwagger.Handler())

		r.Post("/auth/token", app.generateTokenHandler)

		r.Group(func(r chi.Router) {
			r.Use(app.AuthMiddleware)

			r.Route("/posts", func(r chi.Router) {
				r.Post("/", app.createPostHandler)

				r.Route("/{postID}", func(r chi.Router) {
					r.Get("/", app.getPostHandler)
					r.Delete("/", app.deletePostHandler)
					r.Patch("/", app.updatePostHandler)
					r.Post("/comments", app.createCommentHandler)
				})
			})

			r.Route("/users", func(r chi.Router) {
				r.Route("/{userID}", func(r chi.Router) {
					r.Get("/", app.getUserHandler)
					r.Get("/posts", app.getUserPostsHandler)
				})
			})
		})

		// Public endpoint to register a user
		r.Post("/users", app.createUserHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute, // For IO
	}

	app.logger.Info("Server has started", zap.String("addr", app.config.addr))

	return srv.ListenAndServe()
}
