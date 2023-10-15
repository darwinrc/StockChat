package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/model"
	"server/internal/repo"
	"sort"
)

const (
	userJoinMessage  = "<SayHi>"
	userLeaveMessage = "<SayBye>"
	postsLimit       = 50
)

type PostService struct {
	Repo repo.PostRepo
}

// NewPostService builds a service and injects its dependencies
func NewPostService(repo repo.PostRepo) *PostService {
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

// broadcastPosts sends the merged list of posts + commands to the broadcast channel
func broadcastPosts(repo repo.PostRepo, broadcast chan []byte) {
	posts, err := repo.GetRecentPosts(context.Background(), postsLimit)
	if err != nil {
		log.Printf("error getting posts from database: %s", err)
	}

	posts = append(posts, commands...)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp.After(*posts[j].Timestamp)
	})

	bPosts, err := json.Marshal(posts)
	if err != nil {
		log.Printf("error marshaling posts: %s", err)
		return
	}

	broadcast <- bPosts
}
