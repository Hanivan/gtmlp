package gtmlp

import (
	"context"
	"fmt"
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
	XPath    string
	AltXPath []string
	Pipes    []string
}

// Config holds scraping configuration
type Config struct {
	// XPath definitions
	Container    string                 // Repeating element selector
	AltContainer []string               // Alternative container selectors
	Fields       map[string]FieldConfig // Field name → FieldConfig

	// Pagination
	Pagination *PaginationConfig // Optional pagination configuration

	// Security options
	URLValidator    func(string) error // Optional custom URL validation function
	AllowPrivateIPs bool               // Allow scraping private/internal IPs (default: false)

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

// PaginationConfig defines pagination behavior
type PaginationConfig struct {
	Type         string        // "next-link" or "numbered"
	NextSelector string        // XPath for next link (next-link type)
	AltSelectors []string      // Fallback selectors for next link
	PageSelector string        // XPath for all page links (numbered type)
	Pipes        []string      // URL transformation pipes
	MaxPages     int           // Maximum pages to scrape (default: 100)
	Timeout      time.Duration // Total pagination timeout (default: 10m)
}

// PaginatedResults contains page-separated scraping results
type PaginatedResults[T any] struct {
	Pages      []PageResult[T]
	TotalPages int
	TotalItems int
}

// PageResult contains results from a single page
type PageResult[T any] struct {
	URL       string
	PageNum   int
	Items     []T
	ScrapedAt time.Time
}

// PaginationInfo contains extracted pagination URLs
type PaginationInfo struct {
	URLs    []string // All discovered page URLs
	Type    string   // "next-link" or "numbered"
	BaseURL string   // Original base URL
}

// PaginationError represents an error during pagination
type PaginationError struct {
	PageURL      string // URL that failed
	PageNumber   int    // Page number (1-indexed)
	PartialData  any    // Items scraped before failure
	TotalScraped int    // Total items before failure
	Cause        error  // Underlying error
}

func (e *PaginationError) Error() string {
	return fmt.Sprintf("pagination failed at page %d (%s): %v",
		e.PageNumber, e.PageURL, e.Cause)
}

// WithURL adds the base URL to context for parseUrl pipe
func WithURL(ctx context.Context, url string) context.Context {
	return context.WithValue(ctx, contextKey("baseURL"), url)
}

type contextKey string
