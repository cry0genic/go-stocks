package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cry0genic/go-stocks/finance"
	"github.com/cry0genic/go-stocks/history/memory"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	provider = memory.New()
	log      *zap.SugaredLogger
	router   *mux.Router
)

func TestStockHandler(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/stock/blah?last=2", nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("nonexistent symbol results in code: %q", http.StatusText(w.Code))
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/stock/fb?last=blah", nil))
	if w.Code != http.StatusBadRequest {
		t.Errorf("bad 'last' parameter results in code: %q", http.StatusText(w.Code))
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/stock/fb?last=2", nil))
	t.Log(w.Body)

	var actual []finance.Quote
	err := json.NewDecoder(w.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := []finance.Quote{
		{Price: 123.40, Symbol: "fb", Time: time.Now().Add(time.Hour)},
		{Price: 123.42, Symbol: "fb", Time: time.Now().Add(time.Minute)},
	}

	if len(actual) != len(expected) {
		t.Error("actual quote count not equal to expected count")
		t.Logf("expected: %#v", expected)
		t.Logf("actual:   %#v", actual)
		t.Skip()
	}

	for i, q := range actual {
		if q.Price != expected[i].Price {
			t.Errorf("actual price: %.2f; expected: %.2f", q.Price,
				expected[i].Price)
		}
		if q.Symbol != expected[i].Symbol {
			t.Errorf("actual symbol: %q; expected: %q", q.Symbol,
				expected[i].Symbol)
		}
	}
}

func TestStocksHandler(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/stocks?last=blah", nil))
	if w.Code != http.StatusBadRequest {
		t.Errorf("bad 'last' parameter results in code: %q", http.StatusText(w.Code))
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/stocks?last=2", nil))

	var actual finance.QuoteBatch
	err := json.NewDecoder(w.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := finance.QuoteBatch{
		"fb": {
			{Price: 123.40, Symbol: "fb"},
			{Price: 123.42, Symbol: "fb"},
		},
		"goog": {
			{Price: 234.51, Symbol: "goog"},
			{Price: 234.56, Symbol: "goog"},
		},
	}

	if len(actual) != len(expected) {
		t.Error("actual quote count not equal to expected count")
		t.Logf("expected: %#v", expected)
		t.Logf("actual:   %#v", actual)
		t.Skip()
	}

	for symbol := range actual {
		for i, q := range actual[symbol] {
			if q.Price != expected[symbol][i].Price {
				t.Errorf("actual price: %.2f; expected: %.2f", q.Price,
					expected[symbol][i].Price)
			}
			if q.Symbol != expected[symbol][i].Symbol {
				t.Errorf("actual symbol: %q; expected: %q", q.Symbol,
					expected[symbol][i].Symbol)
			}
		}
	}
}

func init() {
	log = zap.NewExample().Sugar()
	quotes := []finance.Quote{
		{Price: 123.45, Symbol: "fb", Time: time.Now()},
		{Price: 123.42, Symbol: "fb", Time: time.Now().Add(time.Minute)},
		{Price: 123.40, Symbol: "fb", Time: time.Now().Add(time.Hour)},
		{Price: 234.56, Symbol: "goog", Time: time.Now()},
		{Price: 234.51, Symbol: "goog", Time: time.Now().Add(time.Minute)},
	}

	if err := provider.SetQuotes(context.Background(), quotes); err != nil {
		log.Fatal(err)
	}
	router = newMux(provider, log, false)
}
