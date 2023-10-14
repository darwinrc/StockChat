package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"server/internal/model"
	"strings"
	"time"
)

type CommandService struct {
	PostRepo model.PostRepo
}

type stockPayload struct {
	StockCode string `json:"stockCode"`
}

type quotePayload struct {
	StockQuote string `json:"stockQuote"`
}

const (
	userID   = "48ccb5c1-9a19-42cd-bd41-3ac5c8af1108"
	username = "StockBot"
)

// NewCommandService builds a service and injects its dependencies
func NewCommandService(postRepo model.PostRepo) *CommandService {
	return &CommandService{
		PostRepo: postRepo,
	}
}

// ProcessCommand processes the command, publishing it to the rabbitmq exchange <stockchat>
func (s *CommandService) ProcessCommand(command string, broadcast chan []byte) {
	log.Println("Processing command: ", command)

	pl := stockPayload{
		StockCode: strings.SplitAfter(command, "=")[1],
	}

	body, err := json.Marshal(pl)
	if err != nil {
		log.Fatalf("error marshaling payload: %s", err)
		return
	}

	conn, ch, err := setupAMQExchange()
	defer conn.Close()
	defer ch.Close()

	if err != nil {
		log.Fatalf("error setting up the amq connection and exchange: %s", err)
	}

	if err := publishAMQMessage(ch, body); err != nil {
		log.Fatalf("error publishing to the exchange: %s", err)
	}

	log.Printf("Stock sent: %s\n", body)
}

// BroadcastCommand subscribes to the rabbitmq exchange <stockchat> and broadcasts the new quotes received
func (s *CommandService) BroadcastCommand(broadcast chan []byte) {
	conn, ch, err := setupAMQExchange()
	defer conn.Close()
	defer ch.Close()

	if err != nil {
		log.Fatalf("error setting up the amq connection and exchange: %s", err)
	}

	messages, err := consumeAMQMessages(ch)
	if err != nil {
		log.Fatalf("error consuming messages: %s", err)
	}

	for message := range messages {
		log.Printf("Quote received: %s\n", string(message.Body))

		var pl quotePayload
		if err := json.Unmarshal(message.Body, &pl); err != nil {
			log.Fatalf("error unmarshaling message: %s", err)
		}

		uID, _ := uuid.FromBytes([]byte(userID))
		ts := time.Now().UTC()

		post := &model.Post{
			UserID: userID,
			User: &model.User{
				ID:       uID,
				Username: username,
			},
			Message:   pl.StockQuote,
			Timestamp: &ts,
		}

		addCommandToMemory(post)
		broadcastPosts(s.PostRepo, broadcast)
	}
}
