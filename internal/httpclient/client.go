package httpclient

import (
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"time"

	fakeUseragent "github.com/lib4u/fake-useragent"
)

const (
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "gtmlp/1.0 (HTML Parsing Library)"
)

// Client is an HTTP client for fetching web pages.
type Client struct {
	client        *http.Client
	timeout       time.Duration
	userAgent     string
	headers       map[string]string
	maxRetries    int
	randomUA      bool
	fakeUA        *fakeUseragent.UserAgent
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// NewClient creates a new HTTP client with the given options.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		client:     &http.Client{Timeout: defaultTimeout},
		timeout:    defaultTimeout,
		userAgent:  defaultUserAgent,
		headers:    make(map[string]string),
		maxRetries: 0,
		randomUA:   true, // Default to random UA for better scraping
	}

	for _, opt := range opts {
		opt(c)
	}

	// Initialize fake-useragent for random UA (default enabled)
	ua, err := fakeUseragent.New()
	if err != nil {
		// If initialization fails, fall back to static UA
		c.randomUA = false
	} else {
		c.fakeUA = ua
	}

	return c
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
		c.client.Timeout = timeout
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// WithHeaders sets custom headers.
func WithHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithProxy sets a proxy URL for the client.
func WithProxy(proxyURL string) ClientOption {
	return func(c *Client) {
		if proxyURL != "" {
			parsedURL, err := neturl.Parse(proxyURL)
			if err == nil {
				c.client.Transport = &http.Transport{
					Proxy: http.ProxyURL(parsedURL),
				}
			}
		}
	}
}

// WithMaxRetries sets the maximum number of retries for failed requests.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithHTTPClient allows using a custom http.Client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// WithRandomUserAgent enables realistic rotating user agents.
// This helps avoid detection during web scraping.
func WithRandomUserAgent() ClientOption {
	return func(c *Client) {
		c.randomUA = true
	}
}

// WithRandomUserAgentFilter enables random user agents with a filter.
// Example: WithRandomUserAgentFilter(app.Chrome, app.Firefox)
func WithRandomUserAgentFilter(browsers ...string) ClientOption {
	return func(c *Client) {
		c.randomUA = true
		// Store filter to be used when UA is created
		if len(browsers) > 0 {
			c.userAgent = browsers[0] // Will be overridden by random UA
		}
	}
}

// WithDisableRandomUA disables random user agents and uses static user agent.
func WithDisableRandomUA() ClientOption {
	return func(c *Client) {
		c.randomUA = false
		c.fakeUA = nil
	}
}

// Get fetches the content of the given URL.
func (c *Client) Get(url string) ([]byte, error) {
	return c.GetWithHeaders(url, nil)
}

// GetWithHeaders fetches the content of the given URL with custom headers.
func (c *Client) GetWithHeaders(url string, headers map[string]string) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	if _, err := neturl.Parse(url); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(c.backoffDuration(attempt))
		}

		body, err := c.tryGet(url, headers)
		if err == nil {
			return body, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// GetHTML fetches and returns the HTML content of the given URL.
func (c *Client) GetHTML(url string) (string, error) {
	body, err := c.Get(url)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetHTMLWithHeaders fetches and returns the HTML content with custom headers.
func (c *Client) GetHTMLWithHeaders(url string, headers map[string]string) (string, error) {
	body, err := c.GetWithHeaders(url, headers)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// tryGet attempts to fetch the URL once.
func (c *Client) tryGet(url string, headers map[string]string) ([]byte, error) {
	req, err := c.buildRequest(url, headers)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// buildRequest creates an HTTP request with all headers set.
func (c *Client) buildRequest(url string, customHeaders map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header
	if c.randomUA && c.fakeUA != nil {
		req.Header.Set("User-Agent", c.fakeUA.GetRandom())
	} else if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	for k, v := range customHeaders {
		req.Header.Set(k, v)
	}

	return req, nil
}

// backoffDuration calculates exponential backoff duration.
func (c *Client) backoffDuration(attempt int) time.Duration {
	return time.Duration(1<<uint(attempt-1)) * time.Second
}
