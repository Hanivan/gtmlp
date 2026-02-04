package gtmlp

import (
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
		Fields: map[string]string{
			"name":  `.//h2/text()`,
			"price": `.//span[@class="price"]/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := Scrape[Product](testHTML, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := Scrape[Product](testHTMLEmpty, config)

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
				Fields: map[string]string{
					"name": `.//h2/text()`,
				},
				Timeout: 30 * time.Second,
			},
			wantErr: "config error: container xpath is required",
		},
		{
			name: "no fields",
			config: &Config{
				Container: `//div[@class="product"]`,
				Fields:    map[string]string{},
				Timeout:   30 * time.Second,
			},
			wantErr: "config error: at least one field is required",
		},
		{
			name: "nil fields",
			config: &Config{
				Container: `//div[@class="product"]`,
				Fields:    nil,
				Timeout:   30 * time.Second,
			},
			wantErr: "config error: at least one field is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Scrape[Product](testHTML, tt.config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := Scrape[Product]("", config)

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
		Fields: map[string]string{
			"name":  `.//h2/text()`,
			"price": `.//span[@class="price"]/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeUntyped(testHTML, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeUntyped(testHTMLEmpty, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeUntyped(testHTML, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeUntyped("", config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeUntyped(testHTML, config)

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
		Fields: map[string]string{
			"name":  `.//h2/text()`,
			"price": `.//span[@class="price"]/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := Scrape[Product](testHTML, config)

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
		Fields: map[string]string{
			"name":  `.//h2/text()`,
			"price": `.//span[@class="price"]/text()`,
			"link":  `.//a/@href`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := Scrape[ProductWithAttr](testHTML, config)

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
		Fields: map[string]string{
			"name":    `.//h2/text()`,
			"price":   `.//span[@class="price"]/text()`,
			"missing": `.//span[@class="missing"]/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeUntyped(testHTML, config)

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
		Fields: map[string]string{
			"name":  `.//h2/text()`,
			"price": `.//span[@class="price"]/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURL[Product](server.URL, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURL[Product](server.URL, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURL[Product](server.URL, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURL[Product]("invalid-url", config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURL[Product](server.URL, config)

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
		Fields: map[string]string{
			"name":  `.//h2/text()`,
			"price": `.//span[@class="price"]/text()`,
			"link":  `.//a/@href`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURL[ProductWithAttr](server.URL, config)

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
		Fields: map[string]string{
			"name":  `.//h2/text()`,
			"price": `.//span[@class="price"]/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURLUntyped(server.URL, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURLUntyped(server.URL, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURLUntyped(server.URL, config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURLUntyped("invalid-url", config)

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
		Fields: map[string]string{
			"name": `.//h2/text()`,
		},
		Timeout: 30 * time.Second,
	}

	results, err := ScrapeURLUntyped(server.URL, config)

	if err == nil {
		t.Fatal("Expected error for HTTP 500, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}
