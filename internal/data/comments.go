package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/RayMC17/comments/internal/validator"
)

type Comment struct {
	//DB        *sql.DB
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

type CommentModel struct {
    DB *sql.DB
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	// check if the Content field is empty
	v.Check(comment.Content != "", "content", "must be provided")
	// check if the Author field is empty
	v.Check(comment.Author != "", "author", "must be provided")
	// check if the Content field is empty
	v.Check(len(comment.Content) <= 100, "content", "must not be more than 100 bytes long")
	// check if the Author field is empty
	v.Check(len(comment.Author) <= 25, "author", "must not be more than 25 bytes long")

}

// Expects a pointer to the actual comment
func (c CommentModel) Insert(comment *Comment) error {
	// the SQL query to be executed against the database table
	query := `
		 INSERT INTO comments (content, author)
		 VALUES ($1, $2)
		 RETURNING id, created_at, version
		 `
	// the actual values to replace $1, and $2
	args := []any{comment.Content, comment.Author}
	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// execute the query against the comments database table. We ask for the the
	// id, created_at, and version to be sent back to us which we will use
	// to update the Comment struct later on
	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt, &comment.Version)

}
