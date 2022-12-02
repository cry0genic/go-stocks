package influxdb

import (
	"context"

	"github.com/cry0genic/go-stocks/finance"
	"github.com/cry0genic/go-stocks/history"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	_ history.Archiver = (*Client)(nil)
	_ history.Provider = (*Client)(nil)
)

type Client struct {
	idb influxdb2.Client
}

func (c Client) Close() error {
	panic("implement me")
}

func (c Client) GetQuotes(ctx context.Context, symbol string, last int) (
	[]finance.Quote, error) {
	panic("implement me")
}

func (c Client) GetQuotesBatch(ctx context.Context, symbols []string,
	last int) (finance.QuoteBatch, error) {
	panic("implement me")
}

func (c Client) SetQuotes(ctx context.Context, quotes []finance.Quote) error {
	panic("implement me")
}
