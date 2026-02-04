package gtmlp

import (
	"context"
	"time"
)

// EnvMapping defines configurable environment variable names
type EnvMapping struct {
	Timeout    string
	UserAgent  string
	RandomUA   string
	MaxRetries string
	Proxy      string
}

// DefaultEnvMapping provides default env var names
var DefaultEnvMapping = &EnvMapping{
	Timeout:    "GTMLP_TIMEOUT",
	UserAgent:  "GTMLP_USER_AGENT",
	RandomUA:   "GTMLP_RANDOM_UA",
	MaxRetries: "GTMLP_MAX_RETRIES",
	Proxy:      "GTMLP_PROXY",
}

// FieldConfig defines a single field's XPath and optional pipes
type FieldConfig struct {
	XPath string
	Pipes []string
}

// Config holds scraping configuration
type Config struct {
	// XPath definitions
	Container string                 // Repeating element selector
	Fields    map[string]FieldConfig // Field name → FieldConfig

	// HTTP options
	Timeout    time.Duration
	UserAgent  string
	RandomUA   bool
	MaxRetries int
	Proxy      string
	Headers    map[string]string
}

// PartialResult contains data and field-level errors
type PartialResult[T any] struct {
	Data   []T
	Errors map[string]error
}

// PipeFunc defines a pipe transformation function
type PipeFunc func(ctx context.Context, input string, params []string) (any, error)

// WithURL adds the base URL to context for parseUrl pipe
func WithURL(ctx context.Context, url string) context.Context {
	return context.WithValue(ctx, contextKey("baseURL"), url)
}

type contextKey string
