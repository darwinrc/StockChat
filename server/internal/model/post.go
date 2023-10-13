package model

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID        uuid.UUID  `json:"id"`
	UserID    string     `json:"userID"`
	User      *User      `json:"user" pg:"rel:has-one"`
	Message   string     `json:"message"`
	Timestamp *time.Time `json:"timestamp"`
}

type PostRepo interface {
	CreatePost(ctx context.Context, post *Post) (*Post, error)
	GetRecentPosts(ctx context.Context, limit int) ([]*Post, error)
}

type PostService interface {
	CreatePost(ctx context.Context, post *Post) (*Post, error)
	GetRecentPosts(ctx context.Context) ([]*Post, error)
}
