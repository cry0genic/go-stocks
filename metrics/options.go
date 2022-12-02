package metrics

import (
	"net/http"
	"time"
)

type Option func(*http.Server)

func IdleTimeout(d time.Duration) Option {
	return func(s *http.Server) {
		s.IdleTimeout = d
	}
}

func ListenAddress(addr string) Option {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

func ReadHeaderTimeout(d time.Duration) Option {
	return func(s *http.Server) {
		s.ReadHeaderTimeout = d
	}
}
