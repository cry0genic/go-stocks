package history

import (
	"context"
	"fmt"

	"github.com/cry0genic/go-stocks/finance"
)

var ErrNotFound = fmt.Errorf("not found")

type Provider interface {
	GetQuotes(ctx context.Context, symbol string, last int) ([]finance.Quote, error)
	GetQuotesBatch(ctx context.Context, symbols []string, last int) (finance.QuoteBatch, error)
}
