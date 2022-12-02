package finance

import "context"

type Provider interface {
	GetQuotes(ctx context.Context, symbol ...string) ([]Quote, error)
}
