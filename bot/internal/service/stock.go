package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type StockService struct{}

type stockPayload struct {
	StockCode string `json:"stockCode"`
}

type quotePayload struct {
	StockQuote string `json:"stockQuote"`
}

const (
	amqpUrl      = "amqp://guest:guest@localhost:5672/"
	exchangeName = "stockchat"
	queueReq     = "stockchat-queue-req"

	stooqUrl = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
)

// NewStockService builds a service
func NewStockService() *StockService {
	return &StockService{}
}

// ProcessMessages ...
func (s *StockService) ProcessMessages() {
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

	qreq, err := ch.QueueDeclare(queueReq, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("error declaring queue: %s", err)
	}

	if err = ch.QueueBind(qreq.Name, "messages.stock", exchangeName, false, nil); err != nil {
		log.Fatalf("error binding exchange to queue: %s", err)
	}

	messages, err := ch.Consume(qreq.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("error consuming queued messages: %s", err)
	}

	for message := range messages {
		log.Printf("Stock received: %s\n", string(message.Body))

		var spl stockPayload
		err := json.Unmarshal(message.Body, &spl)
		if err != nil {
			log.Fatalf("error unmarshaling payload: %s", err)
			return
		}

		stockQuote := getStockQuote(spl.StockCode)

		qpl := quotePayload{
			StockQuote: fmt.Sprintf("%s quote is $%.2f per share", strings.ToUpper(spl.StockCode), stockQuote),
		}

		body, err := json.Marshal(qpl)
		if err != nil {
			log.Fatalf("error unmarshaling payload: %s", err)
			return
		}

		msg := amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		}

		if err := ch.Publish(exchangeName, "messages.quote", false, false, msg); err != nil {
			log.Fatalf("error publishing to the exchange: %s", err)
		}

		log.Printf("Quote sent: %s\n", string(msg.Body))
	}
}

// getStockQuote fetches the stooq API and parses the returned CSV to extract the `Close` stock value
func getStockQuote(stockCode string) float64 {
	url := fmt.Sprintf(stooqUrl, stockCode)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("error publishing to the exchange: %s", err)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	// skip header row
	if _, err = reader.Read(); err != nil {
		log.Fatalf("error reading header row: %s", err)
	}

	records, err := reader.Read()
	if err != nil {
		log.Fatalf("error reading row: %s", err)
	}

	quote, err := strconv.ParseFloat(records[6], 64)
	if err != nil {
		log.Fatalf("error parsing quote value: %s", err)
	}

	return quote
}
