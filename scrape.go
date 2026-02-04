package gtmlp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

// Scrape extracts data from HTML using XPath with a typed result.
// It finds all container nodes and extracts fields from each one.
// Returns an empty slice if no containers are found.
func Scrape[T any](ctx context.Context, html string, config *Config) ([]T, error) {
	getLogger().Debug("scraping started",
		"container", config.Container,
		"fields", len(config.Fields),
		"html_size", len(html))

	// Validate config
	if err := config.Validate(); err != nil {
		getLogger().Error("config validation failed",
			"error", err.Error())
		return nil, err
	}

	// Parse HTML
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		getLogger().Error("html parsing failed",
			"error", err.Error())
		return nil, &ScrapeError{
			Type:    ErrTypeParsing,
			Message: "failed to parse HTML",
			Cause:   err,
		}
	}

	// Find container nodes with fallback
	containerNodes, err := findContainers(doc, config.Container, config.AltContainer)
	if err != nil {
		return nil, err
	}

	var results []T

	for containerNodes.MoveNext() {
		containerNode := containerNodes.Current().(*htmlquery.NodeNavigator).Current()

		// Extract fields from this container
		fieldData := make(map[string]any)
		for fieldName, fieldConfig := range config.Fields {
			value, err := extractFieldWithPipes(ctx, containerNode, fieldConfig)
			if err != nil {
				return nil, err
			}
			fieldData[fieldName] = value
		}

		// Convert map to struct
		var result T
		if err := mapToStruct(fieldData, &result); err != nil {
			return nil, &ScrapeError{
				Type:    ErrTypeParsing,
				Message: "failed to convert map to struct",
				Cause:   err,
			}
		}

		results = append(results, result)
	}

	// Return empty slice if no containers found
	if results == nil {
		results = []T{}
	}

	getLogger().Info("scraping completed",
		"items", len(results),
		"container", config.Container)

	return results, nil
}

// ScrapeUntyped extracts data from HTML using XPath, returning map slices.
// It finds all container nodes and extracts fields from each one.
// Returns an empty slice if no containers are found.
func ScrapeUntyped(ctx context.Context, html string, config *Config) ([]map[string]any, error) {
	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Parse HTML
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return nil, &ScrapeError{
			Type:    ErrTypeParsing,
			Message: "failed to parse HTML",
			Cause:   err,
		}
	}

	// Find container nodes with fallback
	containerNodes, err := findContainers(doc, config.Container, config.AltContainer)
	if err != nil {
		return nil, err
	}

	var results []map[string]any

	for containerNodes.MoveNext() {
		containerNode := containerNodes.Current().(*htmlquery.NodeNavigator).Current()

		// Extract fields from this container
		fieldData := make(map[string]any)
		for fieldName, fieldConfig := range config.Fields {
			value, err := extractFieldWithPipes(ctx, containerNode, fieldConfig)
			if err != nil {
				return nil, err
			}
			fieldData[fieldName] = value
		}

		results = append(results, fieldData)
	}

	// Return empty slice if no containers found
	if results == nil {
		results = []map[string]any{}
	}

	return results, nil
}

// findContainers finds container nodes with altContainer fallback support
func findContainers(doc *html.Node, container string, altContainers []string) (*xpath.NodeIterator, error) {
	// Build list of container XPaths to try
	containers := []string{container}
	containers = append(containers, altContainers...)

	getLogger().Debug("finding containers",
		"primary", container,
		"alternatives", len(altContainers))

	// Try each container XPath in sequence
	for i, containerXPath := range containers {
		// Compile container XPath
		containerExpr, err := xpath.Compile(containerXPath)
		if err != nil {
			getLogger().Error("container xpath compilation failed",
				"xpath", containerXPath,
				"error", err.Error())
			return nil, &ScrapeError{
				Type:    ErrTypeXPath,
				Message: "invalid container XPath",
				XPath:   containerXPath,
				Cause:   err,
			}
		}

		// Find container nodes
		containerNodes := containerExpr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)

		// Check if we found any containers
		if containerNodes.MoveNext() {
			// Found at least one container, re-evaluate to get a fresh iterator
			if i > 0 {
				getLogger().Warn("container fallback used",
					"primary", container,
					"used", containerXPath,
					"fallback_index", i)
			}
			getLogger().Debug("containers found",
				"xpath", containerXPath)
			return containerExpr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator), nil
		}

		getLogger().Debug("container xpath returned empty",
			"xpath", containerXPath)
		// No containers found, try next XPath
	}

	// All container XPaths failed, return empty iterator
	// Create an empty iterator by using an XPath that matches nothing
	emptyExpr, _ := xpath.Compile("/*[false()]")
	return emptyExpr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator), nil
}

// extractField extracts a value from a node using XPath
func extractField(containerNode *html.Node, fieldXPath string) any {
	// Compile field XPath
	expr, err := xpath.Compile(fieldXPath)
	if err != nil {
		return ""
	}

	// Evaluate XPath relative to container node
	nodeIterator := expr.Evaluate(htmlquery.CreateXPathNavigator(containerNode)).(*xpath.NodeIterator)

	// Move to first result
	if !nodeIterator.MoveNext() {
		return ""
	}

	navigator := nodeIterator.Current().(*htmlquery.NodeNavigator)

	// Check if it's an attribute node
	if navigator.NodeType() == xpath.AttributeNode {
		return navigator.Value()
	}

	// For other nodes, get the HTML node
	node := navigator.Current()

	// Extract value based on node type
	return extractValue(node)
}

// extractValue extracts text or attribute value from a node
func extractValue(node *html.Node) any {
	if node == nil {
		return ""
	}

	// Return the node data (text content or attribute value)
	// For attribute nodes selected via @attr, Data contains the value
	// For element nodes, we extract text content
	if node.Type == html.TextNode || node.Type == html.CommentNode {
		return strings.TrimSpace(node.Data)
	}

	// For element nodes, get text content
	return strings.TrimSpace(htmlquery.InnerText(node))
}

// extractFieldWithPipes extracts a value and applies pipes, with altXpath fallback
func extractFieldWithPipes(ctx context.Context, containerNode *html.Node, fieldConfig FieldConfig) (any, error) {
	// Build list of XPaths to try: primary + alternatives
	xpaths := []string{fieldConfig.XPath}
	xpaths = append(xpaths, fieldConfig.AltXPath...)

	// Try each XPath in sequence
	for xpathIdx, xpath := range xpaths {
		// Extract raw value with XPath
		rawValue := extractField(containerNode, xpath)

		// Convert to string for pipe processing
		inputStr, ok := rawValue.(string)
		if !ok {
			inputStr = fmt.Sprintf("%v", rawValue)
		}

		// Apply pipes if defined
		value := any(inputStr)
		if len(fieldConfig.Pipes) > 0 {
			for _, pipeDef := range fieldConfig.Pipes {
				pipeName, params := parsePipeDefinition(pipeDef)
				pipe := getPipe(pipeName)

				if pipe == nil {
					return "", &ScrapeError{
						Type:    ErrTypePipe,
						Message: fmt.Sprintf("unknown pipe '%s'", pipeName),
					}
				}

				result, err := pipe(ctx, inputStr, params)
				if err != nil {
					return "", &ScrapeError{
						Type:    ErrTypePipe,
						Message: fmt.Sprintf("pipe '%s' failed", pipeName),
						Cause:   &PipeError{PipeName: pipeName, Input: inputStr, Params: params, Cause: err},
					}
				}

				value = result
				// Convert result to string for next pipe
				inputStr = fmt.Sprintf("%v", result)
			}
		}

		// Check if result is non-empty after pipes
		if !isEmpty(value) {
			if len(xpaths) > 1 && xpathIdx > 0 {
				getLogger().Warn("field fallback used",
					"primary", fieldConfig.XPath,
					"used", xpath,
					"fallback_index", xpathIdx)
			}
			return value, nil
		}

		getLogger().Debug("field xpath returned empty after pipes",
			"xpath", xpath)
		// Result is empty, try next XPath
	}

	// All XPaths failed, return empty string
	getLogger().Warn("all xpaths failed for field",
		"primary", fieldConfig.XPath,
		"alternatives", len(fieldConfig.AltXPath))
	return "", nil
}

// isEmpty checks if a value is considered empty after pipe processing
func isEmpty(value any) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case int, int64:
		return v == 0
	case float64:
		return v == 0.0
	default:
		// For other types, convert to string and check
		str := fmt.Sprintf("%v", v)
		return str == "" || str == "0"
	}
}

// mapToStruct converts a map to a struct using JSON tags
func mapToStruct(m map[string]any, target any) error {
	// Convert map to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Convert JSON to struct
	return json.Unmarshal(jsonData, target)
}

// ScrapeURL fetches a URL and scrapes it with config (typed)
func ScrapeURL[T any](ctx context.Context, url string, config *Config) ([]T, error) {
	// Check if pagination is configured
	if config.Pagination != nil {
		// Use pagination logic
		results, err := scrapeWithPagination[T](ctx, url, config, false)
		if err != nil {
			return nil, err
		}
		// Return combined items from all pages
		var allItems []T
		for _, page := range results.Pages {
			allItems = append(allItems, page.Items...)
		}
		return allItems, nil
	}

	// No pagination, single page scraping (backward compatible)
	html, err := fetchHTML(url, config)
	if err != nil {
		return nil, err
	}
	// Add URL to context for parseUrl pipe
	ctx = WithURL(ctx, url)
	return Scrape[T](ctx, html, config)
}

// ScrapeURLUntyped fetches a URL and scrapes it, returning maps (no type parameter)
func ScrapeURLUntyped(ctx context.Context, url string, config *Config) ([]map[string]any, error) {
	// Check if pagination is configured
	if config.Pagination != nil {
		// Use pagination logic
		results, err := scrapeWithPagination[map[string]any](ctx, url, config, false)
		if err != nil {
			return nil, err
		}
		// Return combined items from all pages
		var allItems []map[string]any
		for _, page := range results.Pages {
			allItems = append(allItems, page.Items...)
		}
		return allItems, nil
	}

	// No pagination, single page scraping (backward compatible)
	html, err := fetchHTML(url, config)
	if err != nil {
		return nil, err
	}
	// Add URL to context for parseUrl pipe
	ctx = WithURL(ctx, url)
	return ScrapeUntyped(ctx, html, config)
}
