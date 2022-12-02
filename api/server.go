package api

import (
	"context"
	"net/http"
	"time"

	"github.com/cry0genic/go-stocks/history"
	"go.uber.org/zap"
)

const (
	DefaultIdleTimeout = time.Minute

	DefaultListenAddress = ":18081"

	DefaultReadHeaderTimeout = 30 * time.Second
)

type Server struct {
	ctx               context.Context
	srv               *http.Server
	log               *zap.SugaredLogger
	listenAddr        string
	idleTimeout       time.Duration
	readHeaderTimeout time.Duration
	instrumentation   bool
}

func (s *Server) ListenAndServe() error {
	go func() {
		<-s.ctx.Done()
		s.log.Info("shutting down ...")
		sCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.srv.Shutdown(sCtx)
	}()

	s.log.Infof("Listening on %q", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) ListenAndServeTLS(cert, pkey string) error {
	go func() {
		<-s.ctx.Done()
		s.log.Info("shutting down ...")
		sCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.srv.Shutdown(sCtx)
	}()

	s.log.Infof("Listening on %q (TLS)", s.srv.Addr)
	return s.srv.ListenAndServeTLS(cert, pkey)
}

func New(ctx context.Context, p history.Provider, log *zap.SugaredLogger,
	options ...Option) (
	*Server, error) {
	s := &Server{
		ctx:               ctx,
		log:               log.Named("api"),
		listenAddr:        DefaultListenAddress,
		idleTimeout:       DefaultIdleTimeout,
		readHeaderTimeout: DefaultReadHeaderTimeout,
		instrumentation:   true,
	}

	for _, option := range options {
		if option != nil {
			option(s)
		}
	}

	s.srv = &http.Server{
		Addr:              s.listenAddr,
		IdleTimeout:       s.idleTimeout,
		ReadHeaderTimeout: s.readHeaderTimeout,
		Handler:           newMux(p, log, s.instrumentation),
	}

	return s, nil
}
