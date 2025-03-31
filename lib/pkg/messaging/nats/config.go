package nats

import (
	"time"

	"github.com/nats-io/nats.go"
)

type Config struct {
	URL            string
	MaxReconnects  int
	ReconnectWait  time.Duration
	ConnectTimeout time.Duration
	Options        []nats.Option
}

type ConfigOption func(*Config)

func WithURL(url string) ConfigOption {
	return func(c *Config) {
		c.URL = url
	}
}

func WithMaxReconnects(max int) ConfigOption {
	return func(c *Config) {
		c.MaxReconnects = max
	}
}

func WithReconnectWait(wait time.Duration) ConfigOption {
	return func(c *Config) {
		c.ReconnectWait = wait
	}
}

func WithConnectTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.ConnectTimeout = timeout
	}
}

func WithNATSOptions(opts ...nats.Option) ConfigOption {
	return func(c *Config) {
		c.Options = append(c.Options, opts...)
	}
}

func WithNoReconnect() ConfigOption {
	return func(c *Config) {
		c.Options = append(c.Options, nats.NoReconnect())
	}
}

func WithPingInterval(interval time.Duration) ConfigOption {
	return func(c *Config) {
		c.Options = append(c.Options, nats.PingInterval(interval))
	}
}

func WithFlushTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.Options = append(c.Options, nats.FlusherTimeout(timeout))
	}
}

func NewConfig(opts ...ConfigOption) *Config {
	cfg := &Config{
		URL:            nats.DefaultURL,
		MaxReconnects:  5,
		ReconnectWait:  time.Second * 2,
		ConnectTimeout: time.Second * 5,
		Options:        make([]nats.Option, 0),
	}

	cfg.Options = append(cfg.Options,
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.Timeout(cfg.ConnectTimeout),
	)

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
