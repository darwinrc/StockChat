package service

import (
	"context"
	"server/internal/model"
)

type UserService struct {
	Repo model.UserRepo
}

// NewUserService builds a service and injects its dependencies
func NewUserService(repo model.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

// CreateUser inserts a new user into the database
func (s *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return s.Repo.CreateUser(ctx, user)
}

// LoginUser queries a user using username and password and returns it if found
func (s *UserService) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {
	return s.Repo.GetUserByCredentials(ctx, user)
}
