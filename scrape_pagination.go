package gtmlp

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

// Default pagination configuration
const (
	DefaultMaxPages      = 100
	DefaultPaginationTimeout = 10 * time.Minute
)

// ScrapeURLWithPages fetches a URL and scrapes it with pagination, returning page-separated results
func ScrapeURLWithPages[T any](ctx context.Context, url string, config *Config) (*PaginatedResults[T], error) {
	if config.Pagination == nil {
		// No pagination config, scrape single page
		items, err := ScrapeURL[T](ctx, url, config)
		if err != nil {
			return nil, err
		}
		return &PaginatedResults[T]{
			Pages: []PageResult[T]{{
				URL:       url,
				PageNum:   1,
				Items:     items,
				ScrapedAt: time.Now(),
			}},
			TotalPages: 1,
			TotalItems: len(items),
		}, nil
	}

	return scrapeWithPagination[T](ctx, url, config, true)
}

// ExtractPaginationURLs extracts all pagination URLs without scraping
func ExtractPaginationURLs(ctx context.Context, url string, config *Config) (*PaginationInfo, error) {
	if config.Pagination == nil {
		return nil, &ScrapeError{
			Type:    ErrTypeConfig,
			Message: "pagination config is required",
		}
	}

	// Fetch first page
	htmlContent, err := fetchHTML(url, config)
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, &ScrapeError{
			Type:    ErrTypeParsing,
			Message: "failed to parse HTML",
			Cause:   err,
		}
	}

	var urls []string
	switch config.Pagination.Type {
	case "next-link":
		urls, err = extractNextLinkChain(ctx, url, doc, config)
	case "numbered":
		urls, err = extractNumberedPages(ctx, url, doc, config)
	default:
		return nil, &ScrapeError{
			Type:    ErrTypeConfig,
			Message: fmt.Sprintf("unknown pagination type: %s", config.Pagination.Type),
		}
	}

	if err != nil {
		return nil, err
	}

	return &PaginationInfo{
		URLs:    urls,
		Type:    config.Pagination.Type,
		BaseURL: url,
	}, nil
}

// scrapeWithPagination handles pagination logic for auto-follow mode
func scrapeWithPagination[T any](ctx context.Context, startURL string, config *Config, separatePages bool) (*PaginatedResults[T], error) {
	applyPaginationDefaults(config.Pagination)

	var allPages []PageResult[T]
	var allItems []T
	visitedURLs := make(map[string]bool)
	pageNum := 1
	currentURL := startURL
	startTime := time.Now()

	getLogger().Info("pagination starting",
		"url", startURL,
		"type", config.Pagination.Type,
		"max_pages", config.Pagination.MaxPages)

	for {
		// Check timeout
		if time.Since(startTime) > config.Pagination.Timeout {
			getLogger().Warn("pagination timeout exceeded",
				"timeout", config.Pagination.Timeout,
				"elapsed", time.Since(startTime),
				"pages_scraped", pageNum-1)
			break
		}

		// Check max pages
		if pageNum > config.Pagination.MaxPages {
			getLogger().Warn("pagination max pages reached",
				"max_pages", config.Pagination.MaxPages,
				"total_items", len(allItems))
			break
		}

		// Mark URL as visited
		normalized := normalizeURL(currentURL)
		if visitedURLs[normalized] {
			getLogger().Warn("pagination duplicate url",
				"url", currentURL,
				"page", pageNum)
			break
		}
		visitedURLs[normalized] = true

		// Scrape current page
		ctx = WithURL(ctx, currentURL)
		items, err := scrapeCurrentPage[T](ctx, currentURL, config)
		if err != nil {
			// Return error with partial data
			return nil, &PaginationError{
				PageURL:      currentURL,
				PageNumber:   pageNum,
				PartialData:  allItems,
				TotalScraped: len(allItems),
				Cause:        err,
			}
		}

		// Log progress
		getLogger().Info("pagination page scraped",
			"page", pageNum,
			"items", len(items),
			"total_items", len(allItems)+len(items),
			"url", currentURL)

		// Store results
		pageResult := PageResult[T]{
			URL:       currentURL,
			PageNum:   pageNum,
			Items:     items,
			ScrapedAt: time.Now(),
		}
		allPages = append(allPages, pageResult)
		allItems = append(allItems, items...)

		// Get next URL
		nextURL, err := getNextPageURL(ctx, currentURL, config)
		if err != nil {
			return nil, &PaginationError{
				PageURL:      currentURL,
				PageNumber:   pageNum,
				PartialData:  allItems,
				TotalScraped: len(allItems),
				Cause:        err,
			}
		}

		if nextURL == "" {
			// No more pages
			getLogger().Info("pagination complete",
				"reason", "no_next_link",
				"pages", len(allPages))
			break
		}

		getLogger().Info("pagination following next link",
			"url", nextURL,
			"page", pageNum+1)
		currentURL = nextURL
		pageNum++
	}

	getLogger().Info("pagination complete",
		"pages", len(allPages),
		"total_items", len(allItems),
		"duration", time.Since(startTime).String())

	return &PaginatedResults[T]{
		Pages:      allPages,
		TotalPages: len(allPages),
		TotalItems: len(allItems),
	}, nil
}

// scrapeCurrentPage scrapes a single page without pagination logic
func scrapeCurrentPage[T any](ctx context.Context, url string, config *Config) ([]T, error) {
	htmlContent, err := fetchHTML(url, config)
	if err != nil {
		return nil, err
	}
	return Scrape[T](ctx, htmlContent, config)
}

// getNextPageURL extracts the next page URL based on pagination type
func getNextPageURL(ctx context.Context, currentURL string, config *Config) (string, error) {
	htmlContent, err := fetchHTML(currentURL, config)
	if err != nil {
		return "", err
	}

	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", &ScrapeError{
			Type:    ErrTypeParsing,
			Message: "failed to parse HTML",
			Cause:   err,
		}
	}

	switch config.Pagination.Type {
	case "next-link":
		return extractNextURL(ctx, currentURL, doc, config)
	case "numbered":
		// For numbered pagination in auto-follow, we don't use getNextPageURL
		// Instead, we extract all URLs upfront
		return "", nil
	default:
		return "", &ScrapeError{
			Type:    ErrTypeConfig,
			Message: fmt.Sprintf("unknown pagination type: %s", config.Pagination.Type),
		}
	}
}

// extractNextURL extracts the next page URL using NextSelector and AltSelectors
func extractNextURL(ctx context.Context, baseURL string, doc *html.Node, config *Config) (string, error) {
	selectors := []string{config.Pagination.NextSelector}
	selectors = append(selectors, config.Pagination.AltSelectors...)

	for _, selector := range selectors {
		if selector == "" {
			continue
		}

		// Compile XPath
		expr, err := xpath.Compile(selector)
		if err != nil {
			continue // Try next selector
		}

		// Evaluate XPath
		nodeIterator := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
		if !nodeIterator.MoveNext() {
			continue // Try next selector
		}

		navigator := nodeIterator.Current().(*htmlquery.NodeNavigator)
		rawURL := navigator.Value()

		if rawURL == "" {
			continue
		}

		// Apply pipes
		processedURL, err := applyPipesToURL(ctx, rawURL, config.Pagination.Pipes)
		if err != nil {
			continue
		}

		if processedURL == "" {
			continue
		}

		// Resolve relative URL
		absoluteURL, err := resolveURL(baseURL, processedURL)
		if err != nil {
			continue
		}

		return absoluteURL, nil
	}

	// No next link found
	return "", nil
}

// extractNumberedPages extracts all page URLs for numbered pagination
func extractNumberedPages(ctx context.Context, baseURL string, doc *html.Node, config *Config) ([]string, error) {
	if config.Pagination.PageSelector == "" {
		return nil, &ScrapeError{
			Type:    ErrTypeConfig,
			Message: "pageSelector is required for numbered pagination",
		}
	}

	// Compile XPath
	expr, err := xpath.Compile(config.Pagination.PageSelector)
	if err != nil {
		return nil, &ScrapeError{
			Type:    ErrTypeXPath,
			Message: "invalid pageSelector",
			XPath:   config.Pagination.PageSelector,
			Cause:   err,
		}
	}

	// Evaluate XPath
	nodeIterator := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)

	var urls []string
	seenURLs := make(map[string]bool)

	for nodeIterator.MoveNext() {
		navigator := nodeIterator.Current().(*htmlquery.NodeNavigator)
		rawURL := navigator.Value()

		if rawURL == "" {
			continue
		}

		// Apply pipes
		processedURL, err := applyPipesToURL(ctx, rawURL, config.Pagination.Pipes)
		if err != nil {
			continue
		}

		if processedURL == "" {
			continue
		}

		// Resolve relative URL
		absoluteURL, err := resolveURL(baseURL, processedURL)
		if err != nil {
			continue
		}

		// Deduplicate
		normalized := normalizeURL(absoluteURL)
		if seenURLs[normalized] {
			continue
		}
		seenURLs[normalized] = true

		urls = append(urls, absoluteURL)
	}

	return urls, nil
}

// extractNextLinkChain follows next links to build a list of all page URLs
func extractNextLinkChain(ctx context.Context, startURL string, doc *html.Node, config *Config) ([]string, error) {
	var urls []string
	visitedURLs := make(map[string]bool)
	currentURL := startURL
	currentDoc := doc
	pageCount := 0

	applyPaginationDefaults(config.Pagination)

	for {
		if pageCount >= config.Pagination.MaxPages {
			break
		}

		// Mark as visited
		normalized := normalizeURL(currentURL)
		if visitedURLs[normalized] {
			break // Circular reference
		}
		visitedURLs[normalized] = true
		urls = append(urls, currentURL)
		pageCount++

		// Get next URL
		nextURL, err := extractNextURL(ctx, currentURL, currentDoc, config)
		if err != nil || nextURL == "" {
			break
		}

		// Fetch next page
		htmlContent, err := fetchHTML(nextURL, config)
		if err != nil {
			break
		}

		currentDoc, err = htmlquery.Parse(strings.NewReader(htmlContent))
		if err != nil {
			break
		}

		currentURL = nextURL
	}

	return urls, nil
}

// applyPipesToURL applies pipes to a URL string
func applyPipesToURL(ctx context.Context, rawURL string, pipes []string) (string, error) {
	if len(pipes) == 0 {
		return rawURL, nil
	}

	result := rawURL
	for _, pipeDef := range pipes {
		pipeName, params := parsePipeDefinition(pipeDef)
		pipe := getPipe(pipeName)

		if pipe == nil {
			return "", &ScrapeError{
				Type:    ErrTypePipe,
				Message: fmt.Sprintf("unknown pipe '%s'", pipeName),
			}
		}

		processed, err := pipe(ctx, result, params)
		if err != nil {
			return "", err
		}

		result = fmt.Sprintf("%v", processed)
	}

	return result, nil
}

// resolveURL converts a relative URL to absolute using the base URL
func resolveURL(baseURL, relativeURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	rel, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}

	resolved := base.ResolveReference(rel)
	return resolved.String(), nil
}

// normalizeURL normalizes a URL for duplicate detection
func normalizeURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// Strip fragment
	u.Fragment = ""

	// Normalize path (remove trailing slash if present)
	if u.Path != "/" && strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	// Sort query parameters
	if u.RawQuery != "" {
		query := u.Query()
		u.RawQuery = sortQueryParams(query)
	}

	return u.String()
}

// sortQueryParams sorts query parameters for consistent comparison
func sortQueryParams(query url.Values) string {
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		values := query[k]
		sort.Strings(values)
		for _, v := range values {
			parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
		}
	}

	return strings.Join(parts, "&")
}

// applyPaginationDefaults applies default values to pagination config
func applyPaginationDefaults(config *PaginationConfig) {
	if config.MaxPages == 0 {
		config.MaxPages = DefaultMaxPages
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultPaginationTimeout
	}
}

// logPagination is deprecated and removed
// Use gtmlp.SetLogLevel(slog.LevelInfo) to see pagination logs
