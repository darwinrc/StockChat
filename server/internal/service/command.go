package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"server/internal/infra"
	"server/internal/model"
	"server/internal/repo"
	"strings"
	"time"
)

type CommmandService interface {
	ProcessCommand(command string)
	BroadcastCommand(broadcast chan []byte)
	ParseCommand(command string) (string, error)
}

type commandService struct {
	PostRepo   repo.PostRepo
	AMQPClient infra.AMQPClient
}

type stockPayload struct {
	StockCode string `json:"stockCode"`
}

type quotePayload struct {
	StockQuote string `json:"stockQuote"`
}

const (
	userID       = "48ccb5c1-9a19-42cd-bd41-3ac5c8af1108"
	username     = "StockBot"
	stockCommand = "/stock="
)

var commands []*model.Post

// NewCommandService builds a service and injects its dependencies
func NewCommandService(postRepo repo.PostRepo, amqpClient infra.AMQPClient) CommmandService {
	return &commandService{
		PostRepo:   postRepo,
		AMQPClient: amqpClient,
	}
}

// ParseCommand extracts the stock code from the command
// return "" if it is not a stock command, or it is malformed
func (s *commandService) ParseCommand(command string) (string, error) {
	if strings.Index(command, stockCommand) != 0 {
		return "", nil
	}

	tokens := strings.SplitAfter(command, "=")
	if len(tokens) != 2 {
		return "", errors.New(fmt.Sprintf("invalid command: %s. It should be something like /stock=aapl.us", command))
	}

	return tokens[1], nil
}

// ProcessCommand processes the command, publishing it to the rabbitmq exchange <stockchat>
func (s *commandService) ProcessCommand(stockCode string) {
	log.Println("Processing command for: ", stockCode)

	pl := stockPayload{
		StockCode: stockCode,
	}

	body, err := json.Marshal(pl)
	if err != nil {
		log.Printf("error marshaling payload: %s", err)
		return
	}

	err = s.AMQPClient.SetupAMQExchange()
	if err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
	}

	if err := s.AMQPClient.PublishAMQMessage(body); err != nil {
		log.Printf("error publishing to the exchange: %s", err)
	}

	log.Printf("Stock sent: %s\n", body)
}

// BroadcastCommand subscribes to the rabbitmq exchange <stockchat> and broadcasts the new quotes received
func (s *commandService) BroadcastCommand(broadcast chan []byte) {
	err := s.AMQPClient.SetupAMQExchange()
	defer s.AMQPClient.Close()

	if err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
	}

	messages, err := s.AMQPClient.ConsumeAMQMessages()
	if err != nil {
		log.Printf("error consuming messages: %s", err)
	}

	for message := range messages {
		var pl quotePayload
		if err := json.Unmarshal(message.Body, &pl); err != nil {
			log.Printf("error unmarshaling message: %s", err)
		}

		log.Printf("Quote received: %s\n", string(message.Body))

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

// addCommandToMemory adds a post to the commands in-memory list
func addCommandToMemory(post *model.Post) {
	commands = append([]*model.Post{post}, commands...)
}
