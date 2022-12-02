package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(
		ClientAPIRequests,
		ClientDNSDuration,
		ClientInFlightRequests,
		ClientRequestDuration,
		ClientTLSDuration,
		ServerAPIRequests,
		ServerInFlightRequests,
		ServerRequestDuration,
		ServerResponseBytes,
	)
}

var ClientAPIRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "client_api_requests_total",
		Help: "A counter for requests from the wrapped client.",
	}, []string{},
)

var ClientDNSDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "client_dns_duration_seconds",
		Help:    "Trace DNS latency histogram.",
		Buckets: prometheus.DefBuckets,
	}, []string{},
)

var ClientInFlightRequests = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "client_in_flight_requests",
		Help: "A gauge of in-flight requests for the wrapped client.",
	},
)

var ClientRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "client_request_duration_seconds",
		Help:    "A histogram of request latencies.",
		Buckets: prometheus.DefBuckets,
	}, []string{},
)

var ClientTLSDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "client_tls_duration_seconds",
		Help:    "Trace TLS latency histogram.",
		Buckets: prometheus.DefBuckets,
	}, []string{},
)

var ServerAPIRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_requests_total",
		Help: "A counter for requests to the wrapped handler.",
	},
	[]string{},
)

var ServerInFlightRequests = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "server_in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	},
)

var ServerRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "A histogram of latencies for requests.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{},
)

var ServerResponseBytes = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "response_size_bytes",
		Help:    "A histogram of response sizes for requests.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{},
)
