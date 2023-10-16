package repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"server/internal/model"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostRepository(db)

	postID, _ := uuid.FromBytes([]byte("cfab745c-25d2-4a48-a94c-d3f84ef9167a"))

	mock.ExpectQuery("INSERT INTO posts").
		WithArgs("48ccb5c1-9a19-42cd-bd41-3ac5c8af1108", "Test Message").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(postID))

	post := &model.Post{
		UserID:  "48ccb5c1-9a19-42cd-bd41-3ac5c8af1108",
		Message: "Test Message",
	}

	createdPost, err := repo.CreatePost(context.Background(), post)

	assert.NoError(t, err)
	assert.NotNil(t, createdPost)
	assert.Equal(t, postID, createdPost.ID)
}

func TestGetRecentPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostRepository(db)

	postID, _ := uuid.FromBytes([]byte("cfab745c-25d2-4a48-a94c-d3f84ef9167a"))
	userID, _ := uuid.FromBytes([]byte("48ccb5c1-9a19-42cd-bd41-3ac5c8af1108"))

	mock.ExpectQuery("SELECT posts.id, posts.user_id, posts.message, posts.timestamp, users.id, users.username").WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "message", "timestamp", "users_id", "username"}).
			AddRow(postID, "48ccb5c1-9a19-42cd-bd41-3ac5c8af1108", "Test Message", time.Now(), userID, "Alice"))

	limit := 5
	recentPosts, err := repo.GetRecentPosts(context.Background(), limit)

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, recentPosts)
	assert.Len(t, recentPosts, 1)
}
