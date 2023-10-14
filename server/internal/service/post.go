package service

import (
	"context"
	"fmt"
	"server/internal/model"
)

const (
	userJoinMessage  = "<SayHi>"
	userLeaveMessage = "<SayBye>"
	postsLimit       = 49
	stockBotMessage  = "/stock="
)

type PostService struct {
	Repo model.PostRepo
}

// NewPostService builds a service and injects its dependencies
func NewPostService(repo model.PostRepo) *PostService {
	return &PostService{
		Repo: repo,
	}
}

// CreatePost inserts a new post into the database, except for the message to query a stock
func (s *PostService) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	if post.Message == userJoinMessage {
		post.Message = fmt.Sprintf("%s joined the chatroom!", post.User.Username)
	}

	if post.Message == userLeaveMessage {
		post.Message = fmt.Sprintf("%s left the chatroom!", post.User.Username)
	}

	return s.Repo.CreatePost(ctx, post)
}

// GetRecentPosts returns the last 50 posts from the database
func (s *PostService) GetRecentPosts(ctx context.Context) ([]*model.Post, error) {
	return s.Repo.GetRecentPosts(ctx, postsLimit)
}
