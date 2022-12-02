package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/cry0genic/go-stocks/finance"
	"github.com/cry0genic/go-stocks/history"
	"go.uber.org/multierr"
)

var (
	_ history.Archiver = (*Client)(nil)
	_ history.Provider = (*Client)(nil)
)

type Client struct {
	mu     sync.RWMutex
	quotes map[string][]finance.Quote
}


func (c *Client) Close() error {
	return nil
}

func (c *Client) GetQuotes(_ context.Context, symbol string, last int) (
	[]finance.Quote, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	quotes, ok := c.quotes[strings.ToLower(symbol)]
	if !ok {
		return nil, history.ErrNotFound
	}

	if last < 1 {
		last = 1
	}
	if len(quotes) < last {
		last = len(quotes)
	}

	out := make([]finance.Quote, last)
	copy(out, quotes)

	return out, nil
}


func (c *Client) GetQuotesBatch(_ context.Context, symbols []string,
	last int) (finance.QuoteBatch, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(symbols) == 0 {
		symbols = finance.DefaultSymbols
	}

	batch := make(finance.QuoteBatch)
	for _, symbol := range symbols {
		quotes, ok := c.quotes[strings.ToLower(symbol)]
		if !ok {
			return nil, history.ErrNotFound
		}
		if len(quotes) == 0 {
			continue
		}

		if last < 1 {
			last = 1
		}
		if len(quotes) < last {
			last = len(quotes)
		}

		batch[symbol] = make([]finance.Quote, last)
		copy(batch[symbol], quotes)
	}

	return batch, nil
}

func (c *Client) SetQuotes(_ context.Context, quotes []finance.Quote) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error
	for _, quote := range quotes {
		symbol := strings.ToLower(quote.Symbol)
		quotes, ok := c.quotes[symbol]
		if !ok {
			multierr.AppendInto(&err, fmt.Errorf("symbol %q not found", quote.Symbol))
			continue
		}

		c.quotes[symbol] = append([]finance.Quote{quote}, quotes...)
	}

	return err
}

func New(options ...Option) *Client {
	c := &Client{quotes: make(map[string][]finance.Quote)}

	for _, symbol := range finance.DefaultSymbols {
		c.quotes[strings.ToLower(symbol)] = []finance.Quote{}
	}

	for _, option := range options {
		option(c)
	}

	return c
}
