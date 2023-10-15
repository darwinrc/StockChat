package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password,omitempty"`
	Posts    *Post     `pg:"rel:belongs-to" json:"posts,omitempty"`
}
