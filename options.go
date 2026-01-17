package gtmlp

import (
	"time"

	"github.com/Hanivan/gtmlp/internal/httpclient"
)

const defaultUserAgent = "gtmlp/1.0 (HTML Parsing Library)"

// Option is a function that configures parsing behavior.
type Option func(*config)

// config holds configuration for parsing.
type config struct {
	timeout        time.Duration
	userAgent      string
	headers        map[string]string
	proxyURL       string
	maxRetries     int
	suppressErrors bool
}

// defaultConfig returns the default configuration.
func defaultConfig() *config {
	return &config{
		timeout:   30 * time.Second,
		userAgent: defaultUserAgent,
		headers:   make(map[string]string),
	}
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		c.timeout = d
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *config) {
		c.userAgent = ua
	}
}

// WithHeaders sets custom HTTP headers.
func WithHeaders(h map[string]string) Option {
	return func(c *config) {
		for k, v := range h {
			c.headers[k] = v
		}
	}
}

// WithProxy sets a proxy URL for HTTP requests.
func WithProxy(proxyURL string) Option {
	return func(c *config) {
		c.proxyURL = proxyURL
	}
}

// WithMaxRetries sets the maximum number of retries for failed HTTP requests.
func WithMaxRetries(maxRetries int) Option {
	return func(c *config) {
		c.maxRetries = maxRetries
	}
}

// WithRandomUserAgent explicitly enables random user agents (enabled by default).
func WithRandomUserAgent() Option {
	return func(c *config) {
		// Random UA is enabled by default, this is explicit control
		c.userAgent = "" // Will trigger random UA in client
	}
}

// WithDisableRandomUA disables random user agents and uses static user agent.
func WithDisableRandomUA() Option {
	return func(c *config) {
		c.userAgent = defaultUserAgent
	}
}

// WithSuppressErrors enables error suppression for XPath queries.
// When enabled, XPath errors return nil instead of error values.
func WithSuppressErrors() Option {
	return func(c *config) {
		c.suppressErrors = true
	}
}

// applyOptions applies options to config and returns an HTTP client.
func applyOptions(opts ...Option) (*config, *httpclient.Client) {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	clientOpts := []httpclient.ClientOption{
		httpclient.WithTimeout(cfg.timeout),
		httpclient.WithMaxRetries(cfg.maxRetries),
	}

	// Set user agent
	if cfg.userAgent == "" {
		// Empty string means use random UA (default behavior)
		clientOpts = append(clientOpts, httpclient.WithRandomUserAgent())
	} else {
		// Specific user agent provided
		clientOpts = append(clientOpts, httpclient.WithUserAgent(cfg.userAgent))
	}

	// Add headers
	if len(cfg.headers) > 0 {
		clientOpts = append(clientOpts, httpclient.WithHeaders(cfg.headers))
	}

	if cfg.proxyURL != "" {
		clientOpts = append(clientOpts, httpclient.WithProxy(cfg.proxyURL))
	}

	client := httpclient.NewClient(clientOpts...)

	return cfg, client
}
