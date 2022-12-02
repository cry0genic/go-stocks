package iexcloud

import (
	"net/http"
	"time"

	"github.com/cry0genic/go-stocks/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Option func(*Client)

func BatchEndpoint(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.batchEndpoint = url
		}
	}
}

func CallTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			c.timeout = timeout
		}
	}
}

func InstrumentHTTPClient() Option {
	return func(c *Client) {
		
		trace := &promhttp.InstrumentTrace{
			DNSStart: func(t float64) {
				metrics.ClientDNSDuration.WithLabelValues("dns_start").Observe(t)
			},
			DNSDone: func(t float64) {
				metrics.ClientDNSDuration.WithLabelValues("dns_done").Observe(t)
			},
			TLSHandshakeStart: func(t float64) {
				metrics.ClientTLSDuration.WithLabelValues("tls_handshake_start").Observe(t)
			},
			TLSHandshakeDone: func(t float64) {
				metrics.ClientTLSDuration.WithLabelValues("tls_handshake_done").Observe(t)
			},
		}

		roundTripper := promhttp.InstrumentRoundTripperInFlight(metrics.ClientInFlightRequests,
			promhttp.InstrumentRoundTripperCounter(metrics.ClientAPIRequests,
				promhttp.InstrumentRoundTripperTrace(trace,
					promhttp.InstrumentRoundTripperDuration(metrics.ClientRequestDuration,
						http.DefaultTransport,
					),
				),
			),
		)

		c.httpClient.Transport = roundTripper
	}
}
