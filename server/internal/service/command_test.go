package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	mock_infra "server/internal/infra/mocks"
	"server/internal/model"
	mock_repo "server/internal/repo/mocks"
	"testing"
	"time"
)

func TestProcessCommandSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := []byte("{\"stockCode\":\"aapl.us\"}")

	mockAMQP := mock_infra.NewMockAMQPClient(ctrl)
	mockAMQP.EXPECT().SetupAMQExchange().Return(nil)
	mockAMQP.EXPECT().PublishAMQMessage(payload).Return(nil)

	mockPostRepo := &mock_repo.MockPostRepo{}

	scanner, reader, writer := mockLogger(t)
	defer resetLogger(reader, writer)

	service := NewCommandService(mockPostRepo, mockAMQP)
	go service.ProcessCommand("/stock=aapl.us")

	scanner.Scan() // first log: Processing command ...
	scanner.Scan() // last log: Stock sent ...
	got := scanner.Text()
	msg := "Stock sent: {\"stockCode\":\"aapl.us\"}"
	assert.Contains(t, got, msg, "Expected to have a valid command such as  \"Stock sent: {\\\"stockCode\\\":\\\"aapl.us\\\"}\"")
}

func TestProcessCommandFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := &mock_repo.MockPostRepo{}
	mockAMQP := mock_infra.NewMockAMQPClient(ctrl)

	scanner, reader, writer := mockLogger(t)
	defer resetLogger(reader, writer)

	service := NewCommandService(mockPostRepo, mockAMQP)
	go service.ProcessCommand("invalid-command")

	scanner.Scan() // first log: Processing command ...
	scanner.Scan() // last log: invalid command ...
	got := scanner.Text()
	msg := "invalid command: invalid-command. It should be something like /stock=aapl.us"
	assert.Contains(t, got, msg, "Expected to have an invalid command")
}

func TestBroadcastCommandSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msg := "AAPL.US quote is $178.85 per share"

	messages := make(chan amqp.Delivery)
	go func() {
		messages <- amqp.Delivery{Body: []byte(fmt.Sprintf("{\"stockQuote\":\"%s\"}", msg))}
	}()

	mockAMQP := mock_infra.NewMockAMQPClient(ctrl)
	mockAMQP.EXPECT().SetupAMQExchange().Return(nil)
	mockAMQP.EXPECT().ConsumeAMQMessages().Return(messages, nil)

	uID, _ := uuid.FromBytes([]byte(userID))
	ts1 := time.Now().Add(time.Hour * -1)
	ts2 := time.Now().Add(time.Hour)

	posts := []*model.Post{
		{
			UserID: userID,
			User: &model.User{
				ID:       uID,
				Username: username,
			},
			Message:   "Dummy Quote 1",
			Timestamp: &ts1,
		},
		{
			UserID: userID,
			User: &model.User{
				ID:       uID,
				Username: username,
			},
			Message:   "Dummy Quote 2",
			Timestamp: &ts2,
		},
	}

	mockPostRepo := mock_repo.NewMockPostRepo(ctrl)
	mockPostRepo.EXPECT().GetRecentPosts(context.Background(), postsLimit).Return(posts, nil)

	broadcast := make(chan []byte)
	defer close(broadcast)
	scanner, reader, writer := mockLogger(t)
	defer resetLogger(reader, writer)

	service := NewCommandService(mockPostRepo, mockAMQP)
	go service.BroadcastCommand(broadcast)

	scanner.Scan()
	got := scanner.Text()
	txt := fmt.Sprintf("Quote received: {\"stockQuote\":\"%s\"}", msg)

	assert.Contains(t, got, txt)

	// assert addCommandToMemory(post)
	assert.True(t, containsMessage(commands, msg), "Expected to find the message in the memory db")

	// assert broadcastPosts(s.PostRepo, broadcast)
	receivedMessages := make([]*model.Post, 0)

	for m := range broadcast {
		var pl []*model.Post
		if err := json.Unmarshal(m, &pl); err != nil {
			t.Errorf("Failed to unmarshal message: %s", err)
		} else {
			receivedMessages = append(receivedMessages, pl...)
			break
		}
	}

	fmt.Println("Received messages: ", receivedMessages)
	assert.True(t, containsMessage(receivedMessages, msg), "Expected to find the message in the broadcast channel")
}

// Util functions
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

func containsMessage(posts []*model.Post, message string) bool {
	for _, c := range posts {
		if c.Message == message {
			return true
		}
	}

	return false
}
