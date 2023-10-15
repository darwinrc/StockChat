package model

import (
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
