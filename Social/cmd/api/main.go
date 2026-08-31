// Equivalent to Program.cs in .NET (Entry point, configuration, and DI setup)
package main

import (
	"log"

	"github.com/gobackend/social/internal/db"
	"github.com/gobackend/social/internal/env"
	"github.com/gobackend/social/internal/store"
	"github.com/joho/godotenv"
)

//	@title			Social API
//	@description	A social network API.
//	@version		1.0
//	@host			localhost:4040
//	@BasePath		/v1
func main() {
	// Loads .env file (Similar to reading appsettings.json or Environment Variables in .NET)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
	}

	db, err := db.New(cfg.db.conn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)

	// Go doesn't have a built-in IoC container like IServiceCollection.
	// Dependencies are explicitly passed via structs (Constructor Injection).
	app := &application{
		config: cfg,
		store:  store.NewStorage(db),
	}

	defer db.Close()

	log.Println("database connection pool established")

	mux := app.mount() // Middleware

	log.Fatal(app.run(mux))
}
