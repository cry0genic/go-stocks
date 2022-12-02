package iexcloud

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cry0genic/go-stocks/finance"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestNewClientDefaults(t *testing.T) {
	t.Parallel()

	_, err := New("")
	if err != ErrInvalidToken {
		t.Errorf("expected: ErrInvalidToken; actual: %v", err)
	}

	c, err := New("stonks!")
	if err != nil {
		t.Error(err)
	}
	if c != nil {
		if c.batchEndpoint != DefaultBatchEndpoint {
			t.Errorf("expected endpoint: %q; actual endpoint: %q",
				DefaultBatchEndpoint, c.batchEndpoint)
		}
		if c.timeout != DefaultTimeout {
			t.Errorf("expected timeout: %v; actual timeout: %v",
				DefaultTimeout, c.timeout)
		}
		if c.httpClient == nil {
			t.Error("nil HTTP client")
		}
	}
}

func TestNewClientOptions(t *testing.T) {
	t.Parallel()

	endpoint := "https://nonexistent.domain"
	timeout := 42 * time.Second

	c, err := New(
		"to the moon!",
		BatchEndpoint(endpoint),
		CallTimeout(timeout),
		InstrumentHTTPClient(),
	)
	if err != nil {
		t.Error(err)
	}
	if c != nil {
		if c.batchEndpoint != endpoint {
			t.Errorf("expected endpoint: %q; actual endpoint: %q",
				endpoint, c.batchEndpoint)
		}
		if c.timeout != timeout {
			t.Errorf("expected timeout: %v; actual timeout: %v",
				timeout, c.timeout)
		}
		if c.httpClient == nil {
			t.Error("nil HTTP client")
		} else {
			if _, ok := c.httpClient.Transport.(promhttp.RoundTripperFunc); !ok {
				t.Error("underlying HTTP transport is not instrumented")
			}
		}
	}
}

func TestNewClientInvalidEndpoint(t *testing.T) {
	t.Parallel()

	_, err := New("this is fine", BatchEndpoint("blah\n"))
	if err == nil {
		t.Errorf("expected a batchQuotes endpoint error; actual: %q", err)
	}
}

func TestClientGetQuotes(t *testing.T) {
	t.Parallel()

	token := os.Getenv("STONKS_IEX_TOKEN")
	if token == "" {
		t.Logf("token not found in the STONKS_IEX_TOKEN environment variable")
		t.SkipNow()
	}

	c, err := New(token)
	if err != nil {
		t.Fatal(err)
	}

	quotes, err := c.GetQuotes(context.Background(), finance.DefaultSymbols...)
	if err != nil {
		t.Fatal(err)
	}

	for i, quote := range quotes {
		t.Logf("%d: %#v", i, quote)
		if quote.Price == 0 {
			t.Errorf("%d: price failed to decode", i)
		}
		if quote.Symbol == "" {
			t.Errorf("%d: symbol failed to decode", i)
		}
		if quote.Time.Equal(time.Time{}) {
			t.Errorf("%d: time failed to decode", i)
		}
	}
}
