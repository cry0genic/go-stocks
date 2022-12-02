package finance

import "time"

var DefaultSymbols = []string{"fb", "amzn", "aapl", "nflx", "goog"}

type Quote struct {
	Price  float64   `json:"price"`
	Symbol string    `json:"symbol"`
	Time   time.Time `json:"time"`
}

type QuoteBatch map[string][]Quote
