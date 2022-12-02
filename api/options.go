package api

import "time"

type Option func(*Server)

func DisableInstrumentation() Option {
	return func(s *Server) {
		s.instrumentation = false
	}
}

func IdleTimeout(d time.Duration) Option {
	return func(s *Server) {
		if d > 0 {
			s.idleTimeout = d
		}
	}
}

func ListenAddress(addr string) Option {
	return func(s *Server) {
		if addr != "" {
			s.listenAddr = addr
		}
	}
}

func ReadHeaderTimeout(d time.Duration) Option {
	return func(s *Server) {
		if d > 0 {
			s.readHeaderTimeout = d
		}
	}
}
