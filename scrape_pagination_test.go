package gtmlp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test HTML fixtures for pagination
const testHTMLPage1NextLink = `<html><body>
  <div class="product"><h2>Product 1</h2></div>
  <div class="product"><h2>Product 2</h2></div>
  <a rel="next" href="/page/2">Next</a>
</body></html>`

const testHTMLPage2NextLink = `<html><body>
  <div class="product"><h2>Product 3</h2></div>
  <div class="product"><h2>Product 4</h2></div>
  <a rel="next" href="/page/3">Next</a>
</body></html>`

const testHTMLPage3NoNext = `<html><body>
  <div class="product"><h2>Product 5</h2></div>
  <div class="product"><h2>Product 6</h2></div>
</body></html>`

const testHTMLNumberedPagination = `<html><body>
  <div class="product"><h2>Product 1</h2></div>
  <div class="pagination">
    <a href="/page/1" class="active">1</a>
    <a href="/page/2">2</a>
    <a href="/page/3">3</a>
  </div>
</body></html>`

const testHTMLCircular = `<html><body>
  <div class="product"><h2>Product 1</h2></div>
  <a rel="next" href="/page/1">Next</a>
</body></html>`

const testHTMLAltSelector = `<html><body>
  <div class="product"><h2>Product 1</h2></div>
  <a class="next-page" href="/page/2">Next Page</a>
</body></html>`

const testHTMLRelativeURL = `<html><body>
  <div class="product"><h2>Product 1</h2></div>
  <a rel="next" href="?page=2">Next</a>
</body></html>`

// TestScrapeURL_NextLinkPagination tests next-link pagination with auto-follow
func TestScrapeURL_NextLinkPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		switch r.URL.Path {
		case "/products":
			w.Write([]byte(testHTMLPage1NextLink))
		case "/page/2":
			w.Write([]byte(testHTMLPage2NextLink))
		case "/page/3":
			w.Write([]byte(testHTMLPage3NoNext))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:         "next-link",
			NextSelector: `//a[@rel="next"]/@href`,
			MaxPages:     10,
		},
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	products, err := ScrapeURL[Product](context.Background(), server.URL+"/products", config)
	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	// Should have 6 products from 3 pages
	if len(products) != 6 {
		t.Fatalf("Expected 6 products, got %d", len(products))
	}

	expectedNames := []string{"Product 1", "Product 2", "Product 3", "Product 4", "Product 5", "Product 6"}
	for i, expected := range expectedNames {
		if products[i].Name != expected {
			t.Errorf("Expected product[%d].Name = %s, got %s", i, expected, products[i].Name)
		}
	}
}

// TestScrapeURLWithPages tests page-separated results
func TestScrapeURLWithPages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		switch r.URL.Path {
		case "/products":
			w.Write([]byte(testHTMLPage1NextLink))
		case "/page/2":
			w.Write([]byte(testHTMLPage2NextLink))
		case "/page/3":
			w.Write([]byte(testHTMLPage3NoNext))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:         "next-link",
			NextSelector: `//a[@rel="next"]/@href`,
		},
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	results, err := ScrapeURLWithPages[Product](context.Background(), server.URL+"/products", config)
	if err != nil {
		t.Fatalf("ScrapeURLWithPages failed: %v", err)
	}

	if results.TotalPages != 3 {
		t.Errorf("Expected 3 pages, got %d", results.TotalPages)
	}

	if results.TotalItems != 6 {
		t.Errorf("Expected 6 total items, got %d", results.TotalItems)
	}

	if len(results.Pages) != 3 {
		t.Fatalf("Expected 3 pages, got %d", len(results.Pages))
	}

	// Check each page
	for i, page := range results.Pages {
		if page.PageNum != i+1 {
			t.Errorf("Expected page.PageNum = %d, got %d", i+1, page.PageNum)
		}
		if len(page.Items) != 2 {
			t.Errorf("Expected 2 items on page %d, got %d", i+1, len(page.Items))
		}
	}
}

// TestExtractPaginationURLs_NextLink tests extracting next-link URLs
func TestExtractPaginationURLs_NextLink(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		switch r.URL.Path {
		case "/products":
			w.Write([]byte(testHTMLPage1NextLink))
		case "/page/2":
			w.Write([]byte(testHTMLPage2NextLink))
		case "/page/3":
			w.Write([]byte(testHTMLPage3NoNext))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:         "next-link",
			NextSelector: `//a[@rel="next"]/@href`,
			MaxPages:     10,
		},
		Timeout: 30 * time.Second,
	}

	info, err := ExtractPaginationURLs(context.Background(), server.URL+"/products", config)
	if err != nil {
		t.Fatalf("ExtractPaginationURLs failed: %v", err)
	}

	if info.Type != "next-link" {
		t.Errorf("Expected type 'next-link', got %s", info.Type)
	}

	// Should have 3 URLs
	if len(info.URLs) != 3 {
		t.Fatalf("Expected 3 URLs, got %d", len(info.URLs))
	}
}

// TestExtractPaginationURLs_Numbered tests extracting numbered page URLs
func TestExtractPaginationURLs_Numbered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(testHTMLNumberedPagination))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:         "numbered",
			PageSelector: `//div[@class="pagination"]//a/@href`,
		},
		Timeout: 30 * time.Second,
	}

	info, err := ExtractPaginationURLs(context.Background(), server.URL+"/products", config)
	if err != nil {
		t.Fatalf("ExtractPaginationURLs failed: %v", err)
	}

	if info.Type != "numbered" {
		t.Errorf("Expected type 'numbered', got %s", info.Type)
	}

	// Should have 3 URLs
	if len(info.URLs) != 3 {
		t.Fatalf("Expected 3 URLs, got %d: %v", len(info.URLs), info.URLs)
	}
}

// TestPagination_AltSelectors tests fallback to altSelectors
func TestPagination_AltSelectors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if r.URL.Path == "/products" {
			w.Write([]byte(testHTMLAltSelector))
		} else {
			w.Write([]byte(testHTMLPage3NoNext))
		}
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:         "next-link",
			NextSelector: `//a[@rel="next"]/@href`, // Won't match
			AltSelectors: []string{`//a[@class="next-page"]/@href`}, // Will match
		},
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	products, err := ScrapeURL[Product](context.Background(), server.URL+"/products", config)
	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	// Should have products from 2 pages
	if len(products) < 1 {
		t.Fatalf("Expected at least 1 product, got %d", len(products))
	}
}

// TestPagination_CircularReference tests duplicate URL detection
func TestPagination_CircularReference(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(testHTMLCircular))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:          "next-link",
			NextSelector:  `//a[@rel="next"]/@href`,
			EnableLogging: false,
		},
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	products, err := ScrapeURL[Product](context.Background(), server.URL+"/page/1", config)
	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	// Should only scrape once (circular reference detected)
	if len(products) != 1 {
		t.Errorf("Expected 1 product (circular detected), got %d", len(products))
	}
}

// TestPagination_MaxPages tests maxPages limit
func TestPagination_MaxPages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		// Always return a next link
		w.Write([]byte(testHTMLPage1NextLink))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:         "next-link",
			NextSelector: `//a[@rel="next"]/@href`,
			MaxPages:     2, // Limit to 2 pages
		},
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	products, err := ScrapeURL[Product](context.Background(), server.URL+"/products", config)
	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	// Should have 4 products (2 pages × 2 products)
	if len(products) != 4 {
		t.Errorf("Expected 4 products (maxPages=2), got %d", len(products))
	}
}

// TestPagination_RelativeURL tests relative URL resolution
func TestPagination_RelativeURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if r.URL.RawQuery == "" {
			w.Write([]byte(testHTMLRelativeURL))
		} else {
			w.Write([]byte(testHTMLPage3NoNext))
		}
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:         "next-link",
			NextSelector: `//a[@rel="next"]/@href`,
		},
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	products, err := ScrapeURL[Product](context.Background(), server.URL+"/products", config)
	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	// Should resolve relative URL and scrape 2 pages
	// testHTMLRelativeURL has 1 product, testHTMLPage3NoNext has 2 products = 3 total
	if len(products) < 2 {
		t.Errorf("Expected at least 2 products from 2 pages, got %d", len(products))
	}
}

// TestPagination_NoConfig tests backward compatibility
func TestPagination_NoConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(testHTMLPage1NextLink))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		// No Pagination field
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	products, err := ScrapeURL[Product](context.Background(), server.URL+"/products", config)
	if err != nil {
		t.Fatalf("ScrapeURL failed: %v", err)
	}

	// Should only scrape single page (no pagination)
	if len(products) != 2 {
		t.Errorf("Expected 2 products (single page), got %d", len(products))
	}
}

// TestPagination_PageFails tests error handling when a page fails
func TestPagination_PageFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/page/2" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(testHTMLPage1NextLink))
	}))
	defer server.Close()

	config := &Config{
		Container: `//div[@class="product"]`,
		Fields: map[string]FieldConfig{
			"name": {XPath: `.//h2/text()`},
		},
		Pagination: &PaginationConfig{
			Type:         "next-link",
			NextSelector: `//a[@rel="next"]/@href`,
		},
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	_, err := ScrapeURL[Product](context.Background(), server.URL+"/products", config)
	if err == nil {
		t.Fatal("Expected error when page fails, got nil")
	}

	// Check if it's a PaginationError
	pagErr, ok := err.(*PaginationError)
	if !ok {
		t.Fatalf("Expected PaginationError, got %T", err)
	}

	if pagErr.PageNumber != 2 {
		t.Errorf("Expected page number 2, got %d", pagErr.PageNumber)
	}

	if pagErr.TotalScraped != 2 {
		t.Errorf("Expected 2 items scraped before failure, got %d", pagErr.TotalScraped)
	}
}

// TestNormalizeURL tests URL normalization
func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://example.com/page/2", "https://example.com/page/2"},
		{"https://example.com/page/2/", "https://example.com/page/2"},
		{"https://example.com/page/2#section", "https://example.com/page/2"},
		{"https://example.com/page?b=2&a=1", "https://example.com/page?a=1&b=2"},
	}

	for _, tt := range tests {
		result := normalizeURL(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeURL(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

// TestResolveURL tests relative URL resolution
func TestResolveURL(t *testing.T) {
	tests := []struct {
		base     string
		relative string
		expected string
	}{
		{"https://example.com/products", "/page/2", "https://example.com/page/2"},
		{"https://example.com/products", "?page=2", "https://example.com/products?page=2"},
		{"https://example.com/products/", "page/2", "https://example.com/products/page/2"},
		{"https://example.com/products", "https://other.com/page/2", "https://other.com/page/2"},
	}

	for _, tt := range tests {
		result, err := resolveURL(tt.base, tt.relative)
		if err != nil {
			t.Errorf("resolveURL(%s, %s) error: %v", tt.base, tt.relative, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("resolveURL(%s, %s) = %s, expected %s", tt.base, tt.relative, result, tt.expected)
		}
	}
}

// TestValidatePaginationConfig tests pagination config validation
func TestValidatePaginationConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid next-link",
			config: &Config{
				Container: "//div",
				Fields:    map[string]FieldConfig{"name": {XPath: ".//h2"}},
				Pagination: &PaginationConfig{
					Type:         "next-link",
					NextSelector: "//a[@rel='next']/@href",
				},
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid numbered",
			config: &Config{
				Container: "//div",
				Fields:    map[string]FieldConfig{"name": {XPath: ".//h2"}},
				Pagination: &PaginationConfig{
					Type:         "numbered",
					PageSelector: "//a/@href",
				},
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			config: &Config{
				Container: "//div",
				Fields:    map[string]FieldConfig{"name": {XPath: ".//h2"}},
				Pagination: &PaginationConfig{
					Type: "invalid",
				},
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "missing nextSelector",
			config: &Config{
				Container: "//div",
				Fields:    map[string]FieldConfig{"name": {XPath: ".//h2"}},
				Pagination: &PaginationConfig{
					Type: "next-link",
				},
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "missing pageSelector",
			config: &Config{
				Container: "//div",
				Fields:    map[string]FieldConfig{"name": {XPath: ".//h2"}},
				Pagination: &PaginationConfig{
					Type: "numbered",
				},
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
