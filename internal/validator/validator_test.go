package validator

import (
	"testing"

	"github.com/Hanivan/gtmlp/internal/parser"
)

func TestValidateXPath(t *testing.T) {
	html := `<html><body>
		<h1>Title</h1>
		<p class="desc">Description</p>
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
	</body></html>`

	xpaths := []string{
		"//h1",
		"//p[@class='desc']",
		"//li",
		"//nonexistent",
		"",
	}

	results := ValidateXPath(html, xpaths, false)

	if len(results) != 5 {
		t.Fatalf("Expected 5 results, got: %d", len(results))
	}

	// Check first XPath (//h1)
	if !results[0].Valid {
		t.Errorf("XPath //h1 should be valid")
	}
	if results[0].MatchCount != 1 {
		t.Errorf("XPath //h1 should match 1 element, got: %d", results[0].MatchCount)
	}
	if results[0].Sample != "Title" {
		t.Errorf("XPath //h1 sample should be 'Title', got: '%s'", results[0].Sample)
	}

	// Check second XPath (//p[@class='desc'])
	if !results[1].Valid {
		t.Errorf("XPath //p[@class='desc'] should be valid")
	}
	if results[1].Sample != "Description" {
		t.Errorf("XPath //p[@class='desc'] sample should be 'Description', got: '%s'", results[1].Sample)
	}

	// Check third XPath (//li)
	if !results[2].Valid {
		t.Errorf("XPath //li should be valid")
	}
	if results[2].MatchCount != 2 {
		t.Errorf("XPath //li should match 2 elements, got: %d", results[2].MatchCount)
	}

	// Check fourth XPath (//nonexistent)
	if !results[3].Valid {
		t.Errorf("XPath //nonexistent should be valid (no matches is still valid)")
	}
	if results[3].MatchCount != 0 {
		t.Errorf("XPath //nonexistent should match 0 elements, got: %d", results[3].MatchCount)
	}

	// Check fifth XPath (empty string)
	if results[4].Valid {
		t.Errorf("Empty XPath should not be valid")
	}
	if results[4].Error == "" {
		t.Errorf("Empty XPath should have an error message")
	}
}

func TestValidateXPathWithSuppressErrors(t *testing.T) {
	html := `<html><body><h1>Title</h1></body></html>`

	xpaths := []string{
		"//h1",
		"", // invalid XPath
	}

	results := ValidateXPath(html, xpaths, true)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got: %d", len(results))
	}

	// With suppressErrors, empty XPath should return as valid but with 0 matches
	if !results[1].Valid {
		t.Errorf("With suppressErrors, empty XPath should be marked as valid")
	}
}

func TestValidateXPathWithParser(t *testing.T) {
	html := `<html><body><h1>Title</h1></body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("parser.New() failed: %v", err)
	}

	xpaths := []string{"//h1"}

	results := ValidateXPathWithParser(p, xpaths)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got: %d", len(results))
	}

	if !results[0].Valid {
		t.Errorf("XPath //h1 should be valid")
	}
}
