package gtmlp

import "time"

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

// Config holds scraping configuration
type Config struct {
	// XPath definitions
	Container string            // Repeating element selector
	Fields    map[string]string // Field name → XPath expression

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
