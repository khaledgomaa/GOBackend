package store

import (
	"context"
	"database/sql"
)

// Storage acts as a centralized data access layer.
// Equivalent to DbContext in EF Core or a Unit of Work / Repository pattern wrapper in .NET.
type Storage struct {
	// Interfaces are defined where they are used in Go.
	// Similar to defining IPostRepository and IUserRepository in .NET.
	Posts interface {
		Create(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: NewPostsStore(db),
		Users: NewUsersStore(db),
	}
}
