package infra

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type AMQPClient interface {
	SetupAMQExchange() error
	PublishAMQMessage(message []byte) error
	ConsumeAMQMessages() (<-chan amqp.Delivery, error)
	Close()
}

const (
	exchangeName = "stockchat"
	stockKey     = "messages.stock"
	quoteKey     = "messages.quote"
	queueName    = "stockchat-queue-quotes"
)

type amqpClient struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// NewAMQPClient builds an amqp client
func NewAMQPClient() AMQPClient {
	return &amqpClient{}
}

// SetupAMQExchange configures and returns a connection and exchange to rabbitmq
func (c *amqpClient) SetupAMQExchange() error {
	user, password, host := os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST")

	amqpUrl := fmt.Sprintf("amqp://%s:%s@%s/", user, password, host)

	var err error
	c.Conn, err = amqp.Dial(amqpUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("error dialing amqp: %s", err))
	}

	c.Ch, err = c.Conn.Channel()
	if err != nil {
		return errors.New(fmt.Sprintf("error opening amqp channel: %s", err))
	}

	if err = c.Ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Printf("error declaring exchange: %s", err)
		return errors.New(fmt.Sprintf("error declaring amqp exchange: %s", err))
	}

	return nil
}

// PublishAMQMessage publishes a message to the amq exchange
func (c *amqpClient) PublishAMQMessage(message []byte) error {
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	}
	if err := c.Ch.Publish(exchangeName, stockKey, false, false, msg); err != nil {
		return err
	}

	return nil
}

// ConsumeAMQMessages returns the messages from the subscribed queue
func (c *amqpClient) ConsumeAMQMessages() (<-chan amqp.Delivery, error) {
	q, err := c.Ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error declaring queue: %s", err))
	}

	if err = c.Ch.QueueBind(q.Name, quoteKey, exchangeName, false, nil); err != nil {
		return nil, errors.New(fmt.Sprintf("error binding exchange to queue: %s", err))
	}

	messages, err := c.Ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error consuming queued messages: %s", err))
	}

	return messages, nil
}

func (c *amqpClient) Close() {
	c.Conn.Close()
	c.Ch.Close()
}
