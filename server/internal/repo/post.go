package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"server/db"
	"server/internal/model"
)

type PostRepo interface {
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	GetRecentPosts(ctx context.Context, limit int) ([]*model.Post, error)
}

type postRepository struct {
	db db.DB
}

// NewPostRepository builds a postRepository and injects its dependencies
func NewPostRepository(db db.DB) PostRepo {
	return &postRepository{db: db}
}

// CreatePost insert a new post into the database
func (r *postRepository) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	var lastInsertId uuid.UUID
	query := `INSERT INTO posts(user_id, message) VALUES ($1, $2) returning (id)`

	err := r.db.QueryRowContext(ctx, query, post.UserID, post.Message).Scan(&lastInsertId)
	if err != nil {
		return &model.Post{}, err
	}

	post.ID = lastInsertId
	return post, nil
}

// GetRecentPosts returns the last <limit> posts from the database, including the associated user data
func (r *postRepository) GetRecentPosts(ctx context.Context, limit int) ([]*model.Post, error) {
	query := `
		SELECT posts.id, posts.user_id, posts.message, posts.timestamp, users.id, users.username 
		FROM posts
		INNER JOIN users ON users.id = posts.user_id
		ORDER BY posts.timestamp DESC LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return []*model.Post{}, nil
	}
	defer rows.Close()

	var posts []*model.Post

	for rows.Next() {
		post := &model.Post{
			User: &model.User{},
		}
		if err := rows.Scan(&post.ID, &post.UserID, &post.Message, &post.Timestamp, &post.User.ID, &post.User.Username); err != nil {
			return nil, errors.New(fmt.Sprintf("error scanning rows: %s", err))
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error during rows iteration: %s", err))
	}

	return posts, nil
}
