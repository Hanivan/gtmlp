package gtmlp

import (
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

// fetch fetches a URL and returns the HTTP response
func fetch(url string, config *Config) (*http.Response, error) {
	// Validate URL
	if url == "" {
		return nil, &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "URL cannot be empty",
		}
	}

	// Parse and validate URL
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return nil, &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "invalid URL format",
			URL:     url,
			Cause:   err,
		}
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "URL scheme must be http or https",
			URL:     url,
		}
	}

	// Create HTTP client with configured timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// Configure proxy if specified
	if config.Proxy != "" {
		proxyURL, err := neturl.Parse(config.Proxy)
		if err != nil {
			return nil, &ScrapeError{
				Type:    ErrTypeNetwork,
				Message: "invalid proxy URL",
				URL:     url,
				Cause:   err,
			}
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	// Perform request with retry logic
	var lastErr error
	maxAttempts := config.MaxRetries + 1

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Add exponential backoff delay between retries
		if attempt > 0 {
			backoffDuration := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(backoffDuration)
		}

		// Build request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = &ScrapeError{
				Type:    ErrTypeNetwork,
				Message: "failed to create HTTP request",
				URL:     url,
				Cause:   err,
			}
			continue
		}

		// Set User-Agent
		if config.UserAgent != "" {
			req.Header.Set("User-Agent", config.UserAgent)
		}

		// Set default headers
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")

		// Set custom headers
		for key, value := range config.Headers {
			req.Header.Set(key, value)
		}

		// Execute request
		resp, err := client.Do(req)
		if err != nil {
			lastErr = &ScrapeError{
				Type:    ErrTypeNetwork,
				Message: "HTTP request failed",
				URL:     url,
				Cause:   err,
			}
			continue
		}

		// Check status code
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			lastErr = &ScrapeError{
				Type:    ErrTypeNetwork,
				Message: fmt.Sprintf("HTTP request failed with status: %d %s", resp.StatusCode, resp.Status),
				URL:     url,
			}
			continue
		}

		// Success
		return resp, nil
	}

	// All retries exhausted
	return nil, lastErr
}

// fetchHTML fetches a URL and returns the HTML content as a string
func fetchHTML(url string, config *Config) (string, error) {
	resp, err := fetch(url, config)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "failed to read response body",
			URL:     url,
			Cause:   err,
		}
	}

	// Convert to string and trim whitespace
	html := strings.TrimSpace(string(body))

	return html, nil
}
