package poll

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/cry0genic/go-stocks/finance"
	"github.com/cry0genic/go-stocks/history"
	"go.uber.org/zap/zaptest"
)

func TestNewPoller(t *testing.T) {
	t.Parallel()

	m := new(mockProviderArchiver)
	_, err := New(nil, nil, nil)
	if err != ErrNilProvider {
		t.Errorf("expected ErrNilProvider: %v", err)
	}

	_, err = New(m, nil, nil)
	if err != ErrNilArchiver {
		t.Errorf("expected ErrNilArchiver: %v", err)
	}

	_, err = New(m, m, nil)
	if err != ErrNilLogger {
		t.Errorf("expected ErrNilLogger: %v", err)
	}
}

func TestPoller(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now()
	expected := []finance.Quote{
		{Price: 123.42, Symbol: "fb", Time: now},
		{Price: 123.45, Symbol: "fb", Time: now},
	}
	m := &mockProviderArchiver{
		cancel: cancel,
		quotes: []finance.Quote{
			{Price: 123.45, Symbol: "fb", Time: now},
			{Price: 123.42, Symbol: "fb", Time: now},
		},
		storage: make([]finance.Quote, 0, 2),
	}

	p, err := New(m, m, zaptest.NewLogger(t).Sugar())
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		p.Poll(ctx, 100*time.Millisecond, "fb")
		close(done)
	}()
	<-done

	if !reflect.DeepEqual(m.storage, expected) {
		t.Error("storage does not equal expected")
		t.Logf("storage:  %#v", m.storage)
		t.Logf("expected: %#v", expected)
		t.Skip()
	}
}

var (
	_ finance.Provider = (*mockProviderArchiver)(nil)
	_ history.Archiver = (*mockProviderArchiver)(nil)
)

type mockProviderArchiver struct {
	cancel          context.CancelFunc
	quotes, storage []finance.Quote
}

func (m mockProviderArchiver) Close() error { return nil }

func (m *mockProviderArchiver) GetQuotes(_ context.Context, _ ...string) (
	[]finance.Quote, error) {
	var q finance.Quote
	if len(m.quotes) > 0 {
		q, m.quotes = m.quotes[0], m.quotes[1:]
	}
	if len(m.quotes) == 0 {
		m.cancel()
	}

	return []finance.Quote{q}, nil
}

func (m *mockProviderArchiver) SetQuotes(_ context.Context,
	quotes []finance.Quote) error {
	m.storage = append(quotes, m.storage...)

	return nil
}
