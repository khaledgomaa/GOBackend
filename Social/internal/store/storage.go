package store

import (
	"context"
	"database/sql"
)

type Storage struct { // equivalent to dbcontext in EF core
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
