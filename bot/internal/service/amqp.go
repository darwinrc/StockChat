package service

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

const (
	exchangeName = "stockchat"
	queueName    = "stockchat-queue-stocks"
	stockKey     = "messages.stock"
	quoteKey     = "messages.quote"
)

// setupAMQExchange configures and returns a connection and exchange to rabbitmq
func setupAMQExchange() (*amqp.Connection, *amqp.Channel, error) {
	user, password, host := os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST")

	amqpUrl := fmt.Sprintf("amqp://%s:%s@%s/", user, password, host)

	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("error dialing amqp: %s", err))
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("error opening amqp channel: %s", err))
	}

	if err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("error declaring exchange: %s", err)
		return nil, nil, errors.New(fmt.Sprintf("error declaring amqp exchange: %s", err))
	}

	return conn, ch, nil
}

// publishAMQMessage publishes a message to the amq exchange
func publishAMQMessage(ch *amqp.Channel, message []byte) error {
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	}
	if err := ch.Publish(exchangeName, quoteKey, false, false, msg); err != nil {
		return err
	}

	return nil
}

// consumeAMQMessages returns the messages from the subscribed queue
func consumeAMQMessages(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error declaring queue: %s", err))
	}

	if err = ch.QueueBind(q.Name, stockKey, exchangeName, false, nil); err != nil {
		return nil, errors.New(fmt.Sprintf("error binding exchange to queue: %s", err))
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error consuming queued messages: %s", err))
	}

	return messages, nil
}
