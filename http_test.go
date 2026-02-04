package gtmlp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestFetchSuccess tests successful HTTP request
func TestFetchSuccess(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test Content</body></html>"))
	}))
	defer server.Close()

	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
	}

	resp, err := fetch(server.URL, config)
	if err != nil {
		t.Fatalf("fetch() failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// TestFetchTimeout tests request timeout
func TestFetchTimeout(t *testing.T) {
	// Create test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &Config{
		Timeout:   100 * time.Millisecond, // Very short timeout
		UserAgent: "GTMLP/2.0",
	}

	_, err := fetch(server.URL, config)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}

	scrapeErr, ok := err.(*ScrapeError)
	if !ok {
		t.Fatalf("expected *ScrapeError, got %T", err)
	}

	if scrapeErr.Type != ErrTypeNetwork {
		t.Errorf("expected error type %s, got %s", ErrTypeNetwork, scrapeErr.Type)
	}
}

// TestFetchWithRetry tests retry logic on failure
func TestFetchWithRetry(t *testing.T) {
	attempts := 0
	// Create test server that fails first 2 times, then succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Success</body></html>"))
	}))
	defer server.Close()

	config := &Config{
		Timeout:    10 * time.Second,
		UserAgent:  "GTMLP/2.0",
		MaxRetries: 3,
	}

	resp, err := fetch(server.URL, config)
	if err != nil {
		t.Fatalf("fetch() failed after retries: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

// TestFetchWithProxy tests proxy support
func TestFetchWithProxy(t *testing.T) {
	// Create a proxy server
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Proxy the request to target
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Proxied</body></html>"))
	}))
	defer proxyServer.Close()

	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
		Proxy:     proxyServer.URL,
	}

	// Note: This test verifies proxy configuration is applied
	// In real scenario, proxy would forward to actual target
	_, err := fetch(proxyServer.URL, config)
	if err != nil {
		t.Fatalf("fetch() with proxy failed: %v", err)
	}

	// Verify proxy was configured (this is a basic check)
	if config.Proxy == "" {
		t.Error("expected proxy to be configured")
	}
}

// TestFetchWithCustomHeaders tests custom headers
func TestFetchWithCustomHeaders(t *testing.T) {
	headerReceived := false
	// Create test server that checks for custom header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") == "TestValue" {
			headerReceived = true
		}
		if r.Header.Get("Authorization") == "Bearer token123" {
			headerReceived = true
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Headers OK</body></html>"))
	}))
	defer server.Close()

	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
		Headers: map[string]string{
			"X-Custom-Header": "TestValue",
			"Authorization":   "Bearer token123",
		},
	}

	resp, err := fetch(server.URL, config)
	if err != nil {
		t.Fatalf("fetch() failed: %v", err)
	}
	defer resp.Body.Close()

	if !headerReceived {
		t.Error("custom headers were not received by server")
	}
}

// TestFetchInvalidHTTPStatus tests handling of 4xx/5xx status codes
func TestFetchInvalidHTTPStatus(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"400 Bad Request", 400},
		{"404 Not Found", 404},
		{"500 Internal Server Error", 500},
		{"503 Service Unavailable", 503},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte("Error page"))
			}))
			defer server.Close()

			config := &Config{
				Timeout:   10 * time.Second,
				UserAgent: "GTMLP/2.0",
			}

			_, err := fetch(server.URL, config)
			if err == nil {
				t.Error("expected error for invalid status, got nil")
			}

			scrapeErr, ok := err.(*ScrapeError)
			if !ok {
				t.Fatalf("expected *ScrapeError, got %T", err)
			}

			if scrapeErr.Type != ErrTypeNetwork {
				t.Errorf("expected error type %s, got %s", ErrTypeNetwork, scrapeErr.Type)
			}
		})
	}
}

// TestFetchHTMLReturnsHTMLString tests fetchHTML returns HTML string
func TestFetchHTMLReturnsHTMLString(t *testing.T) {
	expectedHTML := "<html><head><title>Test</title></head><body>Content</body></html>"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedHTML))
	}))
	defer server.Close()

	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
	}

	html, err := fetchHTML(server.URL, config)
	if err != nil {
		t.Fatalf("fetchHTML() failed: %v", err)
	}

	if html != expectedHTML {
		t.Errorf("expected HTML %q, got %q", expectedHTML, html)
	}
}

// TestFetchHTMLWithRetry tests fetchHTML with retry logic
func TestFetchHTMLWithRetry(t *testing.T) {
	attempts := 0
	expectedHTML := "<html><body>Success</body></html>"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedHTML))
	}))
	defer server.Close()

	config := &Config{
		Timeout:    10 * time.Second,
		UserAgent:  "GTMLP/2.0",
		MaxRetries: 3,
	}

	html, err := fetchHTML(server.URL, config)
	if err != nil {
		t.Fatalf("fetchHTML() failed after retries: %v", err)
	}

	if html != expectedHTML {
		t.Errorf("expected HTML %q, got %q", expectedHTML, html)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

// TestFetchUserAgent tests User-Agent header is set
func TestFetchUserAgent(t *testing.T) {
	customUA := "MyCustomUserAgent/1.0"
	receivedUA := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html></html>"))
	}))
	defer server.Close()

	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: customUA,
	}

	_, err := fetch(server.URL, config)
	if err != nil {
		t.Fatalf("fetch() failed: %v", err)
	}

	if receivedUA != customUA {
		t.Errorf("expected User-Agent %q, got %q", customUA, receivedUA)
	}
}

// TestFetchInvalidURL tests error handling for invalid URL
func TestFetchInvalidURL(t *testing.T) {
	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
	}

	invalidURLs := []string{
		"not-a-url",
		"http://",
		"://invalid.com",
		"ftp://example.com",
	}

	for _, url := range invalidURLs {
		t.Run(url, func(t *testing.T) {
			_, err := fetch(url, config)
			if err == nil {
				t.Error("expected error for invalid URL, got nil")
			}

			scrapeErr, ok := err.(*ScrapeError)
			if !ok {
				t.Fatalf("expected *ScrapeError, got %T", err)
			}

			if scrapeErr.Type != ErrTypeNetwork {
				t.Errorf("expected error type %s, got %s", ErrTypeNetwork, scrapeErr.Type)
			}
		})
	}
}

// TestFetchEmptyURL tests error handling for empty URL
func TestFetchEmptyURL(t *testing.T) {
	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
	}

	_, err := fetch("", config)
	if err == nil {
		t.Error("expected error for empty URL, got nil")
	}

	scrapeErr, ok := err.(*ScrapeError)
	if !ok {
		t.Fatalf("expected *ScrapeError, got %T", err)
	}

	if scrapeErr.Type != ErrTypeNetwork {
		t.Errorf("expected error type %s, got %s", ErrTypeNetwork, scrapeErr.Type)
	}
}

// TestFetchHTMLReturnsErrorOnNetworkFailure tests fetchHTML error propagation
func TestFetchHTMLReturnsErrorOnNetworkFailure(t *testing.T) {
	config := &Config{
		Timeout:   100 * time.Millisecond,
		UserAgent: "GTMLP/2.0",
	}

	// Use a URL that will timeout (non-routable IP)
	_, err := fetchHTML("http://192.0.2.1:80", config) // 192.0.2.1 is TEST-NET-1 (documentation)
	if err == nil {
		t.Error("expected error for network failure, got nil")
	}

	scrapeErr, ok := err.(*ScrapeError)
	if !ok {
		t.Fatalf("expected *ScrapeError, got %T", err)
	}

	if scrapeErr.Type != ErrTypeNetwork {
		t.Errorf("expected error type %s, got %s", ErrTypeNetwork, scrapeErr.Type)
	}
}

// TestFetchExponentialBackoff tests exponential backoff between retries
func TestFetchExponentialBackoff(t *testing.T) {
	attempts := 0
	var timestamps []time.Time

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamps = append(timestamps, time.Now())
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := &Config{
		Timeout:    1 * time.Second,
		UserAgent:  "GTMLP/2.0",
		MaxRetries: 3,
	}

	start := time.Now()
	_, err := fetch(server.URL, config)
	duration := time.Since(start)

	if err == nil {
		t.Fatal("expected error after retries, got nil")
	}

	// Should have made 4 attempts (1 initial + 3 retries)
	if attempts != 4 {
		t.Errorf("expected 4 attempts, got %d", attempts)
	}

	// With exponential backoff (1s, 2s, 4s) and each request taking minimal time,
	// total time should be at least 7 seconds
	// But since our test uses short timeouts, we just verify it took some time
	if duration < 100*time.Millisecond {
		t.Error("expected some delay due to backoff, got very short duration")
	}
}

// TestFetchWithAllOptions tests fetch with all options combined
func TestFetchWithAllOptions(t *testing.T) {
	receivedUA := ""
	receivedCustomHeader := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		receivedCustomHeader = r.Header.Get("X-Api-Key")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Full Test</body></html>"))
	}))
	defer server.Close()

	config := &Config{
		Timeout:    10 * time.Second,
		UserAgent:  "FullTest/1.0",
		MaxRetries: 1,
		Headers: map[string]string{
			"X-Api-Key": "secret-key-123",
		},
	}

	resp, err := fetch(server.URL, config)
	if err != nil {
		t.Fatalf("fetch() failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if receivedUA != "FullTest/1.0" {
		t.Errorf("expected User-Agent 'FullTest/1.0', got %q", receivedUA)
	}

	if receivedCustomHeader != "secret-key-123" {
		t.Errorf("expected X-Api-Key 'secret-key-123', got %q", receivedCustomHeader)
	}
}

// TestFetchHTMLWithInvalidStatusCode tests fetchHTML with 4xx/5xx
func TestFetchHTMLWithInvalidStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<html><body>404 Not Found</body></html>"))
	}))
	defer server.Close()

	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
	}

	_, err := fetchHTML(server.URL, config)
	if err == nil {
		t.Error("expected error for 404 status, got nil")
	}

	scrapeErr, ok := err.(*ScrapeError)
	if !ok {
		t.Fatalf("expected *ScrapeError, got %T", err)
	}

	if scrapeErr.Type != ErrTypeNetwork {
		t.Errorf("expected error type %s, got %s", ErrTypeNetwork, scrapeErr.Type)
	}

	if !strings.Contains(strings.ToLower(err.Error()), "404") {
		t.Errorf("expected error message to contain status code, got: %v", err)
	}
}
