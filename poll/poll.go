package poll

import (
	"context"
	"fmt"
	"time"

	"github.com/cry0genic/go-stocks/finance"
	"github.com/cry0genic/go-stocks/history"
	"go.uber.org/zap"
)

const DefaultPollDuration = time.Minute

var (
	ErrNilArchiver = fmt.Errorf("archiver cannot be nil")
	ErrNilLogger   = fmt.Errorf("logger cannot be nil")
	ErrNilProvider = fmt.Errorf("finance provider cannot be nil")
)

type Poller struct {
	log      *zap.SugaredLogger
	archiver history.Archiver
	provider finance.Provider
}

func (p Poller) Poll(ctx context.Context, interval time.Duration,
	symbols ...string) {
	if len(symbols) == 0 {
		p.log.Warn("no symbols to poll")
		return
	}
	if interval <= 0 {
		p.log.Warn("invalid interval; using default 1 minute")
		interval = DefaultPollDuration
	}

	p.log.Infof("polling interval: %s", interval)
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		quotes, err := p.provider.GetQuotes(ctx, symbols...)
		if err != nil {
			p.log.Errorf("polling provider: %v", err)
		} else {
			p.log.Debugf("received: %#v", quotes)
			err = p.archiver.SetQuotes(ctx, quotes)
			if err != nil {
				p.log.Errorf("updating history: %v", err)
				continue
			}
			p.log.Debug("stored")
		}

		select {
		case <-ctx.Done():
			p.log.Debug("stopping poller")
			return
		case <-t.C:
		}
	}
}

func New(p finance.Provider, a history.Archiver, l *zap.SugaredLogger) (
	*Poller, error) {
	switch {
	case p == nil:
		return nil, ErrNilProvider
	case a == nil:
		return nil, ErrNilArchiver
	case l == nil:
		return nil, ErrNilLogger
	}

	return &Poller{
		log:      l.Named("poll"),
		archiver: a,
		provider: p,
	}, nil
}
