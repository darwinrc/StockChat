package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"server/db"
	"server/internal/model"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByName(ctx context.Context, user *model.User) (*model.User, error)
}

type userRepository struct {
	db db.DB
}

// NewUserRepository builds a userRepository and injects its dependencies
func NewUserRepository(db db.DB) UserRepo {
	return &userRepository{db: db}
}

// CreateUser insert a new user into the database
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	var lastInsertId uuid.UUID
	query := `INSERT INTO users(username, password) VALUES ($1, $2) returning (id)`

	if err := r.db.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&lastInsertId); err != nil {
		return &model.User{}, err
	}

	user.ID = lastInsertId
	return user, nil
}

// GetUserByName searches for a user in the database given the username
func (r *userRepository) GetUserByName(ctx context.Context, user *model.User) (*model.User, error) {
	dbUser := &model.User{}
	query := `SELECT id, username, password FROM users WHERE username = $1`

	if err := r.db.QueryRowContext(ctx, query, user.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password); err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, errors.New(fmt.Sprintf("error querying the users table: %s", err))
		}
	}

	return dbUser, nil
}
