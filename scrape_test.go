package gtmlp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test types for scraping
type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type ProductWithAttr struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

// Test HTML with products
const testHTML = `<html><body>
  <div class="product">
    <h2>Product 1</h2>
    <span class="price">$10</span>
    <a href="/product/1">View</a>
  </div>
  <div class="product">
    <h2>Product 2</h2>
    <span class="price">$20</span>
    <a href="/product/2">View</a>
  </div>
</body></html>`

const testHTMLEmpty = `<html><body></body></html>`

// TestScrapeTyped_Success tests successful scraping with typed result
func TestScrapeTyped_Success(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name":  {XPath: `.//h2/text()`},
			"price": {XPath: `.//span[@class="price"]/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := Scrape[Product](context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Check first product
	if results[0].Name != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%s'", results[0].Name)
	}
	if results[0].Price != "$10" {
		t.Errorf("Expected price '$10', got '%s'", results[0].Price)
	}

	// Check second product
	if results[1].Name != "Product 2" {
		t.Errorf("Expected name 'Product 2', got '%s'", results[1].Name)
	}
	if results[1].Price != "$20" {
		t.Errorf("Expected price '$20', got '%s'", results[1].Price)
	}
}

// TestScrapeTyped_EmptyResults tests scraping when no containers found
func TestScrapeTyped_EmptyResults(t *testing.T) {
	config := &Config{
		Container: `//div[@class="nonexistent"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := Scrape[Product](context.Background(), testHTMLEmpty, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 results, got %d", len(results))
	}
}

// TestScrapeTyped_InvalidConfig tests scraping with invalid config
func TestScrapeTyped_InvalidConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr string
	}{
		{
			name: "empty container",
			config: &Config{
				Container: "",
				Fields: map[string]FieldConfig{
					"name": {XPath: `.//h2/text()`},
				},
				Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
			},
			wantErr: "config error: container xpath is required",
		},
		{
			name: "no fields",
			config: &Config{
				Container: `//div[@class="product"]`,
				Fields:    map[string]FieldConfig{},
				Timeout:   30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
			},
			wantErr: "config error: at least one field is required",
		},
		{
			name: "nil fields",
			config: &Config{
				Container: `//div[@class="product"]`,
				Fields:    nil,
				Timeout:   30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
			},
			wantErr: "config error: at least one field is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Scrape[Product](context.Background(), testHTML, tt.config)

			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if err.Error() != tt.wantErr {
				t.Errorf("Expected error message '%s', got '%s'", tt.wantErr, err.Error())
			}

			if results != nil {
				t.Errorf("Expected nil results, got %v", results)
			}
		})
	}
}

// TestScrapeTyped_EmptyHTML tests scraping with empty HTML
func TestScrapeTyped_EmptyHTML(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := Scrape[Product](context.Background(), "", config)

	// Should not error, just return empty results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestScrapeUntyped_Success tests successful scraping with untyped result
func TestScrapeUntyped_Success(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name":  {XPath: `.//h2/text()`},
			"price": {XPath: `.//span[@class="price"]/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("ScrapeUntyped failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Check first result
	if results[0]["name"] != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%v'", results[0]["name"])
	}
	if results[0]["price"] != "$10" {
		t.Errorf("Expected price '$10', got '%v'", results[0]["price"])
	}

	// Check second result
	if results[1]["name"] != "Product 2" {
		t.Errorf("Expected name 'Product 2', got '%v'", results[1]["name"])
	}
	if results[1]["price"] != "$20" {
		t.Errorf("Expected price '$20', got '%v'", results[1]["price"])
	}
}

// TestScrapeUntyped_EmptyResults tests untyped scraping when no containers found
func TestScrapeUntyped_EmptyResults(t *testing.T) {
	config := &Config{
		Container: `//div[@class="nonexistent"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTMLEmpty, config)

	if err != nil {
		t.Fatalf("ScrapeUntype failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 results, got %d", len(results))
	}
}

// TestScrapeUntyped_InvalidConfig tests untyped scraping with invalid config
func TestScrapeUntyped_InvalidConfig(t *testing.T) {
	config := &Config{
		Container: "",
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTML, config)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

// TestScrapeUntyped_EmptyHTML tests untyped scraping with empty HTML
func TestScrapeUntyped_EmptyHTML(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), "", config)

	// Should not error, just return empty results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestExtractValue_Text tests extracting text content
func TestExtractValue_Text(t *testing.T) {
	// This test will be implemented once we have the extractValue helper
	// For now, we'll test it indirectly through Scrape
	config := &Config{
		Container: `//div[@class="product"][1]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("ScrapeUntyped failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0]["name"] != "Product 1" {
		t.Errorf("Expected 'Product 1', got '%v'", results[0]["name"])
	}
}

// TestMapToStruct tests converting map to struct
func TestMapToStruct(t *testing.T) {
	// This test will be implemented once we have the mapToStruct helper
	// For now, we'll test it indirectly through Scrape
	config := &Config{
		Container: `//div[@class="product"][1]`,
		Fields: map[string]FieldConfig{
			"name":  {XPath: `.//h2/text()`},
			"price": {XPath: `.//span[@class="price"]/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := Scrape[Product](context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Name != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%s'", results[0].Name)
	}
	if results[0].Price != "$10" {
		t.Errorf("Expected price '$10', got '%s'", results[0].Price)
	}
}

// TestScrape_WithAttributes tests scraping with attribute extraction
func TestScrape_WithAttributes(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name":  {XPath: `.//h2/text()`},
			"price": {XPath: `.//span[@class="price"]/text()`},
			"link":  {XPath: `.//a/@href`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := Scrape[ProductWithAttr](context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Check first product
	if results[0].Name != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%s'", results[0].Name)
	}
	if results[0].Price != "$10" {
		t.Errorf("Expected price '$10', got '%s'", results[0].Price)
	}
	if results[0].Link != "/product/1" {
		t.Errorf("Expected link '/product/1', got '%s'", results[0].Link)
	}

	// Check second product
	if results[1].Link != "/product/2" {
		t.Errorf("Expected link '/product/2', got '%s'", results[1].Link)
	}
}

// TestScrape_MissingFields tests scraping when fields are missing
func TestScrape_MissingFields(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name":    {XPath: `.//h2/text()`},
			"price":   {XPath: `.//span[@class="price"]/text()`},
			"missing": {XPath: `.//span[@class="missing"]/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("ScrapeUntyped failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Missing field should be nil or empty string
	if val, ok := results[0]["missing"]; ok && val != "" {
		t.Errorf("Expected missing field to be empty, got '%v'", val)
	}
}

// TestScrapeURL_Success tests successful URL scraping with typed result
func TestScrapeURL_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name":  {XPath: `.//h2/text()`},
			"price": {XPath: `.//span[@class="price"]/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURL[Product](context.Background(), server.URL, config)

	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Check first product
	if results[0].Name != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%s'", results[0].Name)
	}
	if results[0].Price != "$10" {
		t.Errorf("Expected price '$10', got '%s'", results[0].Price)
	}

	// Check second product
	if results[1].Name != "Product 2" {
		t.Errorf("Expected name 'Product 2', got '%s'", results[1].Name)
	}
	if results[1].Price != "$20" {
		t.Errorf("Expected price '$20', got '%s'", results[1].Price)
	}
}

// TestScrapeURL_EmptyResults tests URL scraping when no containers found
func TestScrapeURL_EmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTMLEmpty))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="nonexistent"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURL[Product](context.Background(), server.URL, config)

	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 results, got %d", len(results))
	}
}

// TestScrapeURL_InvalidConfig tests URL scraping with invalid config
func TestScrapeURL_InvalidConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	config := &Config{
		Container: "",
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURL[Product](context.Background(), server.URL, config)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

// TestScrapeURL_InvalidURL tests URL scraping with invalid URL
func TestScrapeURL_InvalidURL(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURL[Product](context.Background(), "invalid-url", config)

	if err == nil {
		t.Fatal("Expected error for invalid URL, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

// TestScrapeURL_HTTPError tests URL scraping when HTTP request fails
func TestScrapeURL_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURL[Product](context.Background(), server.URL, config)

	if err == nil {
		t.Fatal("Expected error for HTTP 500, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

// TestScrapeURL_WithAttributes tests URL scraping with attribute extraction
func TestScrapeURL_WithAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name":  {XPath: `.//h2/text()`},
			"price": {XPath: `.//span[@class="price"]/text()`},
			"link":  {XPath: `.//a/@href`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURL[ProductWithAttr](context.Background(), server.URL, config)

	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Check first product
	if results[0].Name != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%s'", results[0].Name)
	}
	if results[0].Link != "/product/1" {
		t.Errorf("Expected link '/product/1', got '%s'", results[0].Link)
	}
}

// TestScrapeURLUntyped_Success tests successful URL scraping with untyped result
func TestScrapeURLUntyped_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name":  {XPath: `.//h2/text()`},
			"price": {XPath: `.//span[@class="price"]/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURLUntyped(context.Background(), server.URL, config)

	if err != nil {
		t.Fatalf("ScrapeURLUntyped failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Check first result
	if results[0]["name"] != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%v'", results[0]["name"])
	}
	if results[0]["price"] != "$10" {
		t.Errorf("Expected price '$10', got '%v'", results[0]["price"])
	}

	// Check second result
	if results[1]["name"] != "Product 2" {
		t.Errorf("Expected name 'Product 2', got '%v'", results[1]["name"])
	}
	if results[1]["price"] != "$20" {
		t.Errorf("Expected price '$20', got '%v'", results[1]["price"])
	}
}

// TestScrapeURLUntyped_EmptyResults tests untyped URL scraping when no containers found
func TestScrapeURLUntyped_EmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTMLEmpty))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="nonexistent"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURLUntyped(context.Background(), server.URL, config)

	if err != nil {
		t.Fatalf("ScrapeURLUntyped failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 results, got %d", len(results))
	}
}

// TestScrapeURLUntyped_InvalidConfig tests untyped URL scraping with invalid config
func TestScrapeURLUntyped_InvalidConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	config := &Config{
		Container: "",
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURLUntyped(context.Background(), server.URL, config)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

// TestScrapeURLUntyped_InvalidURL tests untyped URL scraping with invalid URL
func TestScrapeURLUntyped_InvalidURL(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURLUntyped(context.Background(), "invalid-url", config)

	if err == nil {
		t.Fatal("Expected error for invalid URL, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

// TestScrapeURLUntyped_HTTPError tests untyped URL scraping when HTTP request fails
func TestScrapeURLUntyped_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeURLUntyped(context.Background(), server.URL, config)

	if err == nil {
		t.Fatal("Expected error for HTTP 500, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

// Test HTML with varying structures for altXpath tests
const testHTMLAltStructure = `<html><body>
  <div class="product">
    <h1>Product A</h1>
    <span class="price">$15</span>
  </div>
  <div class="item">
    <h2>Product B</h2>
    <div class="price">$25</div>
  </div>
  <article class="listing">
    <h2>Product C</h2>
    <span class="cost">$35</span>
  </article>
</body></html>`

const testHTMLWhitespace = `<html><body>
  <div class="product">
    <h2>   </h2>
    <span class="alt-name">Product With Fallback</span>
    <span class="price">$50</span>
  </div>
</body></html>`

// TestAltContainer_Fallback tests container fallback when primary is empty
func TestAltContainer_Fallback(t *testing.T) {
	config := &Config{
		Container:    `//div[@class="nonexistent"]`,
		AltContainer: []string{`//div[@class="item"]`, `//article[@class="listing"]`},
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTMLAltStructure, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result (from first altContainer), got %d", len(results))
	}

	if results[0]["name"] != "Product B" {
		t.Errorf("Expected name 'Product B', got '%v'", results[0]["name"])
	}
}

// TestAltContainer_AllFail tests when all containers fail
func TestAltContainer_AllFail(t *testing.T) {
	config := &Config{
		Container:    `//div[@class="missing"]`,
		AltContainer: []string{`//div[@class="also-missing"]`, `//article[@class="nope"]`},
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTMLAltStructure, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 results when all containers fail, got %d", len(results))
	}
}

// TestAltXPath_Fallback tests field XPath fallback when primary is empty
func TestAltXPath_Fallback(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {
				XPath:    `.//h3/text()`,
				AltXPath: []string{`.//h2/text()`, `.//h1/text()`},
			},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTMLAltStructure, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	// Primary h3 doesn't exist, should fallback to h2 -> empty, then h1 -> "Product A"
	if results[0]["name"] != "Product A" {
		t.Errorf("Expected name 'Product A' from altXpath fallback, got '%v'", results[0]["name"])
	}
}

// TestAltXPath_WithPipes tests altXpath with pipe validation
func TestAltXPath_WithPipes(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {
				XPath:    `.//h2/text()`,
				AltXPath: []string{`.//span[@class="alt-name"]/text()`},
				Pipes:    []string{"trim"},
			},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTMLWhitespace, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	// Primary h2 has only whitespace "   " -> trim -> "" -> fallback to alt-name
	if results[0]["name"] != "Product With Fallback" {
		t.Errorf("Expected 'Product With Fallback' after whitespace fallback, got '%v'", results[0]["name"])
	}
}

// TestAltXPath_AllFail tests when all XPaths fail
func TestAltXPath_AllFail(t *testing.T) {
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {
				XPath:    `.//h5/text()`,
				AltXPath: []string{`.//h6/text()`, `.//h7/text()`},
			},
			"price": {XPath: `.//span[@class="price"]/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// All XPaths for 'name' fail, should return empty string
	if results[0]["name"] != "" {
		t.Errorf("Expected empty name when all XPaths fail, got '%v'", results[0]["name"])
	}

	// Price should still work
	if results[0]["price"] != "$10" {
		t.Errorf("Expected price '$10', got '%v'", results[0]["price"])
	}
}

// TestAltXPath_MultipleAlternatives tests multiple alternatives with first match
func TestAltXPath_MultipleAlternatives(t *testing.T) {
	config := &Config{
		Container: `//article[@class="listing"]`,
		Fields: map[string]FieldConfig{
			"name": {
				XPath:    `.//h3/text()`,
				AltXPath: []string{`.//h1/text()`, `.//h2/text()`, `.//span/text()`},
			},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTMLAltStructure, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	// h3 empty -> h1 empty -> h2 has "Product C" -> should use h2
	if results[0]["name"] != "Product C" {
		t.Errorf("Expected 'Product C' from second altXpath, got '%v'", results[0]["name"])
	}
}

// TestAltXPath_BackwardCompat tests configs without altXpath work unchanged
func TestAltXPath_BackwardCompat(t *testing.T) {
	// Config without any alt fields
	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name":  {XPath: `.//h2/text()`},
			"price": {XPath: `.//span[@class="price"]/text()`},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := Scrape[Product](context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0].Name != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%s'", results[0].Name)
	}
	if results[0].Price != "$10" {
		t.Errorf("Expected price '$10', got '%s'", results[0].Price)
	}
}

// TestAltXPath_EmptyArrays tests empty altXpath/altContainer arrays
func TestAltXPath_EmptyArrays(t *testing.T) {
	config := &Config{
		Container:    `//div[@class="product"]`,
		AltContainer: []string{}, // Empty array
		Fields: map[string]FieldConfig{
			"name": {
				XPath:    `.//h2/text()`,
				AltXPath: []string{}, // Empty array
			},
		},
		Timeout: 30 * time.Second,
		AllowPrivateIPs: true, // Allow localhost for testing
	}

	results, err := ScrapeUntyped(context.Background(), testHTML, config)

	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0]["name"] != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%v'", results[0]["name"])
	}
}

