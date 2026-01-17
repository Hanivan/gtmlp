package validator

import (
	"github.com/Hanivan/gtmlp/internal/parser"
)

// ValidationResult represents the result of validating an XPath expression.
type ValidationResult struct {
	XPath      string // The XPath expression that was validated
	Valid      bool   // Whether the XPath is valid
	MatchCount int    // Number of matches found
	Sample     string // Sample text from first match
	Error      string // Error message if validation failed
}

// ValidateXPath validates XPath expressions against HTML content.
func ValidateXPath(html string, xpaths []string, suppressErrors bool) []ValidationResult {
	p, err := parser.New(html)
	if err != nil {
		// Return all results as invalid
		results := make([]ValidationResult, len(xpaths))
		for i, xpath := range xpaths {
			results[i] = ValidationResult{
				XPath:  xpath,
				Valid:  false,
				Error:  "Failed to parse HTML",
			}
		}
		return results
	}

	if suppressErrors {
		p = p.WithSuppressErrors()
	}

	return ValidateXPathWithParser(p, xpaths)
}

// ValidateXPathWithParser validates XPath expressions using an existing Parser.
func ValidateXPathWithParser(p *parser.Parser, xpaths []string) []ValidationResult {
	results := make([]ValidationResult, 0, len(xpaths))

	for _, xpath := range xpaths {
		result := ValidationResult{XPath: xpath}

		nodes, err := p.XPathAll(xpath)
		if err != nil {
			result.Valid = false
			result.Error = err.Error()
		} else {
			result.Valid = true
			result.MatchCount = len(nodes)
			if len(nodes) > 0 {
				result.Sample = nodes[0].Text()
			}
		}

		results = append(results, result)
	}

	return results
}
