package health

import (
	"net/http"
	"sync"
	"time"
)

// HealthResult represents the health check result for a URL.
type HealthResult struct {
	URL        string // The URL that was checked
	Alive      bool   // Whether the URL is accessible
	StatusCode int    // HTTP status code
	Error      string // Error message if check failed
}

// checkSingleURL performs a health check on a single URL using the specified method.
func checkSingleURL(client *http.Client, url string, useGet bool) HealthResult {
	result := HealthResult{URL: url}

	var resp *http.Response
	var err error

	if useGet {
		resp, err = client.Get(url)
	} else {
		resp, err = client.Head(url)
	}

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Alive = resp.StatusCode >= 200 && resp.StatusCode < 400
	return result
}

// checkURLsConcurrent checks URLs concurrently using the specified HTTP method.
func checkURLsConcurrent(urls []string, timeout time.Duration, useGet bool) []HealthResult {
	results := make([]HealthResult, len(urls))
	client := &http.Client{Timeout: timeout}

	var wg sync.WaitGroup
	for i, u := range urls {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()
			results[idx] = checkSingleURL(client, url, useGet)
		}(i, u)
	}

	wg.Wait()
	return results
}

// CheckURLHealth checks the health of multiple URLs concurrently.
func CheckURLHealth(urls []string, timeout time.Duration) []HealthResult {
	return checkURLsConcurrent(urls, timeout, false)
}

// CheckURLHealthSequential checks URLs sequentially (slower but more controlled).
func CheckURLHealthSequential(urls []string, timeout time.Duration) []HealthResult {
	results := make([]HealthResult, 0, len(urls))
	client := &http.Client{Timeout: timeout}

	for _, u := range urls {
		results = append(results, checkSingleURL(client, u, false))
	}

	return results
}

// CheckURLHealthWithGet checks URL health using GET request instead of HEAD.
// Use this when servers don't support HEAD requests.
func CheckURLHealthWithGet(urls []string, timeout time.Duration) []HealthResult {
	return checkURLsConcurrent(urls, timeout, true)
}
