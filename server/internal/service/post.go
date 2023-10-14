package service

import (
	"context"
	"fmt"
	"server/internal/model"
)

const (
	userJoinMessage  = "<SayHi>"
	userLeaveMessage = "<SayBye>"
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

// CreatePost inserts a new post into the database and sends the updated post list to the broadcast channel
func (s *PostService) CreatePost(ctx context.Context, post *model.Post, broadcast chan []byte) error {
	if post.Message == userJoinMessage {
		post.Message = fmt.Sprintf("%s joined the chatroom!", post.User.Username)
	}

	if post.Message == userLeaveMessage {
		post.Message = fmt.Sprintf("%s left the chatroom!", post.User.Username)
	}

	post, err := s.Repo.CreatePost(ctx, post)
	if err != nil {
		return err
	}

	broadcastPosts(s.Repo, broadcast)

	return nil
}
