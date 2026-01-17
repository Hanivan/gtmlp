package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.timeout != 30*time.Second {
		t.Errorf("Default timeout mismatch, got: %v", client.timeout)
	}

	if client.userAgent == "" {
		t.Error("Default user agent should not be empty")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	opts := []ClientOption{
		WithTimeout(10 * time.Second),
		WithUserAgent("TestBot/1.0"),
		WithHeaders(map[string]string{
			"X-Custom": "value",
		}),
	}

	client := NewClient(opts...)

	if client.timeout != 10*time.Second {
		t.Errorf("Timeout option not applied, got: %v", client.timeout)
	}

	if client.userAgent != "TestBot/1.0" {
		t.Errorf("UserAgent option not applied, got: %s", client.userAgent)
	}

	if client.headers["X-Custom"] != "value" {
		t.Error("Headers option not applied")
	}
}

func TestClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test</body></html>"))
	}))
	defer server.Close()

	client := NewClient()
	body, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if string(body) != "<html><body>Test</body></html>" {
		t.Errorf("Get() body mismatch: %s", string(body))
	}
}

func TestClientGetWithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom") != "test-value" {
			t.Errorf("Custom header not received: %s", r.Header.Get("X-Custom"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := NewClient()
	body, err := client.GetWithHeaders(server.URL, map[string]string{
		"X-Custom": "test-value",
	})
	if err != nil {
		t.Fatalf("GetWithHeaders() failed: %v", err)
	}

	if string(body) != "OK" {
		t.Errorf("GetWithHeaders() body mismatch: %s", string(body))
	}
}

func TestClientGetUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if ua != "CustomBot/1.0" {
			t.Errorf("User-Agent mismatch: %s", ua)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(WithUserAgent("CustomBot/1.0"), WithDisableRandomUA())
	_, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
}

func TestClientGetHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test HTML</body></html>"))
	}))
	defer server.Close()

	client := NewClient()
	html, err := client.GetHTML(server.URL)
	if err != nil {
		t.Fatalf("GetHTML() failed: %v", err)
	}

	if html != "<html><body>Test HTML</body></html>" {
		t.Errorf("GetHTML() mismatch: %s", html)
	}
}

func TestClientGetError(t *testing.T) {
	client := NewClient(WithTimeout(1 * time.Second))

	// Test with invalid URL
	_, err := client.Get("invalid-url")
	if err == nil {
		t.Error("Get() with invalid URL should return error")
	}

	// Test with non-existent server
	_, err = client.Get("http://localhost:9999")
	if err == nil {
		t.Error("Get() to non-existent server should return error")
	}
}

func TestClientGetHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient()
	_, err := client.Get(server.URL)
	if err == nil {
		t.Error("Get() with 404 status should return error")
	}
}

func TestClientGetEmptyURL(t *testing.T) {
	client := NewClient()
	_, err := client.Get("")
	if err == nil {
		t.Error("Get() with empty URL should return error")
	}
}

func TestWithTimeout(t *testing.T) {
	client := NewClient(WithTimeout(5 * time.Second))
	if client.timeout != 5*time.Second {
		t.Errorf("WithTimeout() failed, got: %v", client.timeout)
	}
}

func TestWithUserAgent(t *testing.T) {
	client := NewClient(WithUserAgent("MyBot/1.0"))
	if client.userAgent != "MyBot/1.0" {
		t.Errorf("WithUserAgent() failed, got: %s", client.userAgent)
	}
}

func TestWithHeaders(t *testing.T) {
	client := NewClient(WithHeaders(map[string]string{
		"X-Header-1": "value1",
		"X-Header-2": "value2",
	}))

	if client.headers["X-Header-1"] != "value1" {
		t.Error("WithHeaders() failed to set X-Header-1")
	}

	if client.headers["X-Header-2"] != "value2" {
		t.Error("WithHeaders() failed to set X-Header-2")
	}
}

func TestWithMaxRetries(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Success"))
		}
	}))
	defer server.Close()

	client := NewClient(WithMaxRetries(3))
	body, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() with retries failed: %v", err)
	}

	if string(body) != "Success" {
		t.Errorf("Get() body mismatch after retries: %s", string(body))
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}
