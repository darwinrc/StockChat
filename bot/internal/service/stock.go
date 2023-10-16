package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type StockService struct{}

type stockPayload struct {
	StockCode string `json:"stockCode"`
}

type quotePayload struct {
	StockQuote string `json:"stockQuote"`
}

// NewStockService builds a service
func NewStockService() *StockService {
	return &StockService{}
}

// ProcessMessages subscribes to the rabbitmq exchange <stockchat> to get stock codes
// and publishes back the corresponding quotes fetched from the stooq api
func (s *StockService) ProcessMessages() {
	conn, ch, err := setupAMQExchange()
	defer conn.Close()
	defer ch.Close()

	if err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
	}

	messages, err := consumeAMQMessages(ch)
	if err != nil {
		log.Printf("error consuming messages: %s", err)
	}

	for message := range messages {
		log.Printf("Stock received: %s\n", string(message.Body))

		var spl stockPayload
		err := json.Unmarshal(message.Body, &spl)
		if err != nil {
			log.Printf("error unmarshaling payload: %s", err)
			return
		}

		var quote string

		stockQuote, err := getStockQuote(spl.StockCode)
		if err == nil {
			quote = fmt.Sprintf("%s quote is $%.2f per share", strings.ToUpper(spl.StockCode), stockQuote)
		} else {
			log.Printf("error getting stock quote from stooq: %s", err)

			if err.Error() == "stock code not found" {
				quote = fmt.Sprintf("%s is not a valid stock code. Please check stooq.com for the stock list", strings.ToUpper(spl.StockCode))
			}
		}

		qpl := quotePayload{
			StockQuote: quote,
		}

		body, err := json.Marshal(qpl)
		if err != nil {
			log.Printf("error unmarshaling payload: %s", err)
			return
		}

		if err := publishAMQMessage(ch, body); err != nil {
			log.Printf("error publishing to the exchange: %s", err)
		}

		log.Printf("Quote sent: %s\n", string(body))
	}
}
