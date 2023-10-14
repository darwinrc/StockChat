package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	stooqUrl = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
)

// getStockQuote fetches the stooq API and parses the returned CSV to extract the `Close` stock value
func getStockQuote(stockCode string) (float64, error) {
	url := fmt.Sprintf(stooqUrl, stockCode)

	resp, err := http.Get(url)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("error publishing to the exchange: %s", err))
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	// skip header row
	if _, err = reader.Read(); err != nil {
		return 0, errors.New(fmt.Sprintf("error reading header row: %s", err))
	}

	records, err := reader.Read()
	if err != nil {
		return 0, errors.New(fmt.Sprintf("error reading row: %s", err))
	}

	quote, err := strconv.ParseFloat(records[6], 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("error parsing quote value: %s", err))
	}

	return quote, nil
}
