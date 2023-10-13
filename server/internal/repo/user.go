package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"server/internal/model"
)

type UserRepo struct {
	DB *sql.DB
}

type userRepository struct {
	db Database
}

// NewUserRepository builds a userRepository and injects its dependencies
func NewUserRepository(db Database) model.UserRepo {
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

// GetUserByCredentials searches for a user in the database given the username and password
func (r *userRepository) GetUserByCredentials(ctx context.Context, user *model.User) (*model.User, error) {
	dbUser := &model.User{}
	query := `SELECT id, username FROM users WHERE username = $1 AND password = $2`

	if err := r.db.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&dbUser.ID, &dbUser.Username); err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, errors.New(fmt.Sprintf("error querying the users table: %s", err))
		}
	}

	return dbUser, nil
}
