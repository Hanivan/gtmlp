# GTMLP v2.0 API Reference

Complete API reference for GTMLP v2.0 - the configuration-based, type-safe scraping API.

## Table of Contents

- [Overview](#overview)
- [Core Scraping Functions](#core-scraping-functions)
- [Config Loading](#config-loading)
- [Data Transformation Pipes](#data-transformation-pipes)
- [XPath Validation](#xpath-validation)
- [Health Check](#health-check)
- [Types](#types)
- [Error Handling](#error-handling)
- [Complete Examples](#complete-examples)

## Overview

GTMLP v2.0 provides a streamlined, configuration-based API for web scraping with the following benefits:

- **Type Safety**: Go generics ensure compile-time type checking
- **External Configuration**: Define selectors in JSON/YAML files
- **Environment Variables**: Override settings via environment
- **XPath Validation**: Test selectors before scraping
- **Health Checks**: Verify URL availability before scraping

## Core Scraping Functions

### Scrape

Extracts data from HTML string using typed results.

```go
func Scrape[T any](html string, config *Config) ([]T, error)
```

**Parameters:**
- `html`: HTML content as string
- `config`: Scraping configuration with XPath selectors

**Returns:**
- `[]T`: Slice of extracted data (typed)
- `error`: Error if scraping fails

**Example:**

```go
type Article struct {
    Title   string `json:"title"`
    Content string `json:"content"`
    Author  string `json:"author"`
}

config := &gtmlp.Config{
    Container: "//article[@class='blog-post']",
    Fields: map[string]gtmlp.FieldConfig{
        "title":   {XPath: ".//h2/text()"},
        "content": {XPath: ".//div[@class='content']/text()"},
        "author":  {XPath: ".//span[@class='author']/text()"},
    },
}

html := `<html>...</html>`

articles, err := gtmlp.Scrape[Article](context.Background(), html, config)
if err != nil {
    log.Fatal(err)
}

for _, article := range articles {
    fmt.Printf("%s by %s\n", article.Title, article.Author)
}
```

**Notes:**
- Container XPath defines repeating elements
- Field XPaths are relative to each container
- Returns empty slice if no containers found
- Type parameter `T` must be a struct with JSON tags

### ScrapeUntyped

Extracts data from HTML string returning maps (no type parameter).

```go
func ScrapeUntyped(html string, config *Config) ([]map[string]any, error)
```

**Parameters:**
- `html`: HTML content as string
- `config`: Scraping configuration

**Returns:**
- `[]map[string]any`: Slice of extracted data as maps
- `error`: Error if scraping fails

**Example:**

```go
config := &gtmlp.Config{
    Container: "//div[@class='product']",
    Fields: map[string]gtmlp.FieldConfig{
        "name":  {XPath: ".//h2/text()"},
        "price": {XPath: ".//span[@class='price']/text()"},
    },
}

results, err := gtmlp.ScrapeUntyped(context.Background(), html, config)
if err != nil {
    log.Fatal(err)
}

for _, result := range results {
    fmt.Printf("%s: %v\n", result["name"], result["price"])
}
```

**Use when:**
- You don't need type safety
- Working with dynamic schemas
- Quick prototyping

### ScrapeURL

Fetches a URL and scrapes it with typed results.

```go
func ScrapeURL[T any](url string, config *Config) ([]T, error)
```

**Parameters:**
- `url`: URL to fetch and scrape
- `config`: Scraping configuration (includes HTTP settings)

**Returns:**
- `[]T`: Slice of extracted data (typed)
- `error`: Error if fetching or scraping fails

**Example:**

```go
config, _ := gtmlp.LoadConfig("selectors.json", nil)

products, err := gtmlp.ScrapeURL[Product](
    context.Background(),
    "https://example.com/products",
    config,
)
if err != nil {
    log.Fatal(err)
}
```

**HTTP Options in Config:**
```go
config := &gtmlp.Config{
    // XPath settings
    Container: "//div[@class='product']",
    Fields:    map[string]gtmlp.FieldConfig{...},

    // HTTP settings
    Timeout:    30 * time.Second,
    UserAgent:  "MyBot/1.0",
    RandomUA:   true,
    MaxRetries: 3,
    Proxy:      "http://proxy.example.com:8080",
    Headers: map[string]string{
        "Accept-Language": "en-US",
    },
}
```

### ScrapeURLUntyped

Fetches a URL and scrapes it, returning maps.

```go
func ScrapeURLUntyped(url string, config *Config) ([]map[string]any, error)
```

**Parameters:**
- `url`: URL to fetch and scrape
- `config`: Scraping configuration

**Returns:**
- `[]map[string]any`: Slice of extracted data as maps
- `error`: Error if fetching or scraping fails

**Example:**

```go
config, _ := gtmlp.LoadConfig("selectors.json", nil)

results, err := gtmlp.ScrapeURLUntyped(
    context.Background(),
    "https://example.com/products",
    config,
)
if err != nil {
    log.Fatal(err)
}
```

## Config Loading

### LoadConfig

Loads configuration from file (JSON/YAML auto-detected).

```go
func LoadConfig(path string, envMapping *EnvMapping) (*Config, error)
```

**Parameters:**
- `path`: Path to config file (.json, .yaml, or .yml)
- `envMapping`: Environment variable mapping (nil for defaults)

**Returns:**
- `*Config`: Loaded configuration
- `error`: Error if loading fails

**Example:**

```go
// Use default environment variable names
config, err := gtmlp.LoadConfig("selectors.json", nil)
if err != nil {
    log.Fatal(err)
}

// Use custom environment variable names
customMapping := &gtmlp.EnvMapping{
    Timeout:    "MY_APP_TIMEOUT",
    UserAgent:  "MY_APP_UA",
    RandomUA:   "MY_APP_RANDOM_UA",
    MaxRetries: "MY_APP_RETRIES",
    Proxy:      "MY_APP_PROXY",
}

config, err := gtmlp.LoadConfig("selectors.yaml", customMapping)
```

**JSON Config Example:**
```json
{
  "container": "//div[@class='product']",
  "fields": {
    "name": {"xpath": ".//h2/text()", "pipes": ["trim"]},
    "price": {"xpath": ".//span[@class='price']/text()", "pipes": ["trim", "tofloat"]},
    "link": {"xpath": ".//a/@href", "pipes": ["parseurl"]}
  },
  "timeout": "30s",
  "user_agent": "GTMLP/2.0",
  "random_ua": false,
  "max_retries": 0,
  "proxy": "",
  "headers": {
    "Accept-Language": "en-US"
  }
}
```

**YAML Config Example:**
```yaml
container: "//div[@class='product']"
fields:
  name:
    xpath: ".//h2/text()"
    pipes: ["trim"]
  price:
    xpath: ".//span[@class='price']/text()"
    pipes: ["trim", "tofloat"]
  link:
    xpath: ".//a/@href"
    pipes: ["parseurl"]
timeout: 30s
user_agent: "GTMLP/2.0"
random_ua: false
max_retries: 0
proxy: ""
headers:
  Accept-Language: "en-US"
```

### ParseConfig

Parses configuration from string.

```go
func ParseConfig(data string, format ConfigFormat, envMapping *EnvMapping) (*Config, error)
```

**Parameters:**
- `data`: Configuration data as string
- `format`: Format type (FormatJSON or FormatYAML)
- `envMapping`: Environment variable mapping

**Returns:**
- `*Config`: Parsed configuration
- `error`: Error if parsing fails

**Example:**

```go
jsonConfig := `{
    "container": "//div[@class='product']",
    "fields": {
        "name": ".//h2/text()"
    }
}`

config, err := gtmlp.ParseConfig(
    jsonConfig,
    gtmlp.FormatJSON,
    nil, // use default env mapping
)
if err != nil {
    log.Fatal(err)
}
```

### Config.Validate

Validates the configuration.

```go
func (c *Config) Validate() error
```

**Returns:**
- `error`: Error if configuration is invalid

**Validation Rules:**
- Container XPath must not be empty
- At least one field must be defined
- Timeout must be positive

**Example:**

```go
config := &gtmlp.Config{
    Container: "//div[@class='product']",
    Fields: map[string]gtmlp.FieldConfig{
        "name": {XPath: ".//h2/text()"},
    },
    Timeout: 30 * time.Second,
}

if err := config.Validate(); err != nil {
    log.Fatalf("Invalid config: %v", err)
}
```

## Data Transformation Pipes

Pipes transform extracted field values after XPath extraction. Apply pipes in config using the `pipes` array:

```json
{
  "fields": {
    "name": {"xpath": ".//h2/text()", "pipes": ["trim"]},
    "price": {"xpath": ".//span[@class='price']/text()", "pipes": ["trim", "tofloat"]},
    "url": {"xpath": ".//a/@href", "pipes": ["parseurl"]}
  }
}
```

### Built-in Pipes

#### trim

Removes leading and trailing whitespace.

```json
{"name": {"xpath": ".//h2/text()", "pipes": ["trim"]}}
```

#### toint

Converts string to integer. Strips common currency symbols (`$`, commas).

```json
{"price": {"xpath": ".//span[@class='price']/text()", "pipes": ["toint"]}}
```

**Input:** `"$1,234"` → **Output:** `1234` (int)

#### tofloat

Converts string to float64. Strips currency symbols.

```json
{"price": {"xpath": ".//span[@class='price']/text()", "pipes": ["tofloat"]}}
```

**Input:** `"$1,234.56"` → **Output:** `1234.56` (float64)

#### parseurl

Converts relative URLs to absolute using base URL from context.

```json
{"link": {"xpath": ".//a/@href", "pipes": ["parseurl"]}}
```

**Requires:** Base URL in context (automatically set when using `ScrapeURL`)

**Input:** `"/products/item"` with base `https://example.com` → **Output:** `https://example.com/products/item`

#### parsetime

Parses datetime string with specified layout and timezone.

```json
{"date": {"xpath": ".//time/@datetime", "pipes": ["parsetime:2006-01-02T15:04:05Z:UTC"]}}
```

**Parameters:**
1. `layout` - Go time format (required)
2. `timezone` - IANA timezone name (optional, default: UTC)

#### regexreplace

Performs regex substitution.

```json
{"clean": {"xpath": ".//text()", "pipes": ["regexreplace:\\\\s+:_:i"]}}
```

**Parameters:**
1. `pattern` - Regex pattern (required)
2. `replacement` - Replacement string (required)
3. `flags` - Optional flags (only `i` for case-insensitive supported)

#### humanduration

Converts seconds to human-readable "X ago" format.

```json
{"ago": {"xpath": ".//time/@data-seconds", "pipes": ["humanduration"]}}
```

**Input:** `"120"` → **Output:** `"2 minutes ago"`

### Custom Pipes

Register custom pipes using `RegisterPipe`:

```go
import (
    "context"
    "strings"
    "github.com/Hanivan/gtmlp"
)

func init() {
    gtmlp.RegisterPipe("uppercase", func(ctx context.Context, input string, params []string) (any, error) {
        return strings.ToUpper(input), nil
    })

    gtmlp.RegisterPipe("slugify", func(ctx context.Context, input string, params []string) (any, error) {
        // Custom slugification logic
        slug := strings.ToLower(strings.ReplaceAll(input, " ", "-"))
        return slug, nil
    })
}
```

**Pipe Function Signature:**

```go
type PipeFunc func(ctx context.Context, input string, params []string) (any, error)
```

**Parameters:**
- `ctx` - Context (contains baseURL for parseurl pipe)
- `input` - String value from XPath extraction
- `params` - Pipe parameters (split by `:`)

**Returns:**
- `any` - Transformed value (can be string, int, float64, time.Time, etc.)
- `error` - Error if transformation fails

**Example with parameters:**

```go
gtmlp.RegisterPipe("prefix", func(ctx context.Context, input string, params []string) (any, error) {
    if len(params) < 1 {
        return "", fmt.Errorf("prefix requires parameter")
    }
    return params[0] + input, nil
})

// Usage in config:
{"name": {"xpath": ".//h2/text()", "pipes": ["prefix:Product: "]}}
```

### Pipe Chains

Pipes are applied in order. Each pipe receives the output of the previous pipe:

```json
{"price": {"xpath": ".//span/text()", "pipes": ["trim", "regexReplace:\\$::", "tofloat"]}}
```

**Flow:** `"$1,234.56"` → `trim` → `"1,234.56"` → `regexReplace` → `"1,234.56"` → `tofloat` → `1234.56`

### Error Handling

Pipe errors return `ErrTypePipe` with context:

```go
results, err := gtmlp.ScrapeURL[Product](context.Background(), url, config)
if err != nil {
    if gtmlp.Is(err, gtmlp.ErrTypePipe) {
        pipeErr := err.(*gtmlp.PipeError)
        log.Printf("Pipe '%s' failed on field '%s': %v",
            pipeErr.PipeName, pipeErr.Field, pipeErr.Cause)
    }
}
```

## XPath Validation

### ValidateXPath

Validates XPath expressions against HTML content.

```go
func ValidateXPath(html string, xpaths map[string]string) map[string]ValidationResult
```

**Parameters:**
- `html`: HTML content to test against
- `xpaths`: Map of field names to XPath expressions

**Returns:**
- `map[string]ValidationResult`: Validation results for each field

**Example:**

```go
html := `<html><body><div class="product"><h2>Product A</h2></div></body></html>`

xpaths := map[string]string{
    "container": "//div[@class='product']",
    "name":       ".//h2/text()",
    "price":      ".//span[@class='price']/text()",  // This will fail
}

results := gtmlp.ValidateXPath(html, xpaths)

for field, result := range results {
    if result.Valid {
        fmt.Printf("%s: Valid (%d matches)\n", field, result.MatchCount)
    } else {
        fmt.Printf("%s: Invalid - %v\n", field, result.Error)
    }
}
```

**Output:**
```
container: Valid (1 matches)
name: Valid (1 matches)
price: Invalid - no matches found
```

### ValidateXPathURL

Validates XPath expressions by fetching from a URL.

```go
func ValidateXPathURL(url string, config *Config) (map[string]ValidationResult, error)
```

**Parameters:**
- `url`: URL to fetch and validate against
- `config`: Configuration with XPath expressions

**Returns:**
- `map[string]ValidationResult`: Validation results
- `error`: Error if fetching fails

**Example:**

```go
config, _ := gtmlp.LoadConfig("selectors.json", nil)

results, err := gtmlp.ValidateXPathURL(
    context.Background(),
    "https://example.com/products",
    config,
)
if err != nil {
    log.Fatal(err)
}

for field, result := range results {
    if !result.Valid {
        log.Printf("Field '%s' failed: %v", field, result.Error)
    }
}
```

## Health Check

### CheckHealth

Performs a health check on a single URL.

```go
func CheckHealth(url string) HealthCheckResult
```

**Parameters:**
- `url`: URL to check

**Returns:**
- `HealthCheckResult`: Health check result

**Example:**

```go
result := gtmlp.CheckHealth("https://example.com")

if result.Status == gtmlp.StatusHealthy {
    fmt.Printf("URL is healthy (200 OK)\n")
} else {
    fmt.Printf("URL is %s: %v\n", result.Status, result.Error)
}

fmt.Printf("Status code: %d\n", result.Code)
fmt.Printf("Latency: %v\n", result.Latency)
```

### CheckHealthMulti

Performs health checks on multiple URLs concurrently.

```go
func CheckHealthMulti(urls []string) []HealthCheckResult
```

**Parameters:**
- `urls`: Slice of URLs to check

**Returns:**
- `[]HealthCheckResult`: Health check results (same order as input)

**Example:**

```go
urls := []string{
    "https://example.com",
    "https://api.example.com/health",
    "https://unknown-domain-12345.com",
}

results := gtmlp.CheckHealthMulti(urls)

for i, result := range results {
    fmt.Printf("%s: %s (%d) - %v\n",
        result.URL,
        result.Status,
        result.Code,
        result.Latency,
    )

    if result.Error != nil {
        fmt.Printf("  Error: %v\n", result.Error)
    }
}
```

**Output:**
```
https://example.com: healthy (200) - 45ms
https://api.example.com/health: healthy (200) - 123ms
https://unknown-domain-12345.com: error (0) - 2ms
  Error: network error: HTTP request failed (url: https://unknown-domain-12345.com)
```

### CheckHealthWithOptions

Performs health check with custom configuration.

```go
func CheckHealthWithOptions(url string, config *Config) HealthCheckResult
```

**Parameters:**
- `url`: URL to check
- `config`: Configuration (uses timeout, user agent, proxy, headers)

**Returns:**
- `HealthCheckResult`: Health check result

**Example:**

```go
config := &gtmlp.Config{
    Timeout:   5 * time.Second,
    UserAgent: "MyHealthCheck/1.0",
    Headers: map[string]string{
        "Authorization": "Bearer token123",
    },
}

result := gtmlp.CheckHealthWithOptions("https://api.example.com", config)
```

## Types

### Config

Scraping configuration structure.

```go
type Config struct {
    // XPath definitions
    Container string                    // Repeating element selector
    Fields    map[string]FieldConfig    // Field name → Field configuration

    // HTTP options
    Timeout    time.Duration
    UserAgent  string
    RandomUA   bool
    MaxRetries int
    Proxy      string
    Headers    map[string]string
}
```

### FieldConfig

Field configuration with XPath and optional pipes.

```go
type FieldConfig struct {
    XPath string   // XPath expression
    Pipes []string // Pipe chain (e.g., ["trim", "tofloat"])
}
```

### PipeFunc

Pipe transformation function signature.

```go
type PipeFunc func(ctx context.Context, input string, params []string) (any, error)
```

**Field Descriptions:**
- `Container`: XPath expression for repeating container elements
- `Fields`: Map of field names to relative XPath expressions
- `Timeout`: HTTP request timeout (default: 30s)
- `UserAgent`: HTTP User-Agent header (default: "GTMLP/2.0")
- `RandomUA`: Use random user-agent (default: false)
- `MaxRetries`: Number of retries (default: 0)
- `Proxy`: HTTP proxy URL
- `Headers`: Additional HTTP headers

### EnvMapping

Environment variable name mapping.

```go
type EnvMapping struct {
    Timeout    string
    UserAgent  string
    RandomUA   string
    MaxRetries string
    Proxy      string
}
```

**Default Mapping:**
```go
DefaultEnvMapping = &EnvMapping{
    Timeout:    "GTMLP_TIMEOUT",
    UserAgent:  "GTMLP_USER_AGENT",
    RandomUA:   "GTMLP_RANDOM_UA",
    MaxRetries: "GTMLP_MAX_RETRIES",
    Proxy:      "GTMLP_PROXY",
}
```

### HealthStatus

Health check status enumeration.

```go
type HealthStatus int

const (
    StatusHealthy   HealthStatus = iota  // 2xx status codes
    StatusUnhealthy                      // 4xx/5xx/3xx status codes
    StatusError                          // Network or other errors
)
```

**String Representation:**
```go
func (s HealthStatus) String() string
```

**Example:**
```go
status := gtmlp.StatusHealthy
fmt.Println(status)  // Output: "healthy"
```

### HealthCheckResult

Health check result structure.

```go
type HealthCheckResult struct {
    URL     string        // The URL that was checked
    Status  HealthStatus  // The health status
    Code    int           // HTTP status code (0 if error)
    Latency time.Duration // Time taken for the check
    Error   error         // Error message if check failed
}
```

### ValidationResult

XPath validation result structure.

```go
type ValidationResult struct {
    XPath      string // The XPath expression
    Valid      bool   // Whether the XPath is valid
    MatchCount int    // Number of matches found
    Error      error  // Error if validation failed
}
```

### ScrapeError

Typed error with context.

```go
type ScrapeError struct {
    Type    ErrorType
    Message string
    XPath   string
    URL     string
    Cause   error
}
```

**Error Types:**
```go
const (
    ErrTypeNetwork    ErrorType = "network"
    ErrTypeParsing    ErrorType = "parsing"
    ErrTypeXPath      ErrorType = "xpath"
    ErrTypeConfig     ErrorType = "config"
    ErrTypeValidation ErrorType = "validation"
    ErrTypePipe       ErrorType = "pipe"
)
```

**Methods:**
```go
func (e *ScrapeError) Error() string
func (e *ScrapeError) Unwrap() error
```

**Error Type Checking:**
```go
func Is(err error, errorType ErrorType) bool
```

**Example:**
```go
results, err := gtmlp.ScrapeURL[Product](url, config)
if err != nil {
    if gtmlp.Is(err, gtmlp.ErrTypeNetwork) {
        log.Printf("Network error: %v", err)
    } else if gtmlp.Is(err, gtmlp.ErrTypeXPath) {
        log.Printf("XPath error: %v", err)
    }
}
```

## Error Handling

All v2.0 API functions return `*ScrapeError` for structured error handling.

### Error Type Checking

```go
results, err := gtmlp.ScrapeURL[Product](url, config)
if err != nil {
    switch {
    case gtmlp.Is(err, gtmlp.ErrTypeNetwork):
        // Handle network errors (timeout, connection failed, etc.)
    case gtmlp.Is(err, gtmlp.ErrTypeParsing):
        // Handle HTML parsing errors
    case gtmlp.Is(err, gtmlp.ErrTypeXPath):
        // Handle XPath syntax errors
    case gtmlp.Is(err, gtmlp.ErrTypeConfig):
        // Handle configuration errors
    default:
        // Unknown error
    }
}
```

### Error Context

`ScrapeError` provides rich context:

```go
if err != nil {
    var scrapeErr *gtmlp.ScrapeError
    if errors.As(err, &scrapeErr) {
        log.Printf("Error type: %s", scrapeErr.Type)
        log.Printf("Message: %s", scrapeErr.Message)
        if scrapeErr.URL != "" {
            log.Printf("URL: %s", scrapeErr.URL)
        }
        if scrapeErr.XPath != "" {
            log.Printf("XPath: %s", scrapeErr.XPath)
        }
        if scrapeErr.Cause != nil {
            log.Printf("Cause: %v", scrapeErr.Cause)
        }
    }
}
```

## Complete Examples

### Example 1: E-commerce Product Scraping

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/Hanivan/gtmlp"
)

type Product struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Link  string  `json:"link"`
}

func main() {
    // Load config
    config, err := gtmlp.LoadConfig("products.json", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Scrape products
    products, err := gtmlp.ScrapeURL[Product](
        context.Background(),
        "https://example.com/shop",
        config,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Process results
    for _, p := range products {
        fmt.Printf("%s - %.2f (%s)\n", p.Name, p.Price, p.Link)
    }
}
```

**products.json:**
```json
{
  "container": "//div[@class='product-item']",
  "fields": {
    "name": {"xpath": ".//h3[@class='product-title']/text()", "pipes": ["trim"]},
    "price": {"xpath": ".//span[@class='price']/text()", "pipes": ["trim", "tofloat"]},
    "link": {"xpath": ".//a[@class='product-link']/@href", "pipes": ["parseurl"]}
  },
  "timeout": "30s",
  "user_agent": "MyBot/1.0",
  "random_ua": true,
  "max_retries": 3
}
```

### Example 2: Article Scraping with Validation

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/Hanivan/gtmlp"
)

type Article struct {
    Title   string `json:"title"`
    Author  string `json:"author"`
    Content string `json:"content"`
}

func main() {
    // Load config
    config, err := gtmlp.LoadConfig("articles.yaml", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Validate XPath before scraping
    url := "https://blog.example.com"
    results, err := gtmlp.ValidateXPathURL(context.Background(), url, config)
    if err != nil {
        log.Fatal(err)
    }

    // Check validation results
    for field, result := range results {
        if !result.Valid {
            log.Printf("Warning: Field '%s' validation failed: %v", field, result.Error)
        } else if result.MatchCount == 0 {
            log.Printf("Warning: Field '%s' has no matches", field)
        }
    }

    // Scrape if validation passed
    articles, err := gtmlp.ScrapeURL[Article](context.Background(), url, config)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Scraped %d articles\n", len(articles))
}
```

### Example 3: Batch URL Health Check

```go
package main

import (
    "fmt"
    "log"
    "github.com/Hanivan/gtmlp"
)

func main() {
    urls := []string{
        "https://example.com",
        "https://api.example.com/health",
        "https://blog.example.com",
        "https://invalid-url-12345.com",
    }

    results := gtmlp.CheckHealthMulti(urls)

    healthy := 0
    unhealthy := 0
    errors := 0

    for _, result := range results {
        switch result.Status {
        case gtmlp.StatusHealthy:
            healthy++
            fmt.Printf("✓ %s - %d (%v)\n", result.URL, result.Code, result.Latency)
        case gtmlp.StatusUnhealthy:
            unhealthy++
            fmt.Printf("✗ %s - %d (%v)\n", result.URL, result.Code, result.Latency)
        case gtmlp.StatusError:
            errors++
            fmt.Printf("⚠ %s - Error: %v\n", result.URL, result.Error)
        }
    }

    fmt.Printf("\nSummary: %d healthy, %d unhealthy, %d errors\n",
        healthy, unhealthy, errors)
}
```

### Example 4: Environment Variable Configuration

```go
package main

import (
    "context"
    "log"
    "os"
    "github.com/Hanivan/gtmlp"
)

func main() {
    // Set environment variables (or set them externally)
    os.Setenv("GTMLP_TIMEOUT", "60s")
    os.Setenv("GTMLP_USER_AGENT", "MyBot/2.0")
    os.Setenv("GTMLP_RANDOM_UA", "true")
    os.Setenv("GTMLP_MAX_RETRIES", "5")

    // Load config (env vars will override file settings)
    config, err := gtmlp.LoadConfig("selectors.json", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Config now has environment variable values
    // Timeout: 60s (from env)
    // UserAgent: "MyBot/2.0" (from env)
    // RandomUA: true (from env)
    // MaxRetries: 5 (from env)

    results, err := gtmlp.ScrapeURL[Product](context.Background(), "https://example.com", config)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Example 5: Dynamic Scraping with Untyped Results

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/Hanivan/gtmlp"
)

func main() {
    config := &gtmlp.Config{
        Container: "//div[@class='item']",
        Fields: map[string]gtmlp.FieldConfig{
            "title":       {XPath: ".//h2/text()"},
            "description": {XPath: ".//p/text()"},
            "tags":        {XPath: ".//span[@class='tag']/text()"},
        },
    }

    // Use untyped scraping for dynamic data
    results, err := gtmlp.ScrapeURLUntyped(
        context.Background(),
        "https://example.com/items",
        config,
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, item := range results {
        fmt.Printf("Title: %v\n", item["title"])
        fmt.Printf("Description: %v\n", item["description"])
        fmt.Printf("Tags: %v\n", item["tags"])
        fmt.Println("---")
    }
}
```

### Example 6: Proxy Configuration

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/Hanivan/gtmlp"
)

func main() {
    config := &gtmlp.Config{
        Container: "//div[@class='product']",
        Fields: map[string]gtmlp.FieldConfig{
            "name": {XPath: ".//h2/text()"},
        },
        Timeout:   30 * time.Second, // Timeout applies to proxy too
        Proxy:     "http://proxy.example.com:8080",
        UserAgent: "MyBot/1.0",
        Headers: map[string]string{
            "Proxy-Authentication": "Basic base64encoded",
        },
    }

    products, err := gtmlp.ScrapeURL[Product](
        context.Background(),
        "https://example.com",
        config,
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

### Example 7: Error Handling with Retry

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/Hanivan/gtmlp"
)

func main() {
    config := &gtmlp.Config{
        Container: "//div[@class='product']",
        Fields: map[string]gtmlp.FieldConfig{
            "name": {XPath: ".//h2/text()"},
        },
        Timeout:    10 * time.Second,
        MaxRetries: 3, // Retry up to 3 times with exponential backoff
        RandomUA:   true,
    }

    products, err := gtmlp.ScrapeURL[Product](
        context.Background(),
        "https://example.com",
        config,
    )

    if err != nil {
        // Check error type
        if gtmlp.Is(err, gtmlp.ErrTypeNetwork) {
            log.Printf("Network error after retries: %v", err)
        } else if gtmlp.Is(err, gtmlp.ErrTypeXPath) {
            log.Printf("Invalid XPath in config: %v", err)
        } else if gtmlp.Is(err, gtmlp.ErrTypePipe) {
            log.Printf("Pipe transformation failed: %v", err)
        } else {
            log.Printf("Unexpected error: %v", err)
        }
        return
    }

    log.Printf("Successfully scraped %d products", len(products))
}
```

## Best Practices

### 1. Use Type Parameters When Possible

```go
// Good - Type safe
products, err := gtmlp.ScrapeURL[Product](context.Background(), url, config)

// Acceptable for dynamic data
results, err := gtmlp.ScrapeURLUntyped(context.Background(), url, config)
```

### 2. Validate XPath Before Scraping

```go
// Always validate in development
results, err := gtmlp.ValidateXPathURL(context.Background(), url, config)
if err != nil {
    log.Fatal(err)
}

for field, result := range results {
    if !result.Valid || result.MatchCount == 0 {
        log.Printf("Issue with field '%s'", field)
    }
}
```

### 3. Use Configuration Files

Store XPath selectors in JSON/YAML files for maintainability:

```json
{
  "container": "//div[@class='product']",
  "fields": {
    "name": {"xpath": ".//h2/text()", "pipes": ["trim"]}
  }
}
```

### 4. Set Appropriate Timeouts

```go
config := &gtmlp.Config{
    Timeout: 30 * time.Second, // Adjust based on target site
    // ...
}
```

### 5. Handle Errors Gracefully

```go
results, err := gtmlp.ScrapeURL[Product](context.Background(), url, config)
if err != nil {
    if gtmlp.Is(err, gtmlp.ErrTypeNetwork) {
        // Retry with backoff
    } else if gtmlp.Is(err, gtmlp.ErrTypeConfig) {
        // Fix configuration
    } else if gtmlp.Is(err, gtmlp.ErrTypePipe) {
        // Fix pipe configuration
    }
    return
}
```

### 6. Use Health Checks for Monitoring

```go
results := gtmlp.CheckHealthMulti(urls)
for _, result := range results {
    if result.Status != gtmlp.StatusHealthy {
        alert(result.URL, result.Status, result.Error)
    }
}
```

### 7. Leverage Environment Variables

```bash
# Production
export GTMLP_TIMEOUT=60s
export GTMLP_MAX_RETRIES=5
export GTMLP_PROXY=http://proxy.prod:8080

# Development
export GTMLP_TIMEOUT=10s
export GTMLP_MAX_RETRIES=1
```

## Migration from v1.x

The v1.x API remains fully functional. Key differences:

### v1.x Style (Programmatic)

```go
patterns := []gtmlp.PatternField{
    gtmlp.NewContainerPattern("container", "//div[@class='product']"),
    {Key: "name", Patterns: []string{".//h2/text()"}},
    {Key: "price", Patterns: []string{".//span[@class='price']/text()"}},
}

results, err := gtmlp.FromURL("https://example.com").
    WithPatterns(patterns).
    Extract()
```

### v2.0 Style (Configuration-based)

```go
config, _ := gtmlp.LoadConfig("selectors.json", nil)
products, err := gtmlp.ScrapeURL[Product]("https://example.com", config)
```

Both APIs can coexist in the same application. Choose based on your needs:

- **v1.x**: Dynamic scraping, complex patterns, pipe transformations
- **v2.0**: Static scraping, external config, type safety

## See Also

- [README.md](../README.md) - Main project documentation
- [API.md](API.md) - v1.x API reference
- [PATTERNS.md](PATTERNS.md) - Pattern extraction guide
- [EXAMPLES.md](EXAMPLES.md) - Usage examples
