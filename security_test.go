package gtmlp

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestIsPrivateIP tests private IP detection
func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		isPrivate bool
	}{
		// Localhost
		{"localhost IPv4", "127.0.0.1", true},
		{"localhost IPv6", "::1", true},

		// Private IPv4 ranges
		{"10.0.0.0/8", "10.0.0.1", true},
		{"172.16.0.0/12", "172.16.0.1", true},
		{"192.168.0.0/16", "192.168.1.1", true},

		// Link-local
		{"169.254.0.0/16", "169.254.169.254", true}, // AWS metadata

		// Public IPv4
		{"public IP 1", "8.8.8.8", false},
		{"public IP 2", "1.1.1.1", false},
		{"public IP 3", "93.184.216.34", false}, // example.com

		// IPv6
		{"IPv6 public", "2001:4860:4860::8888", false}, // Google DNS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.ip)
			}

			result := isPrivateIP(ip)
			if result != tt.isPrivate {
				t.Errorf("isPrivateIP(%s) = %v, expected %v", tt.ip, result, tt.isPrivate)
			}
		})
	}
}

// TestSSRFProtection_Localhost tests SSRF protection for localhost
func TestSSRFProtection_Localhost(t *testing.T) {
	config := &Config{
		Container:       "//div",
		Fields:          map[string]FieldConfig{"name": {XPath: ".//h2"}},
		AllowPrivateIPs: false, // SSRF protection enabled
		Timeout:         30 * time.Second,
	}

	localhostURLs := []string{
		"http://localhost:8080",
		"http://127.0.0.1",
		"http://127.0.0.1:8080",
	}

	for _, url := range localhostURLs {
		err := validateURL(url, config)
		if err == nil {
			t.Errorf("Expected SSRF protection to block %s, but it was allowed", url)
		}
		if !strings.Contains(err.Error(), "SSRF protection") {
			t.Errorf("Expected SSRF error message, got: %v", err)
		}
	}
}

// TestSSRFProtection_PrivateIPs tests SSRF protection for private IP ranges
func TestSSRFProtection_PrivateIPs(t *testing.T) {
	config := &Config{
		Container:       "//div",
		Fields:          map[string]FieldConfig{"name": {XPath: ".//h2"}},
		AllowPrivateIPs: false,
		Timeout:         30 * time.Second,
	}

	privateIPs := []string{
		"http://10.0.0.1",
		"http://172.16.0.1",
		"http://192.168.1.1",
		"http://169.254.169.254", // AWS metadata service
	}

	for _, url := range privateIPs {
		err := validateURL(url, config)
		if err == nil {
			t.Errorf("Expected SSRF protection to block %s, but it was allowed", url)
		}
		if !strings.Contains(err.Error(), "SSRF protection") {
			t.Errorf("Expected SSRF error message, got: %v", err)
		}
	}
}

// TestAllowPrivateIPs tests that private IPs can be allowed when configured
func TestAllowPrivateIPs(t *testing.T) {
	config := &Config{
		Container:       "//div",
		Fields:          map[string]FieldConfig{"name": {XPath: ".//h2"}},
		AllowPrivateIPs: true, // Allow private IPs
		Timeout:         30 * time.Second,
	}

	privateIPs := []string{
		"http://127.0.0.1",
		"http://10.0.0.1",
		"http://192.168.1.1",
	}

	for _, url := range privateIPs {
		err := validateURL(url, config)
		// Should not error on private IPs when AllowPrivateIPs is true
		// Note: May still error on DNS resolution if these IPs don't exist, but shouldn't error on SSRF
		if err != nil && strings.Contains(err.Error(), "SSRF protection") {
			t.Errorf("Expected %s to be allowed with AllowPrivateIPs=true, got SSRF error: %v", url, err)
		}
	}
}

// TestCustomURLValidator tests custom URL validation
func TestCustomURLValidator(t *testing.T) {
	config := &Config{
		Container: "//div",
		Fields:    map[string]FieldConfig{"name": {XPath: ".//h2"}},
		URLValidator: func(url string) error {
			// Only allow example.com domain
			if !strings.Contains(url, "example.com") {
				return &ScrapeError{
					Type:    ErrTypeNetwork,
					Message: "domain not in allowlist",
				}
			}
			return nil
		},
		Timeout: 30 * time.Second,
	}

	// Should allow example.com
	err := validateURL("https://example.com", config)
	if err != nil {
		t.Errorf("Expected example.com to be allowed, got error: %v", err)
	}

	// Should block other domains
	err = validateURL("https://evil.com", config)
	if err == nil {
		t.Error("Expected evil.com to be blocked by custom validator")
	}
}

// TestHTTPWarning tests that HTTP URLs generate warnings
func TestHTTPWarning(t *testing.T) {
	// This test can't easily verify the warning was logged without capturing logs
	// But we can verify the URL is still validated
	config := &Config{
		Container:       "//div",
		Fields:          map[string]FieldConfig{"name": {XPath: ".//h2"}},
		AllowPrivateIPs: true, // Allow so we don't get SSRF error
		Timeout:         30 * time.Second,
	}

	// HTTP URL should not error, but should log warning
	err := validateURL("http://example.com", config)
	// Should not return error just for HTTP
	if err != nil {
		t.Errorf("HTTP URL should not cause validation error, got: %v", err)
	}

	// HTTPS should also work
	err = validateURL("https://example.com", config)
	if err != nil {
		t.Errorf("HTTPS URL should not cause validation error, got: %v", err)
	}
}

// TestSSRF_RealWorld tests SSRF protection with real HTTP requests
func TestSSRF_RealWorld(t *testing.T) {
	// Skip if in short mode (requires network)
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	config := &Config{
		Container:       "//div",
		Fields:          map[string]FieldConfig{"name": {XPath: ".//h2"}},
		AllowPrivateIPs: false,
		Timeout:         30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	// Try to scrape localhost (should be blocked)
	_, err := ScrapeURL[Product](context.Background(), "http://localhost:8080", config)
	if err == nil {
		t.Error("Expected SSRF protection to block localhost scraping")
	}
	if !strings.Contains(err.Error(), "SSRF protection") && !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("Expected SSRF or validation error, got: %v", err)
	}
}

// TestSSRF_AWSMetadata tests protection against AWS metadata service access
func TestSSRF_AWSMetadata(t *testing.T) {
	config := &Config{
		Container:       "//div",
		Fields:          map[string]FieldConfig{"name": {XPath: ".//h2"}},
		AllowPrivateIPs: false,
		Timeout:         30 * time.Second,
	}

	// AWS metadata service (common SSRF target)
	awsMetadataURLs := []string{
		"http://169.254.169.254/latest/meta-data/",
		"http://169.254.169.254/latest/user-data/",
	}

	for _, url := range awsMetadataURLs {
		err := validateURL(url, config)
		if err == nil {
			t.Errorf("Expected SSRF protection to block AWS metadata URL %s", url)
		}
		if !strings.Contains(err.Error(), "SSRF protection") {
			t.Errorf("Expected SSRF error for %s, got: %v", url, err)
		}
	}
}

// TestURLValidation_InvalidScheme tests that invalid schemes are rejected
func TestURLValidation_InvalidScheme(t *testing.T) {
	config := &Config{
		Container: "//div",
		Fields:    map[string]FieldConfig{"name": {XPath: ".//h2"}},
		Timeout:   30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	invalidURLs := []string{
		"ftp://example.com",
		"file:///etc/passwd",
		"javascript:alert(1)",
	}

	for _, url := range invalidURLs {
		_, err := ScrapeURL[Product](context.Background(), url, config)
		if err == nil {
			t.Errorf("Expected invalid scheme to be rejected for %s", url)
		}
	}
}

// TestSecurityWithPagination tests security features work with pagination
func TestSecurityWithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return HTML with a next link to localhost (SSRF attempt)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body>
			<div class="product"><h2>Product</h2></div>
			<a rel="next" href="http://127.0.0.1:8080/evil">Next</a>
		</body></html>`))
	}))
	defer server.Close()

	config := &Config{
		Container:       "//div[@class='product']",
		Fields:          map[string]FieldConfig{"name": {XPath: ".//h2/text()"}},
		AllowPrivateIPs: false, // SSRF protection enabled
		Pagination: &PaginationConfig{
			Type:         "next-link",
			NextSelector: "//a[@rel='next']/@href",
		},
		Timeout: 30 * time.Second,
	}

	type Product struct {
		Name string `json:"name"`
	}

	// Should fail when trying to follow next link to localhost
	_, err := ScrapeURL[Product](context.Background(), server.URL, config)

	// Should get pagination error with SSRF protection
	if err == nil {
		t.Error("Expected SSRF protection to block malicious pagination link")
	}

	// Check if it's a pagination error
	if pagErr, ok := err.(*PaginationError); ok {
		if !strings.Contains(pagErr.Cause.Error(), "SSRF") && !strings.Contains(pagErr.Cause.Error(), "validation") {
			t.Errorf("Expected SSRF/validation error in pagination, got: %v", pagErr.Cause)
		}
	} else if !strings.Contains(err.Error(), "SSRF") && !strings.Contains(err.Error(), "validation") {
		t.Errorf("Expected SSRF/validation error, got: %v", err)
	}
}

// TestDefaultURLValidator tests the default URL validator
func TestDefaultURLValidator(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		shouldErr bool
	}{
		{"valid HTTP", "http://example.com", false},
		{"valid HTTPS", "https://example.com", false},
		{"invalid scheme", "ftp://example.com", true},
		{"no scheme", "example.com", true},
		{"empty hostname", "http://", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := defaultURLValidator(tt.url)
			if (err != nil) != tt.shouldErr {
				t.Errorf("defaultURLValidator(%s) error = %v, shouldErr = %v", tt.url, err, tt.shouldErr)
			}
		})
	}
}
