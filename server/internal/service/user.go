package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 13)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error hashing password: %s", err))
	}

	user.Password = string(hashedPassword)

	return s.Repo.CreateUser(ctx, user)
}

// LoginUser queries a user using username and password and returns it if found
func (s *UserService) LoginUser(ctx context.Context, user *model.User) *model.User {
	dbUser, err := s.Repo.GetUserByName(ctx, user)
	if err != nil {
		dbUser = &model.User{}
		log.Printf("error getting the user from the database: %s", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		dbUser = &model.User{}
		log.Printf("password does not match: %s", err)
	}

	return dbUser
}
