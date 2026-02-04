package gtmlp

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
)

// ValidationResult represents XPath validation result for config-based validation
type ValidationResult struct {
	XPath      string
	Valid      bool
	MatchCount int
	Error      error
}

// ValidateXPath validates XPath expressions against HTML
func ValidateXPath(html string, xpaths map[string]string) map[string]ValidationResult {
	results := make(map[string]ValidationResult)

	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		// All fail if parsing fails
		for field, xpath := range xpaths {
			results[field] = ValidationResult{
				XPath: xpath,
				Valid: false,
				Error: err,
			}
		}
		return results
	}

	for field, xpathExpr := range xpaths {
		result := ValidationResult{
			XPath: xpathExpr,
		}

		// Try to compile the XPath to check syntax
		_, err := xpath.Compile(xpathExpr)
		if err != nil {
			result.Valid = false
			result.Error = err
			results[field] = result
			continue
		}

		// Test the expression by finding nodes
		nodes := htmlquery.Find(doc, xpathExpr)
		result.Valid = true
		result.MatchCount = len(nodes)
		results[field] = result
	}

	return results
}

// ValidateXPathURL validates XPath expressions from a URL
func ValidateXPathURL(url string, config *Config) (map[string]ValidationResult, error) {
	// Fetch HTML content from URL
	html, err := fetchHTML(url, config)
	if err != nil {
		return nil, err
	}

	// Validate XPath expressions against fetched HTML
	// Use container and fields from config
	xpaths := make(map[string]string)
	xpaths["container"] = config.Container
	for field, xpath := range config.Fields {
		xpaths[field] = xpath
	}

	results := ValidateXPath(html, xpaths)
	return results, nil
}
