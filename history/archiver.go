package history

import (
	"context"
	"io"

	"github.com/cry0genic/go-stocks/finance"
)

type Archiver interface {
	SetQuotes(ctx context.Context, quotes []finance.Quote) error
	io.Closer
}
