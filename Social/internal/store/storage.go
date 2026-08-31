package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("Resource not found")
	QueryTimeoutDuration = time.Second * 5
)

// Storage acts as a centralized data access layer.
// Equivalent to DbContext in EF Core or a Unit of Work / Repository pattern wrapper in .NET.
type Storage struct {
	// Interfaces are defined where they are used in Go.
	// Similar to defining IPostRepository and IUserRepository in .NET.
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    NewPostsStore(db),
		Users:    NewUsersStore(db),
		Comments: NewCommentStore(db),
	}
}
