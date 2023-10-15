package model

import (
	"context"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password,omitempty"`
	Posts    *Post     `pg:"rel:belongs-to" json:"posts,omitempty"`
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByName(ctx context.Context, user *User) (*User, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	LoginUser(ctx context.Context, user *User) *User
}
