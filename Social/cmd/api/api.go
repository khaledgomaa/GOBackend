package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gobackend/social/internal/store"
)

// application struct acts as a container for dependencies.
// Methods on this struct act as HTTP handlers (similar to Controller actions).
type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
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

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
