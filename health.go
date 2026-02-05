package gtmlp

import (
	"io"
	"net/http"
	neturl "net/url"
	"sync"
	"time"
)

// HealthStatus represents the health status of a URL
type HealthStatus int

const (
	// StatusHealthy indicates the URL returned a 2xx status code
	StatusHealthy HealthStatus = iota
	// StatusUnhealthy indicates the URL returned a 4xx or 5xx status code
	StatusUnhealthy
	// StatusError indicates there was a network or other error
	StatusError
)

// String returns the string representation of HealthStatus
func (s HealthStatus) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusUnhealthy:
		return "unhealthy"
	case StatusError:
		return "error"
	default:
		return "unknown"
	}
}

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	URL     string        // The URL that was checked
	Status  HealthStatus  // The health status of the URL
	Code    int           // HTTP status code (0 if error occurred)
	Latency time.Duration // Time taken for the health check
	Error   error         // Error message if check failed
}

// CheckHealth performs a health check on a single URL
func CheckHealth(url string) HealthCheckResult {
	// Use default config
	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
	}
	return CheckHealthWithOptions(url, config)
}

// CheckHealthMulti performs health checks on multiple URLs concurrently
func CheckHealthMulti(urls []string) []HealthCheckResult {
	if len(urls) == 0 {
		return []HealthCheckResult{}
	}

	getLogger().Info("health check multi starting",
		"urls_count", len(urls))

	results := make([]HealthCheckResult, len(urls))
	var wg sync.WaitGroup

	// Use default config for all checks
	config := &Config{
		Timeout:   10 * time.Second,
		UserAgent: "GTMLP/2.0",
	}

	for i, url := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()
			results[idx] = CheckHealthWithOptions(u, config)
		}(i, url)
	}

	wg.Wait()

	// Log summary
	healthy := 0
	unhealthy := 0
	errors := 0
	for _, r := range results {
		switch r.Status {
		case StatusHealthy:
			healthy++
		case StatusUnhealthy:
			unhealthy++
		case StatusError:
			errors++
		}
	}

	getLogger().Info("health check multi completed",
		"total", len(urls),
		"healthy", healthy,
		"unhealthy", unhealthy,
		"errors", errors)

	return results
}

// fetchForHealth fetches a URL and returns the HTTP response, even for 4xx/5xx status codes
// Unlike the regular fetch function, this doesn't treat non-2xx codes as errors
func fetchForHealth(url string, config *Config) (*http.Response, error) {
	getLogger().Debug("health check fetch starting",
		"url", url,
		"timeout", config.Timeout)

	// Validate URL
	if url == "" {
		getLogger().Error("health check empty url")
		return nil, &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "URL cannot be empty",
		}
	}

	// Parse and validate URL
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		getLogger().Error("health check invalid url format",
			"url", url,
			"error", err.Error())
		return nil, &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "invalid URL format",
			URL:     url,
			Cause:   err,
		}
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		getLogger().Error("health check invalid url scheme",
			"url", url,
			"scheme", parsedURL.Scheme)
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
			getLogger().Error("health check invalid proxy url",
				"url", url,
				"proxy", config.Proxy,
				"error", err.Error())
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
		getLogger().Debug("health check using proxy",
			"url", url,
			"proxy", config.Proxy)
	}

	// Build request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		getLogger().Error("health check failed to create request",
			"url", url,
			"error", err.Error())
		return nil, &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "failed to create HTTP request",
			URL:     url,
			Cause:   err,
		}
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

	getLogger().Debug("health check sending request",
		"url", url,
		"user_agent", config.UserAgent)

	// Execute request (no retries for health checks)
	resp, err := client.Do(req)
	if err != nil {
		getLogger().Warn("health check request failed",
			"url", url,
			"error", err.Error())
		return nil, &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "HTTP request failed",
			URL:     url,
			Cause:   err,
		}
	}

	getLogger().Debug("health check response received",
		"url", url,
		"status", resp.StatusCode)

	return resp, nil
}

// CheckHealthWithOptions performs a health check on a single URL with custom configuration
func CheckHealthWithOptions(url string, config *Config) HealthCheckResult {
	result := HealthCheckResult{
		URL:    url,
		Status: StatusError,
		Code:   0,
	}

	getLogger().Info("health check starting",
		"url", url,
		"timeout", config.Timeout)

	// Validate URL
	if url == "" {
		getLogger().Error("health check empty url")
		result.Error = &ScrapeError{
			Type:    ErrTypeNetwork,
			Message: "URL cannot be empty",
		}
		return result
	}

	// Start latency measurement
	startTime := time.Now()

	// Use custom fetch function that doesn't treat 4xx/5xx as errors
	resp, err := fetchForHealth(url, config)

	// Measure latency
	result.Latency = time.Since(startTime)

	if err != nil {
		result.Error = err
		getLogger().Info("health check completed with error",
			"url", url,
			"status", "error",
			"latency_ms", result.Latency.Milliseconds(),
			"error", err.Error())
		return result
	}

	// Close response body
	defer resp.Body.Close()

	// Drain the response body to ensure the connection can be reused
	_, _ = io.Copy(io.Discard, resp.Body)

	// Set status code
	result.Code = resp.StatusCode

	// Determine health status based on status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = StatusHealthy
		getLogger().Info("health check completed",
			"url", url,
			"status", "healthy",
			"status_code", resp.StatusCode,
			"latency_ms", result.Latency.Milliseconds())
	} else if resp.StatusCode >= 400 {
		result.Status = StatusUnhealthy
		getLogger().Info("health check completed",
			"url", url,
			"status", "unhealthy",
			"status_code", resp.StatusCode,
			"latency_ms", result.Latency.Milliseconds())
	} else {
		// 3xx redirects - consider unhealthy for health check purposes
		result.Status = StatusUnhealthy
		getLogger().Warn("health check encountered redirect",
			"url", url,
			"status", "unhealthy",
			"status_code", resp.StatusCode,
			"latency_ms", result.Latency.Milliseconds())
	}

	return result
}
