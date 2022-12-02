package sqlite

import "time"

type Option func(*Client)

func ConnMaxLifetime(d time.Duration) Option {
	return func(c *Client) {
		c.connsMaxLifetime = d
	}
}

func DatabaseFile(f string) Option {
	return func(c *Client) {
		if f != "" {
			c.file = f
		}
	}
}

func MaxIdleConnections(i int) Option {
	return func(c *Client) {
		c.maxIdleConns = i
	}
}

func Symbols(symbols []string) Option {
	return func(c *Client) {
		c.symbols = make(map[string]struct{})

		for _, symbol := range symbols {
			c.symbols[symbol] = struct{}{}
		}
	}
}
