package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// Post represents a database entity and a JSON response object.
// Struct tags (`json:"id"`) are equivalent to attributes like [JsonPropertyName("id")] in .NET.
type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

// PostStore implements the data access logic for Posts.
// Equivalent to a Repository class in .NET (e.g., PostRepository : IPostRepository).
type PostStore struct {
	db *sql.DB // Equivalent to injecting DbContext or IDbConnection
}

func NewPostsStore(db *sql.DB) *PostStore {
	return &PostStore{db: db}
}

// Create executes an INSERT query.
// Equivalent to dbContext.Posts.Add(post); dbContext.SaveChanges(); or using Dapper.
func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT into posts (content,title,user_id,tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
	SELECT id, user_id,title,content,created_at,updated_at,tags
	FROM posts
	WHERE id = $1
	`

	var post Post

	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	comments, err := NewCommentStore(s.db).GetByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	post.Comments = comments

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content = $2, tags = $3
	WHERE id = $4
	RETURNING updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
	).Scan(&post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
