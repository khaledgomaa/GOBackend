// This file equivalent to program.cs in .NET
package main

import (
	"log"

	"github.com/gobackend/social/internal/db"
	"github.com/gobackend/social/internal/env"
	"github.com/gobackend/social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // To load environment variables defined insided .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// To initialize required configuration
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			conn:         env.GetString("DB_CONN", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.conn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)

	// Encapsulates required dependencies with respecting DI
	app := &application{
		config: cfg,
		store:  store.NewStorage(db),
	}

	defer db.Close()

	log.Println("database connection pool established")

	mux := app.mount() // Middleware

	log.Fatal(app.run(mux))
}
