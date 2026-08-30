package db

import (
	"context"
	"database/sql"
	"time"
)

// New creates and configures a database connection pool.
// Equivalent to services.AddDbContextPool<T>() in .NET.
func New(conn string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	// sql.Open doesn't actually connect, it just validates arguments.
	// Ping() is used to verify the connection, similar to context.Database.CanConnect() in EF Core.
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
