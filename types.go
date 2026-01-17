package gtmlp

import (
	"fmt"
)

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message string
	Err     error
}

// Error returns the error message.
func (e *ParseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error.
func (e *ParseError) Unwrap() error {
	return e.Err
}

// NewParseError creates a new ParseError.
func NewParseError(message string, err error) *ParseError {
	return &ParseError{
		Message: message,
		Err:     err,
	}
}
