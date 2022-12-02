package sqlite

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cry0genic/go-stocks/finance"
)

func TestGetQuotes(t *testing.T) {
	t.Parallel()

	now := time.Now()
	testCases := []struct {
		quotes   []finance.Quote
		symbol   string
		last     int
		expected []finance.Quote
	}{
		{ 
			quotes: []finance.Quote{
				{Price: 123.45, Symbol: "fb", Time: now},
				{Price: 123.42, Symbol: "fb", Time: now},
			},
			symbol: "fb",
			last:   0,
			expected: []finance.Quote{
				{Price: 123.42, Symbol: "fb", Time: now},
			},
		},
	}

	dir, err := ioutil.TempDir("", "stonks")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("removing temp dir: %v", err)
		}
	}()

	t.Logf("using temp directory %q", dir)

	c, err := New(DatabaseFile(filepath.Join(dir, DefaultDatabaseFile)))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	for i, tc := range testCases {
		err = c.SetQuotes(context.Background(), tc.quotes)
		if err != nil {
			t.Errorf("%d: set actual: %v", i, err)
			continue
		}

		actual, err := c.GetQuotes(context.Background(), tc.symbol, tc.last)
		if err != nil {
			t.Errorf("%d: get quotes: %v", i, err)
			continue
		}

		if len(actual) != len(tc.expected) {
			t.Errorf("%d: actual quote count not equal to expected count", i)
			continue
		}

		for j, q := range actual {
			if q.Price != tc.expected[i].Price {
				t.Errorf("%d.%d: actual price: %.2f; expected: %.2f", i, j,
					q.Price, tc.expected[i].Price)
			}
			if q.Symbol != tc.expected[i].Symbol {
				t.Errorf("%d.%d: actual symbol: %q; expected: %q", i, j,
					q.Symbol, tc.expected[i].Symbol)
			}
		}
	}
}

func TestGetQuotesBatch(t *testing.T) {
	t.Parallel()

	now := time.Now()
	testCases := []struct {
		quotes   []finance.Quote
		symbols  []string
		last     int
		expected finance.QuoteBatch
	}{
		{ 
			quotes: []finance.Quote{
				{Price: 123.45, Symbol: "fb", Time: now},
				{Price: 123.42, Symbol: "fb", Time: now},
				{Price: 123.40, Symbol: "fb", Time: now},
				{Price: 234.56, Symbol: "goog", Time: now},
				{Price: 234.51, Symbol: "goog", Time: now},
			},
			symbols: []string{"fb", "goog"},
			last:    2,
			expected: finance.QuoteBatch{
				"fb": {
					{Price: 123.40, Symbol: "fb"},
					{Price: 123.42, Symbol: "fb"},
				},
				"goog": {
					{Price: 234.51, Symbol: "goog"},
					{Price: 234.56, Symbol: "goog"},
				},
			},
		},
	}

	dir, err := ioutil.TempDir("", "stonks")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("removing temp dir: %v", err)
		}
	}()

	t.Logf("using temp directory %q", dir)

	c, err := New(DatabaseFile(filepath.Join(dir, DefaultDatabaseFile)))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	for i, tc := range testCases {
		err = c.SetQuotes(context.Background(), tc.quotes)
		if err != nil {
			t.Errorf("%d: set actual: %v", i, err)
			continue
		}

		actual, err := c.GetQuotesBatch(context.Background(), tc.symbols,
			tc.last)
		if err != nil {
			t.Errorf("%d: get quotes: %v", i, err)
			continue
		}

		if len(actual) != len(tc.expected) {
			t.Errorf("%d: actual quote count not equal to expected count", i)
			t.Logf("expected: %#v", tc.expected)
			t.Logf("actual:   %#v", actual)
			continue
		}

		for symbol := range actual {
			for j, q := range actual[symbol] {
				if q.Price != tc.expected[symbol][j].Price {
					t.Errorf("actual price: %.2f; expected: %.2f", q.Price,
						tc.expected[symbol][j].Price)
				}
				if q.Symbol != tc.expected[symbol][j].Symbol {
					t.Errorf("actual symbol: %q; expected: %q", q.Symbol,
						tc.expected[symbol][j].Symbol)
				}
			}
		}
	}
}
