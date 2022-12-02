package memory

import (
	"strings"

	"github.com/cry0genic/go-stocks/finance"
)

type Option func(*Client)

func Symbols(symbols []string) Option {
	s := make([]string, len(symbols))
	copy(s, symbols)

	return func(c *Client) {
		c.quotes = make(map[string][]finance.Quote)
		for _, symbol := range s {
			c.quotes[strings.ToLower(symbol)] = []finance.Quote{}
		}
	}
}
