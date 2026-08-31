// Equivalent to Program.cs in .NET (Entry point, configuration, and DI setup)
package main

import (
	"github.com/gobackend/social/internal/auth"
	"github.com/gobackend/social/internal/db"
	"github.com/gobackend/social/internal/env"
	"github.com/gobackend/social/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//	@title			Social API
//	@description	A social network API.
//	@version		1.0
//	@host			localhost:4040
//	@BasePath		/v1
func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	// Loads .env file (Similar to reading appsettings.json or Environment Variables in .NET)
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}

	// Equivalent to IOptions<T> or binding Configuration sections in .NET
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			conn:         env.GetString("DB_CONN", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "DEVELOPMENT"),
		auth: authConfig{
			secret: env.GetString("AUTH_SECRET", "supersecret"),
			iss:    env.GetString("AUTH_ISS", "social-api"),
			aud:    env.GetString("AUTH_AUD", "social-client"),
		},
	}

	db, err := db.New(cfg.db.conn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)

	authenticator := auth.NewAuthenticator(cfg.auth.secret, cfg.auth.iss, cfg.auth.aud)

	// Go doesn't have a built-in IoC container like IServiceCollection.
	// Dependencies are explicitly passed via structs (Constructor Injection).
	app := &application{
		config:        cfg,
		store:         store.NewStorage(db),
		logger:        logger,
		authenticator: authenticator,
	}

	defer db.Close()

	logger.Info("database connection pool established")

	mux := app.mount() // Middleware

	logger.Fatal("Server stopped", zap.Error(app.run(mux)))
}
