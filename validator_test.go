package gtmlp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateXPath_ValidHTML_ValidXPath(t *testing.T) {
	html := `<html><body>
		<h1>Title</h1>
		<p class="desc">Description</p>
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
	</body></html>`

	xpaths := map[string]string{
		"title":     "//h1",
		"desc":      "//p[@class='desc']",
		"items":     "//li",
		"container": "//ul",
	}

	results := ValidateXPath(html, xpaths)

	if len(results) != 4 {
		t.Fatalf("Expected 4 results, got: %d", len(results))
	}

	// Check title XPath
	titleResult := results["title"]
	if !titleResult.Valid {
		t.Errorf("XPath //h1 should be valid, got valid=false")
	}
	if titleResult.MatchCount != 1 {
		t.Errorf("XPath //h1 should match 1 element, got: %d", titleResult.MatchCount)
	}
	if titleResult.Error != nil {
		t.Errorf("XPath //h1 should not have error, got: %v", titleResult.Error)
	}

	// Check desc XPath
	descResult := results["desc"]
	if !descResult.Valid {
		t.Errorf("XPath //p[@class='desc'] should be valid")
	}
	if descResult.MatchCount != 1 {
		t.Errorf("XPath //p[@class='desc'] should match 1 element, got: %d", descResult.MatchCount)
	}

	// Check items XPath
	itemsResult := results["items"]
	if !itemsResult.Valid {
		t.Errorf("XPath //li should be valid")
	}
	if itemsResult.MatchCount != 2 {
		t.Errorf("XPath //li should match 2 elements, got: %d", itemsResult.MatchCount)
	}

	// Check container XPath
	containerResult := results["container"]
	if !containerResult.Valid {
		t.Errorf("XPath //ul should be valid")
	}
	if containerResult.MatchCount != 1 {
		t.Errorf("XPath //ul should match 1 element, got: %d", containerResult.MatchCount)
	}
}

func TestValidateXPath_InvalidXPathSyntax(t *testing.T) {
	html := `<html><body><h1>Title</h1></body></html>`

	xpaths := map[string]string{
		"valid":   "//h1",
		"invalid": "//[invalid",
	}

	results := ValidateXPath(html, xpaths)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got: %d", len(results))
	}

	// Valid XPath should work
	validResult := results["valid"]
	if !validResult.Valid {
		t.Errorf("XPath //h1 should be valid")
	}
	if validResult.Error != nil {
		t.Errorf("XPath //h1 should not have error, got: %v", validResult.Error)
	}

	// Invalid XPath should fail
	invalidResult := results["invalid"]
	if invalidResult.Valid {
		t.Errorf("XPath //[invalid should be invalid, got valid=true")
	}
	if invalidResult.Error == nil {
		t.Errorf("XPath //[invalid should have an error")
	}
	if invalidResult.MatchCount != 0 {
		t.Errorf("Invalid XPath should have 0 matches, got: %d", invalidResult.MatchCount)
	}
}

func TestValidateXPath_ParsingError(t *testing.T) {
	// Note: htmlquery is very lenient and doesn't error on invalid HTML.
	// This test verifies that our implementation handles empty/invalid HTML gracefully.
	// Instead of errors, we get valid results with 0 matches.
	html := ``

	xpaths := map[string]string{
		"field1": "//h1",
		"field2": "//p",
	}

	results := ValidateXPath(html, xpaths)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got: %d", len(results))
	}

	// With empty HTML, htmlquery creates an empty document (no error)
	// Results should be valid but with 0 matches
	for field, result := range results {
		if !result.Valid {
			t.Errorf("Field %s: should be valid even with empty HTML (htmlquery is lenient)", field)
		}
		if result.Error != nil {
			t.Errorf("Field %s: should not have error with empty HTML, got: %v", field, result.Error)
		}
		if result.MatchCount != 0 {
			t.Errorf("Field %s: should have 0 matches in empty HTML, got: %d", field, result.MatchCount)
		}
	}
}

func TestValidateXPath_EmptyMatchCount(t *testing.T) {
	html := `<html><body><h1>Title</h1></body></html>`

	xpaths := map[string]string{
		"found":   "//h1",
		"missing": "//nonexistent",
	}

	results := ValidateXPath(html, xpaths)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got: %d", len(results))
	}

	// Found XPath
	foundResult := results["found"]
	if !foundResult.Valid {
		t.Errorf("XPath //h1 should be valid")
	}
	if foundResult.MatchCount != 1 {
		t.Errorf("XPath //h1 should match 1 element, got: %d", foundResult.MatchCount)
	}

	// Missing XPath should still be valid but with 0 matches
	missingResult := results["missing"]
	if !missingResult.Valid {
		t.Errorf("XPath //nonexistent should be valid (no matches is still valid XPath)")
	}
	if missingResult.MatchCount != 0 {
		t.Errorf("XPath //nonexistent should match 0 elements, got: %d", missingResult.MatchCount)
	}
	if missingResult.Error != nil {
		t.Errorf("XPath //nonexistent should not have error, got: %v", missingResult.Error)
	}
}

func TestValidateXPathURL(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body>
			<div class="product">
				<h1>Product Title</h1>
				<span class="price">$19.99</span>
			</div>
		</body></html>`))
	}))
	defer server.Close()

	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
		Container: "//div[@class='product']",
		Fields: map[string]string{
			"title": "//h1",
			"price": "//span[@class='price']",
		},
	}

	results, err := ValidateXPathURL(server.URL, config)

	if err != nil {
		t.Fatalf("ValidateXPathURL failed: %v", err)
	}
	if results == nil {
		t.Fatal("Results should not be nil")
	}

	// Should have 3 results: container, title, price
	if len(results) != 3 {
		t.Fatalf("Expected 3 validation results, got: %d", len(results))
	}

	// Check container validation
	containerResult := results["container"]
	if !containerResult.Valid {
		t.Errorf("Container XPath should be valid")
	}
	if containerResult.MatchCount != 1 {
		t.Errorf("Expected 1 container match, got: %d", containerResult.MatchCount)
	}

	// Check title field
	titleResult := results["title"]
	if !titleResult.Valid {
		t.Errorf("Title XPath should be valid")
	}
	if titleResult.MatchCount != 1 {
		t.Errorf("Expected 1 title match, got: %d", titleResult.MatchCount)
	}

	// Check price field
	priceResult := results["price"]
	if !priceResult.Valid {
		t.Errorf("Price XPath should be valid")
	}
	if priceResult.MatchCount != 1 {
		t.Errorf("Expected 1 price match, got: %d", priceResult.MatchCount)
	}
}

func TestValidateXPath_EmptyXPathMap(t *testing.T) {
	html := `<html><body><h1>Title</h1></body></html>`

	xpaths := map[string]string{}

	results := ValidateXPath(html, xpaths)

	if len(results) != 0 {
		t.Fatalf("Expected 0 results for empty xpath map, got: %d", len(results))
	}
}

func TestValidateXPath_ComplexXPath(t *testing.T) {
	html := `<html><body>
		<div class="product" data-id="123">
			<h2 class="title">Product 1</h2>
			<span class="price">$10.00</span>
		</div>
		<div class="product" data-id="456">
			<h2 class="title">Product 2</h2>
			<span class="price">$20.00</span>
		</div>
	</body></html>`

	xpaths := map[string]string{
		"products":      "//div[@class='product']",
		"withAttribute": "//div[@data-id='123']",
		"nested":        "//div[@class='product']//span[@class='price']",
	}

	results := ValidateXPath(html, xpaths)

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got: %d", len(results))
	}

	// Check products
	productsResult := results["products"]
	if !productsResult.Valid {
		t.Errorf("XPath should be valid")
	}
	if productsResult.MatchCount != 2 {
		t.Errorf("Expected 2 products, got: %d", productsResult.MatchCount)
	}

	// Check with attribute
	attrResult := results["withAttribute"]
	if !attrResult.Valid {
		t.Errorf("XPath with attribute should be valid")
	}
	if attrResult.MatchCount != 1 {
		t.Errorf("Expected 1 match with data-id=123, got: %d", attrResult.MatchCount)
	}

	// Check nested
	nestedResult := results["nested"]
	if !nestedResult.Valid {
		t.Errorf("Nested XPath should be valid")
	}
	if nestedResult.MatchCount != 2 {
		t.Errorf("Expected 2 nested spans, got: %d", nestedResult.MatchCount)
	}
}
