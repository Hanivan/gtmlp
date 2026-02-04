package gtmlp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCheckHealth_HealthyURL tests that a healthy URL (2xx status) returns StatusHealthy
func TestCheckHealth_HealthyURL(t *testing.T) {
	// Create a test server that returns 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	result := CheckHealth(server.URL)

	// Verify result
	if result.URL != server.URL {
		t.Errorf("Expected URL %s, got %s", server.URL, result.URL)
	}

	if result.Status != StatusHealthy {
		t.Errorf("Expected status %v, got %v", StatusHealthy, result.Status)
	}

	if result.Code != http.StatusOK {
		t.Errorf("Expected code %d, got %d", http.StatusOK, result.Code)
	}

	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}

	if result.Latency == 0 {
		t.Error("Expected non-zero latency")
	}
}

// TestCheckHealth_UnhealthyURL_4xx tests that a 4xx status returns StatusUnhealthy
func TestCheckHealth_UnhealthyURL_4xx(t *testing.T) {
	// Create a test server that returns 404 Not Found
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	result := CheckHealth(server.URL)

	// Verify result
	if result.Status != StatusUnhealthy {
		t.Errorf("Expected status %v, got %v", StatusUnhealthy, result.Status)
	}

	if result.Code != http.StatusNotFound {
		t.Errorf("Expected code %d, got %d", http.StatusNotFound, result.Code)
	}

	if result.Error != nil {
		t.Errorf("Expected no error for 4xx, got %v", result.Error)
	}

	if result.Latency == 0 {
		t.Error("Expected non-zero latency")
	}
}

// TestCheckHealth_UnhealthyURL_5xx tests that a 5xx status returns StatusUnhealthy
func TestCheckHealth_UnhealthyURL_5xx(t *testing.T) {
	// Create a test server that returns 500 Internal Server Error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	result := CheckHealth(server.URL)

	// Verify result
	if result.Status != StatusUnhealthy {
		t.Errorf("Expected status %v, got %v", StatusUnhealthy, result.Status)
	}

	if result.Code != http.StatusInternalServerError {
		t.Errorf("Expected code %d, got %d", http.StatusInternalServerError, result.Code)
	}

	if result.Error != nil {
		t.Errorf("Expected no error for 5xx, got %v", result.Error)
	}
}

// TestCheckHealth_NetworkError tests that a network failure returns StatusError
func TestCheckHealth_NetworkError(t *testing.T) {
	// Use an invalid URL that will cause a network error
	invalidURL := "http://localhost:99999/nonexistent"

	result := CheckHealth(invalidURL)

	// Verify result
	if result.Status != StatusError {
		t.Errorf("Expected status %v, got %v", StatusError, result.Status)
	}

	if result.Code != 0 {
		t.Errorf("Expected code 0 for network error, got %d", result.Code)
	}

	if result.Error == nil {
		t.Error("Expected error for network failure")
	}
}

// TestCheckHealth_Timeout tests that a timeout returns StatusError
func TestCheckHealth_Timeout(t *testing.T) {
	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create config with short timeout
	config := &Config{
		Timeout: 100 * time.Millisecond,
	}

	result := CheckHealthWithOptions(server.URL, config)

	// Verify result
	if result.Status != StatusError {
		t.Errorf("Expected status %v for timeout, got %v", StatusError, result.Status)
	}

	if result.Error == nil {
		t.Error("Expected error for timeout")
	}
}

// TestCheckHealth_LatencyMeasurement tests that latency is measured correctly
func TestCheckHealth_LatencyMeasurement(t *testing.T) {
	// Create a test server with a small delay
	delay := 100 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := CheckHealth(server.URL)

	// Verify latency was measured
	if result.Latency == 0 {
		t.Error("Expected non-zero latency")
	}

	// Latency should be at least the delay (with some tolerance for processing)
	if result.Latency < delay {
		t.Errorf("Expected latency >= %v, got %v", delay, result.Latency)
	}
}

// TestCheckHealthMulti_MultipleURLs tests checking multiple URLs concurrently
func TestCheckHealthMulti_MultipleURLs(t *testing.T) {
	// Create test servers with different responses
	healthyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer healthyServer.Close()

	unhealthyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer unhealthyServer.Close()

	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	urls := []string{
		healthyServer.URL,
		unhealthyServer.URL,
		errorServer.URL,
		"http://localhost:99999/invalid",
	}

	results := CheckHealthMulti(urls)

	// Verify number of results
	if len(results) != len(urls) {
		t.Fatalf("Expected %d results, got %d", len(urls), len(results))
	}

	// Verify healthy URL
	if results[0].Status != StatusHealthy {
		t.Errorf("Expected URL 0 to be %v, got %v", StatusHealthy, results[0].Status)
	}

	// Verify unhealthy URL (4xx)
	if results[1].Status != StatusUnhealthy {
		t.Errorf("Expected URL 1 to be %v, got %v", StatusUnhealthy, results[1].Status)
	}

	// Verify unhealthy URL (5xx)
	if results[2].Status != StatusUnhealthy {
		t.Errorf("Expected URL 2 to be %v, got %v", StatusUnhealthy, results[2].Status)
	}

	// Verify error URL
	if results[3].Status != StatusError {
		t.Errorf("Expected URL 3 to be %v, got %v", StatusError, results[3].Status)
	}
}

// TestCheckHealthMulti_EmptyList tests that empty URL list returns empty results
func TestCheckHealthMulti_EmptyList(t *testing.T) {
	urls := []string{}

	results := CheckHealthMulti(urls)

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty URL list, got %d", len(results))
	}
}

// TestCheckHealthMulti_Concurrency tests that concurrent checks work correctly
func TestCheckHealthMulti_Concurrency(t *testing.T) {
	// Create multiple slow servers
	numServers := 5
	servers := make([]*httptest.Server, numServers)
	urls := make([]string, numServers)

	for i := 0; i < numServers; i++ {
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		urls[i] = servers[i].URL
		defer servers[i].Close()
	}

	// Measure time taken
	start := time.Now()
	results := CheckHealthMulti(urls)
	elapsed := time.Since(start)

	// Verify results
	if len(results) != numServers {
		t.Fatalf("Expected %d results, got %d", numServers, len(results))
	}

	// All should be healthy
	for i, result := range results {
		if result.Status != StatusHealthy {
			t.Errorf("Expected URL %d to be %v, got %v", i, StatusHealthy, result.Status)
		}
	}

	// With concurrency, should be much faster than sequential (5 * 50ms = 250ms)
	// Allow some overhead, but should be significantly less
	if elapsed > 200*time.Millisecond {
		t.Errorf("Expected concurrent execution to be fast, took %v", elapsed)
	}
}

// TestCheckHealthWithOptions_CustomTimeout tests custom timeout configuration
func TestCheckHealthWithOptions_CustomTimeout(t *testing.T) {
	// Create a server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Test with sufficient timeout
	config := &Config{
		Timeout: 200 * time.Millisecond,
	}

	result := CheckHealthWithOptions(server.URL, config)

	if result.Status != StatusHealthy {
		t.Errorf("Expected status %v with sufficient timeout, got %v", StatusHealthy, result.Status)
	}

	if result.Error != nil {
		t.Errorf("Expected no error with sufficient timeout, got %v", result.Error)
	}
}

// TestCheckHealthWithOptions_CustomUserAgent tests custom user agent
func TestCheckHealthWithOptions_CustomUserAgent(t *testing.T) {
	receivedUA := ""

	// Create a test server that checks User-Agent
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	customUA := "TestAgent/1.0"
	config := &Config{
		Timeout:   5 * time.Second,
		UserAgent: customUA,
	}

	result := CheckHealthWithOptions(server.URL, config)

	if result.Status != StatusHealthy {
		t.Errorf("Expected status %v, got %v", StatusHealthy, result.Status)
	}

	if receivedUA != customUA {
		t.Errorf("Expected User-Agent %s, got %s", customUA, receivedUA)
	}
}

// TestCheckHealthWithOptions_Proxy tests proxy configuration
func TestCheckHealthWithOptions_Proxy(t *testing.T) {
	// This test validates that proxy config is accepted
	// Note: We can't test actual proxy functionality without a proxy server

	config := &Config{
		Timeout: 5 * time.Second,
		Proxy:   "http://proxy.example.com:8080",
	}

	// Create a simple test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// This should not panic or error due to proxy config format
	result := CheckHealthWithOptions(server.URL, config)

	// Result will be error since proxy doesn't exist, but we're just testing
	// that the function accepts the proxy config
	if result.Error == nil {
		t.Log("Note: Proxy test passed - config was accepted")
	}
}

// TestCheckHealth_InvalidURL tests invalid URL format
func TestCheckHealth_InvalidURL(t *testing.T) {
	invalidURLs := []string{
		"",
		"not-a-url",
		"ftp://example.com",
		"://invalid",
	}

	for _, url := range invalidURLs {
		t.Run(fmt.Sprintf("URL_%q", url), func(t *testing.T) {
			result := CheckHealth(url)

			if result.Status != StatusError {
				t.Errorf("Expected status %v for invalid URL, got %v", StatusError, result.Status)
			}

			if result.Error == nil {
				t.Error("Expected error for invalid URL")
			}
		})
	}
}

// TestHealthStatus_String tests HealthStatus string representation
func TestHealthStatus_String(t *testing.T) {
	tests := []struct {
		status   HealthStatus
		expected string
	}{
		{StatusHealthy, "healthy"},
		{StatusUnhealthy, "unhealthy"},
		{StatusError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.status.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.status.String())
			}
		})
	}
}
