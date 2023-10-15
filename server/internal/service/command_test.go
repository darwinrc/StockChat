package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"log"
	"os"
	mock_repo "server/internal/repo/mocks"
	"server/internal/service/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock for PostRepo
//type mockPostRepo struct{}

//func (m *mockPostRepo) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
//	userID := "f1c21d1d-3411-4bfd-a99f-8fc52dc65bb5"
//	uID, _ := uuid.FromBytes([]byte(userID))
//	ts, _ := time.Parse(time.RFC3339, "2023-10-15 06:07:38.467471")
//
//	return &model.Post{
//		ID:     uuid.UUID{},
//		UserID: userID,
//		User: &model.User{
//			ID:       uID,
//			Username: "Alice",
//		},
//		Message:   "test post",
//		Timestamp: &ts,
//	}, nil
//}
//
//func (m *mockPostRepo) GetRecentPosts(ctx context.Context, limit int) ([]*model.Post, error) {
//	userID := "f1c21d1d-3411-4bfd-a99f-8fc52dc65bb5"
//	uID, _ := uuid.FromBytes([]byte(userID))
//	ts, _ := time.Parse(time.RFC3339, "2023-10-15 06:07:38.467471")
//
//	userID2 := "57a28a21-989b-41de-a200-e2dd3e330a26"
//	uID2, _ := uuid.FromBytes([]byte(userID))
//	ts2, _ := time.Parse(time.RFC3339, "2023-10-15 06:17:38.467471")
//
//	return []*model.Post{
//		{
//			ID:     uuid.UUID{},
//			UserID: userID,
//			User: &model.User{
//				ID:       uID,
//				Username: "Alice",
//			},
//			Message:   "test post 1",
//			Timestamp: &ts,
//		},
//		{
//			ID:     uuid.UUID{},
//			UserID: userID2,
//			User: &model.User{
//				ID:       uID2,
//				Username: "Bob",
//			},
//			Message:   "test post 2",
//			Timestamp: &ts2,
//		},
//	}, nil
//}

func TestProcessCommand(t *testing.T) {
	t.Setenv("RABBITMQ_USERNAME", "guest")
	t.Setenv("RABBITMQ_PASSWORD", "guest")
	t.Setenv("RABBITMQ_HOST", "localhost:5672")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new CommandService with mock PostRepo
	service := &CommandService{
		PostRepo: &mocks.MockPostRepo{},
	}

	// Test case 1: Valid command
	command := "/stock=aapl.us"
	broadcast := make(chan []byte)
	scanner, reader, writer := mockLogger(t)
	defer resetLogger(reader, writer)

	go service.ProcessCommand(command, broadcast)

	scanner.Scan() // scans first log: Processing command ...
	scanner.Scan() // scans second log: Stock sent ...

	got := scanner.Text() // the last line written to the scanner
	msg := "Stock sent: {\"stockCode\":\"aapl.us\"}"
	assert.Contains(t, got, msg)

	//assert.ElementsMatch(t, expectedMessages, receivedMessages)

	// Test case 2: Error case - JSON marshal failure
	//command = "invalid-command"
	//broadcast = make(chan []byte)
	//service.ProcessCommand(command, broadcast)
	// You can assert expected error conditions here.
}

func TestBroadcastCommand(t *testing.T) {
	service := &CommandService{
		PostRepo: &mock_repo.MockPostRepo{},
	}

	// Create a channel for broadcasting messages
	broadcast := make(chan []byte)

	// Run the BroadcastCommand function
	go service.BroadcastCommand(broadcast)

	// Simulate the passage of time for testing
	time.Sleep(100 * time.Millisecond)

	// Simulate closing the channel when the BroadcastCommand function completes
	close(broadcast)

	// Verify that messages were broadcasted
	expectedMessages := []string{"Quote 1", "Quote 2", "Quote 3"}
	receivedMessages := make([]string, 0)

	for msg := range broadcast {
		var pl quotePayload
		if err := json.Unmarshal(msg, &pl); err != nil {
			t.Errorf("Failed to unmarshal message: %s", err)
		} else {
			receivedMessages = append(receivedMessages, pl.StockQuote)
		}
	}

	// Check that the received messages match the expected messages
	assert.ElementsMatch(t, expectedMessages, receivedMessages)
}

func mockLogger(t *testing.T) (*bufio.Scanner, *os.File, *os.File) {
	reader, writer, err := os.Pipe()
	if err != nil {
		assert.Fail(t, "couldn't get os Pipe: %v", err)
	}
	log.SetOutput(writer)

	return bufio.NewScanner(reader), reader, writer
}

func resetLogger(reader *os.File, writer *os.File) {
	err := reader.Close()
	if err != nil {
		fmt.Println("error closing reader was ", err)
	}
	if err = writer.Close(); err != nil {
		fmt.Println("error closing writer was ", err)
	}
	log.SetOutput(os.Stderr)
}
