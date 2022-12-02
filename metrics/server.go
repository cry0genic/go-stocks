package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(options ...Option) *http.Server {
	s := &http.Server{
		Addr:              ":2112",
		Handler:           promhttp.Handler(),
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
	}

	for _, option := range options {
		option(s)
	}

	return s
}
