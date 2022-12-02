package memory

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/cry0genic/go-stocks/finance"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	c := New()

	if len(finance.DefaultSymbols) != len(c.quotes) {
		t.Errorf("the quotes map length mismatches the default symbols slice length")
	}

	for _, symbol := range finance.DefaultSymbols {
		if _, ok := c.quotes[strings.ToLower(symbol)]; !ok {
			t.Errorf("%q not found in the quotes map", symbol)
		}
	}
}

func TestNewClientOptions(t *testing.T) {
	t.Parallel()

	symbols := []string{"foo", "bar"}
	c := New(Symbols(symbols))

	if len(symbols) != len(c.quotes) {
		t.Errorf("the quotes map length mismatches the optional symbols slice length")
	}

	for _, symbol := range symbols {
		if _, ok := c.quotes[strings.ToLower(symbol)]; !ok {
			t.Errorf("%q not found in the quotes map", symbol)
		}
	}
}

func TestGetQuotes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		quotes   []finance.Quote
		symbol   string
		last     int
		expected []finance.Quote
	}{
		{
			quotes: []finance.Quote{
				{Price: 123.45, Symbol: "fb"},
				{Price: 123.42, Symbol: "fb"},
			},
			symbol: "fb",
			last:   0,
			expected: []finance.Quote{
				{Price: 123.42, Symbol: "fb"},
			},
		},
	}

	c := New()

	for i, tc := range testCases {
		err := c.SetQuotes(context.Background(), tc.quotes)
		if err != nil {
			t.Errorf("%d: set actual: %v", i, err)
			continue
		}

		actual, err := c.GetQuotes(context.Background(), tc.symbol, tc.last)
		if err != nil {
			t.Errorf("%d: get quotes: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("%d: actual quotes not equal to expected", i)
			t.Logf("expected: %#v", tc.expected)
			t.Logf("actual:   %#v", actual)
		}
	}
}

func TestGetQuotesBatch(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		quotes   []finance.Quote
		symbols  []string
		last     int
		expected finance.QuoteBatch
	}{
		{
			quotes: []finance.Quote{
				{Price: 123.45, Symbol: "fb"},
				{Price: 123.42, Symbol: "fb"},
				{Price: 234.56, Symbol: "goog"},
			},
			symbols: []string{"fb", "goog"},
			last:    0,
			expected: finance.QuoteBatch{
				"fb": {
					{Price: 123.42, Symbol: "fb"},
				},
				"goog": {
					{Price: 234.56, Symbol: "goog"},
				},
			},
		},
	}

	c := New()

	for i, tc := range testCases {
		err := c.SetQuotes(context.Background(), tc.quotes)
		if err != nil {
			t.Errorf("%d: set actual: %v", i, err)
			continue
		}

		actual, err := c.GetQuotesBatch(context.Background(), tc.symbols, tc.last)
		if err != nil {
			t.Errorf("%d: get quotes batch: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("%d: actual quotes not equal to expected", i)
			t.Logf("expected: %#v", tc.expected)
			t.Logf("actual:   %#v", actual)
		}
	}
}
