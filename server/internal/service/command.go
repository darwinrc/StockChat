package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"
	"server/internal/model"
	"strings"
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
	amqpUrl      = "amqp://guest:guest@localhost:5672/"
	exchangeName = "stockchat"
	queueRes     = "stockchat-queue-res"
)

var (
	quotes   = []string{}
	messages <-chan amqp.Delivery
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

	stockCode := strings.SplitAfter(command, "=")[1]
	pl := stockPayload{
		StockCode: stockCode,
	}

	body, err := json.Marshal(pl)
	if err != nil {
		log.Fatalf("error marshaling payload: %s", err)
		return
	}

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	}

	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		log.Fatalf("error dialing amqp: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error opening channel: %s", err)
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("error declaring exchange: %s", err)
	}

	if err = ch.Publish(exchangeName, "messages.stock", false, false, msg); err != nil {
		log.Fatalf("error publishing to the exchange: %s", err)
	}

	log.Printf("Stock sent: %s\n", body)

	//TODO.. Refactor

	qreq, err := ch.QueueDeclare(queueRes, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("error declaring queue: %s", err)
	}

	if err = ch.QueueBind(qreq.Name, "messages.quote", exchangeName, false, nil); err != nil {
		log.Fatalf("error binding exchange to queue: %s", err)
	}

	messages, err = ch.Consume(qreq.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("error consuming queued messages: %s", err)
	}

	for message := range messages {
		log.Printf("Quote received: %s\n", string(message.Body))

		var pl quotePayload

		if err := json.Unmarshal(message.Body, &pl); err != nil {
			log.Fatalf("error unmarshaling message: %s", err)
		}

		uID, _ := uuid.FromBytes([]byte("48ccb5c1-9a19-42cd-bd41-3ac5c8af1108"))

		post := &model.Post{
			UserID: "48ccb5c1-9a19-42cd-bd41-3ac5c8af1108",
			User: &model.User{
				ID:       uID,
				Username: "StockBot",
			},
			Message: pl.StockQuote,
		}

		posts, err := s.PostRepo.GetRecentPosts(context.Background(), postsLimit)
		if err != nil {
			log.Fatalf("error getting posts from database: %s", err)
		}

		posts = append([]*model.Post{post}, posts...)

		jsonPosts, err := json.Marshal(posts)
		if err != nil {
			log.Fatalf("error getting posts to json: %s", err)
			return
		}

		broadcast <- jsonPosts
	}

}
