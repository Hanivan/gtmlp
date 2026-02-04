package gtmlp

import (
	"errors"
	"fmt"
)

// ErrorType represents the category of error
type ErrorType string

const (
	ErrTypeNetwork    ErrorType = "network"
	ErrTypeParsing    ErrorType = "parsing"
	ErrTypeXPath      ErrorType = "xpath"
	ErrTypeConfig     ErrorType = "config"
	ErrTypeValidation ErrorType = "validation"
)

// ScrapeError is a typed error with context
type ScrapeError struct {
	Type    ErrorType
	Message string
	XPath   string
	URL     string
	Cause   error
}

func (e *ScrapeError) Error() string {
	if e.XPath != "" {
		return fmt.Sprintf("%s error: %s (xpath: %s)", e.Type, e.Message, e.XPath)
	}
	if e.URL != "" {
		return fmt.Sprintf("%s error: %s (url: %s)", e.Type, e.Message, e.URL)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

func (e *ScrapeError) Unwrap() error {
	return e.Cause
}

// Is checks if error is of specific type
func Is(err error, errorType ErrorType) bool {
	var scrapeErr *ScrapeError
	if errors.As(err, &scrapeErr) {
		return scrapeErr.Type == errorType
	}
	return false
}
