package health

import (
	"testing"
	"time"
)

func TestCheckURLHealth(t *testing.T) {
	// Use real URLs that are likely to be accessible
	urls := []string{
		"https://example.com",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/404",
		"https://this-domain-definitely-does-not-exist-12345.com",
	}

	results := CheckURLHealth(urls, 5*time.Second)

	if len(results) != 4 {
		t.Fatalf("Expected 4 results, got: %d", len(results))
	}

	// Check example.com (should be alive)
	if !results[0].Alive {
		t.Errorf("example.com should be alive, got error: %s", results[0].Error)
	}

	// Check 404 URL (should not be considered alive since it's not 2xx-3xx)
	if results[2].Alive {
		t.Errorf("404 status should not be considered alive")
	}
	if results[2].StatusCode != 404 {
		t.Errorf("Expected status code 404, got: %d", results[2].StatusCode)
	}

	// Check non-existent domain (should fail)
	if results[3].Alive {
		t.Errorf("Non-existent domain should not be alive")
	}
	if results[3].Error == "" {
		t.Errorf("Non-existent domain should have an error")
	}
}

func TestCheckURLHealthSequential(t *testing.T) {
	urls := []string{
		"https://example.com",
	}

	results := CheckURLHealthSequential(urls, 5*time.Second)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got: %d", len(results))
	}

	if !results[0].Alive {
		t.Errorf("example.com should be alive, got error: %s", results[0].Error)
	}
}

func TestCheckURLHealthWithGet(t *testing.T) {
	urls := []string{
		"https://example.com",
	}

	results := CheckURLHealthWithGet(urls, 5*time.Second)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got: %d", len(results))
	}

	if !results[0].Alive {
		t.Errorf("example.com should be alive with GET, got error: %s", results[0].Error)
	}
}

func TestCheckURLHealthEmpty(t *testing.T) {
	urls := []string{}

	results := CheckURLHealth(urls, 5*time.Second)

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty URL list, got: %d", len(results))
	}
}
