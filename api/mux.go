package api

import (
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/cry0genic/go-stocks/history"
	"github.com/cry0genic/go-stocks/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func newMux(provider history.Provider, log *zap.SugaredLogger,
	instrument bool) *mux.Router {
	log = log.Named("mux")

	r := mux.NewRouter().StrictSlash(true)
	r.Use(gziphandler.GzipHandler, zapLoggerMiddleware(log))

	if instrument {
		r.Use(metricsMiddleware)
		log.Info("API instrumented")
	}

	s := r.Methods("GET").PathPrefix("/v1").Subrouter()
	s.HandleFunc("/stocks", stocks(provider, log))
	s.HandleFunc("/stock/{symbol:[a-zA-Z0-9]+}", stock(provider, log))

	return r
}

func metricsMiddleware(next http.Handler) http.Handler {
	return promhttp.InstrumentHandlerInFlight(metrics.ServerInFlightRequests,
		promhttp.InstrumentHandlerDuration(metrics.ServerRequestDuration,
			promhttp.InstrumentHandlerCounter(metrics.ServerAPIRequests,
				promhttp.InstrumentHandlerResponseSize(metrics.ServerResponseBytes,
					next,
				),
			),
		),
	)
}

func zapLoggerMiddleware(log *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				next.ServeHTTP(w, r)
				log.Debugf("%s - %s %s (%s)", r.RemoteAddr, r.Method,
					r.URL.EscapedPath(), time.Since(start))
			},
		)
	}
}
