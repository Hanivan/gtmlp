package gtmlp

import (
	"errors"
	"fmt"
	"testing"
)

// TestScrapeErrorErrorWithXPath tests ScrapeError.Error() with XPath context
func TestScrapeErrorErrorWithXPath(t *testing.T) {
	err := &ScrapeError{
		Type:    ErrTypeXPath,
		Message: "invalid expression",
		XPath:   "//div[@class='test']",
	}

	expected := "xpath error: invalid expression (xpath: //div[@class='test'])"
	result := err.Error()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestScrapeErrorErrorWithURL tests ScrapeError.Error() with URL context
func TestScrapeErrorErrorWithURL(t *testing.T) {
	err := &ScrapeError{
		Type:    ErrTypeNetwork,
		Message: "connection refused",
		URL:     "https://example.com",
	}

	expected := "network error: connection refused (url: https://example.com)"
	result := err.Error()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestScrapeErrorErrorWithoutContext tests ScrapeError.Error() without XPath or URL
func TestScrapeErrorErrorWithoutContext(t *testing.T) {
	err := &ScrapeError{
		Type:    ErrTypeConfig,
		Message: "missing required field",
	}

	expected := "config error: missing required field"
	result := err.Error()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestScrapeErrorErrorWithBothXPathAndURL tests priority when both XPath and URL are present
func TestScrapeErrorErrorWithBothXPathAndURL(t *testing.T) {
	err := &ScrapeError{
		Type:    ErrTypeXPath,
		Message: "test error",
		XPath:   "//div/text()",
		URL:     "https://example.com",
	}

	// XPath should take priority over URL
	expected := "xpath error: test error (xpath: //div/text())"
	result := err.Error()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestScrapeErrorUnwrap tests ScrapeError.Unwrap() method
func TestScrapeErrorUnwrap(t *testing.T) {
	t.Run("with cause", func(t *testing.T) {
		originalErr := errors.New("original error")
		err := &ScrapeError{
			Type:    ErrTypeParsing,
			Message: "failed to parse",
			Cause:   originalErr,
		}

		unwrapped := err.Unwrap()

		if unwrapped != originalErr {
			t.Errorf("Expected unwrapped error to be original error, got %v", unwrapped)
		}
	})

	t.Run("without cause", func(t *testing.T) {
		err := &ScrapeError{
			Type:    ErrTypeParsing,
			Message: "failed to parse",
			Cause:   nil,
		}

		unwrapped := err.Unwrap()

		if unwrapped != nil {
			t.Errorf("Expected unwrapped error to be nil, got %v", unwrapped)
		}
	})

	t.Run("with chain of errors", func(t *testing.T) {
		// Test that errors.Is and errors.As work correctly
		baseErr := errors.New("base error")
		scrapeErr := &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "network failure",
			Cause:   fmt.Errorf("wrapped: %w", baseErr),
		}

		// Should be able to unwrap to the base error
		if !errors.Is(scrapeErr, baseErr) {
			t.Error("Expected errors.Is to find base error in chain")
		}
	})
}

// TestIsFunction tests the Is() function with correct type
func TestIsFunctionWithCorrectType(t *testing.T) {
	err := &ScrapeError{
		Type:    ErrTypeXPath,
		Message: "invalid xpath",
		XPath:   "//test",
	}

	if !Is(err, ErrTypeXPath) {
		t.Error("Expected Is to return true for matching error type")
	}
}

// TestIsFunctionWithWrongType tests the Is() function with wrong type
func TestIsFunctionWithWrongType(t *testing.T) {
	err := &ScrapeError{
		Type:    ErrTypeParsing,
		Message: "parse error",
	}

	if Is(err, ErrTypeXPath) {
		t.Error("Expected Is to return false for non-matching error type")
	}

	if Is(err, ErrTypeNetwork) {
		t.Error("Expected Is to return false for different error type")
	}
}

// TestIsWithNilError tests the Is() function with nil error
func TestIsWithNilError(t *testing.T) {
	if Is(nil, ErrTypeXPath) {
		t.Error("Expected Is to return false for nil error")
	}

	if Is(nil, ErrTypeNetwork) {
		t.Error("Expected Is to return false for nil error")
	}
}

// TestIsWithNonScrapeError tests the Is() function with non-ScrapeError types
func TestIsWithNonScrapeError(t *testing.T) {
	standardErr := errors.New("standard error")

	if Is(standardErr, ErrTypeXPath) {
		t.Error("Expected Is to return false for non-ScrapeError type")
	}

	if Is(standardErr, ErrTypeNetwork) {
		t.Error("Expected Is to return false for non-ScrapeError type")
	}
}

// TestIsWithWrappedScrapeError tests Is() with wrapped ScrapeError
func TestIsWithWrappedScrapeError(t *testing.T) {
	scrapeErr := &ScrapeError{
		Type:    ErrTypeParsing,
		Message: "parsing failed",
	}
	wrappedErr := fmt.Errorf("wrapped: %w", scrapeErr)

	if !Is(wrappedErr, ErrTypeParsing) {
		t.Error("Expected Is to return true for wrapped ScrapeError")
	}

	if Is(wrappedErr, ErrTypeXPath) {
		t.Error("Expected Is to return false for different error type")
	}
}

// TestAllErrorTypeConstants tests all ErrorType constants
func TestAllErrorTypeConstants(t *testing.T) {
	tests := []struct {
		name      string
		errorType ErrorType
		expected  string
	}{
		{"ErrTypeNetwork", ErrTypeNetwork, "network"},
		{"ErrTypeParsing", ErrTypeParsing, "parsing"},
		{"ErrTypeXPath", ErrTypeXPath, "xpath"},
		{"ErrTypeConfig", ErrTypeConfig, "config"},
		{"ErrTypeValidation", ErrTypeValidation, "validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.errorType) != tt.expected {
				t.Errorf("Expected ErrorType to be %q, got %q", tt.expected, tt.errorType)
			}

			// Test that we can create a ScrapeError with each type
			err := &ScrapeError{
				Type:    tt.errorType,
				Message: "test message",
			}

			errMsg := err.Error()
			expectedPrefix := string(tt.errorType) + " error: test message"
			if errMsg != expectedPrefix {
				t.Errorf("Expected error message %q, got %q", expectedPrefix, errMsg)
			}
		})
	}
}

// TestErrorTypeUniqueness verifies all ErrorType constants are unique
func TestErrorTypeUniqueness(t *testing.T) {
	types := []ErrorType{
		ErrTypeNetwork,
		ErrTypeParsing,
		ErrTypeXPath,
		ErrTypeConfig,
		ErrTypeValidation,
	}

	seen := make(map[ErrorType]bool)
	for _, et := range types {
		if seen[et] {
			t.Errorf("Duplicate ErrorType found: %s", et)
		}
		seen[et] = true
	}

	if len(seen) != len(types) {
		t.Error("Not all ErrorType constants are unique")
	}
}

// TestScrapeErrorCreation tests creating ScrapeError with all fields
func TestScrapeErrorCreation(t *testing.T) {
	cause := errors.New("underlying error")
	err := &ScrapeError{
		Type:    ErrTypeValidation,
		Message: "validation failed",
		XPath:   "//input[@id='email']",
		URL:     "https://example.com/form",
		Cause:   cause,
	}

	if err.Type != ErrTypeValidation {
		t.Error("Type not set correctly")
	}

	if err.Message != "validation failed" {
		t.Error("Message not set correctly")
	}

	if err.XPath != "//input[@id='email']" {
		t.Error("XPath not set correctly")
	}

	if err.URL != "https://example.com/form" {
		t.Error("URL not set correctly")
	}

	if err.Cause != cause {
		t.Error("Cause not set correctly")
	}
}

// TestErrorsAsCompatibility tests compatibility with errors.As
func TestErrorsAsCompatibility(t *testing.T) {
	scrapeErr := &ScrapeError{
		Type:    ErrTypeNetwork,
		Message: "connection timeout",
	}

	// Test errors.As can extract ScrapeError
	var extracted *ScrapeError
	if !errors.As(scrapeErr, &extracted) {
		t.Error("errors.As failed to extract ScrapeError")
	}

	if extracted != scrapeErr {
		t.Error("Extracted error is not the same as original")
	}

	// Test with wrapped error
	wrapped := fmt.Errorf("wrapped: %w", scrapeErr)
	var extracted2 *ScrapeError
	if !errors.As(wrapped, &extracted2) {
		t.Error("errors.As failed to extract wrapped ScrapeError")
	}
}
