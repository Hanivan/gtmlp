# GTMLP v2.0 API Reference

Complete API reference for GTMLP v2.0 - the configuration-based, type-safe scraping API.

## Table of Contents

- [Overview](#overview)
- [Core Scraping Functions](#core-scraping-functions)
- [Config Loading](#config-loading)
- [Logging](#logging)
- [Security](#security)
- [Fallback XPath Chains](#fallback-xpath-chains)
- [Pagination](#pagination)
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

## Logging

GTMLP uses Go's standard `log/slog` package for structured logging with configurable log levels.

### SetLogLevel

Sets the global log level for GTMLP operations.

```go
func SetLogLevel(level slog.Level)
```

**Parameters:**
- `level`: Log level (LevelDebug, LevelInfo, LevelWarn, LevelError)

**Example:**

```go
import "log/slog"

// Development: see HTTP requests and scraping progress
gtmlp.SetLogLevel(slog.LevelInfo)

// Troubleshooting: see detailed XPath evaluation
gtmlp.SetLogLevel(slog.LevelDebug)

// Production: warnings and errors only (default)
gtmlp.SetLogLevel(slog.LevelWarn)

// Silent: errors only
gtmlp.SetLogLevel(slog.LevelError)
```

**Note:** `SetLogLevel` creates a new default handler (TextHandler to stderr). For custom handlers, use `SetLogger` instead.

### SetLogger

Sets a custom logger with full control over output format and destination.

```go
func SetLogger(logger *slog.Logger)
```

**Parameters:**
- `logger`: Custom slog.Logger instance

**Example:**

```go
import (
    "log/slog"
    "os"
)

// JSON format to stdout
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})
gtmlp.SetLogger(slog.New(handler))

// Custom writer (e.g., file)
file, _ := os.OpenFile("gtmlp.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
handler := slog.NewTextHandler(file, &slog.HandlerOptions{
    Level: slog.LevelDebug,
})
gtmlp.SetLogger(slog.New(handler))
```

### GetLogger

Returns the current global logger instance.

```go
func GetLogger() *slog.Logger
```

**Returns:**
- `*slog.Logger`: Current logger

**Example:**

```go
logger := gtmlp.GetLogger()
logger.Info("custom log message", "key", "value")
```

### Log Levels

| Level | Description | What's Logged |
|-------|-------------|---------------|
| **Debug** | Detailed debugging | XPath evaluation, fallback usage, field extraction details |
| **Info** | General information | HTTP requests, scraping progress, pagination status |
| **Warn** | Warnings (default) | Fallback XPath used, HTTP warnings, duplicate URLs |
| **Error** | Errors only | HTTP failures, parsing errors, validation failures |

### Log Output Examples

**Debug Level:**
```
time=2026-02-04T10:30:00Z level=DEBUG msg="scraping started" container="//div" fields=3 html_size=15234
time=2026-02-04T10:30:00Z level=DEBUG msg="containers found" xpath="//div[@class='product']"
time=2026-02-04T10:30:00Z level=WARN msg="field fallback used" primary=".//h2" used=".//h3" fallback_index=1
```

**Info Level:**
```
time=2026-02-04T10:30:00Z level=INFO msg="http request successful" url="https://example.com" status=200 duration_ms=342 attempt=1
time=2026-02-04T10:30:00Z level=INFO msg="scraping completed" items=25 container="//div"
time=2026-02-04T10:30:00Z level=INFO msg="pagination starting" url="https://example.com/page1" type="next-link"
```

**Warn Level (default):**
```
time=2026-02-04T10:30:00Z level=WARN msg="http url used" url="http://example.com" recommendation="use_https"
time=2026-02-04T10:30:00Z level=WARN msg="container fallback used" primary="//div[@class='a']" used="//div[@class='b']"
time=2026-02-04T10:30:00Z level=WARN msg="pagination duplicate url" url="https://example.com/page1" page=5
```

### Best Practices

**Development:**
```go
gtmlp.SetLogLevel(slog.LevelInfo) // See requests and progress
```

**Staging/Testing:**
```go
gtmlp.SetLogLevel(slog.LevelDebug) // Full visibility for debugging
```

**Production:**
```go
// Use default (Warn level) or explicitly set
gtmlp.SetLogLevel(slog.LevelWarn) // Warnings and errors only
```

**Custom Production Logger:**
```go
// JSON logs to stderr for log aggregation
handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelWarn,
    AddSource: true, // Include source file/line
})
gtmlp.SetLogger(slog.New(handler))
```

## Security

GTMLP includes built-in security features to prevent common web scraping vulnerabilities.

### SSRF Protection

Server-Side Request Forgery (SSRF) protection is enabled by default, blocking requests to private IP ranges.

**Blocked by default:**
- Localhost: `127.0.0.0/8`, `::1`
- Private IPv4: `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`
- Link-local: `169.254.0.0/16` (AWS metadata service)

**Configuration:**

```go
config := &gtmlp.Config{
    Container: "//div[@class='product']",
    Fields:    fields,

    // SSRF protection (default: enabled)
    AllowPrivateIPs: false, // set true to allow private IPs

    // Other settings...
    Timeout: 30 * time.Second,
}
```

**Example - Allow Private IPs:**

```go
// Development/testing: allow localhost
config := &gtmlp.Config{
    Container:       "//div",
    Fields:          fields,
    AllowPrivateIPs: true, // Allow private IPs
}

products, _ := gtmlp.ScrapeURL[Product](ctx, "http://localhost:8080", config)
```

**SSRF Error Example:**

```go
products, err := gtmlp.ScrapeURL[Product](ctx, "http://169.254.169.254/", config)
// Error: network error: URL validation failed
//   SSRF protection: link-local IP access blocked (169.254.169.254)
```

### Custom URL Validator

Add custom URL validation logic using `URLValidator`:

```go
config := &gtmlp.Config{
    Container: "//div",
    Fields:    fields,

    // Custom validator (runs after SSRF check)
    URLValidator: func(url string) error {
        // Only allow specific domains
        allowedDomains := []string{"example.com", "api.example.com"}

        u, _ := url.Parse(url)
        for _, domain := range allowedDomains {
            if u.Host == domain || strings.HasSuffix(u.Host, "."+domain) {
                return nil
            }
        }

        return fmt.Errorf("domain %s not in allowlist", u.Host)
    },
}
```

**Common Use Cases:**

**1. Domain Allowlist:**
```go
URLValidator: func(url string) error {
    if !strings.Contains(url, "trusted-domain.com") {
        return errors.New("domain not allowed")
    }
    return nil
}
```

**2. Protocol Restriction:**
```go
URLValidator: func(url string) error {
    if !strings.HasPrefix(url, "https://") {
        return errors.New("only HTTPS allowed")
    }
    return nil
}
```

**3. Rate Limiting Check:**
```go
URLValidator: func(url string) error {
    if isRateLimited(url) {
        return errors.New("rate limit exceeded")
    }
    return nil
}
```

### Security Best Practices

**1. Always use HTTPS in production:**
```go
URLValidator: func(url string) error {
    if strings.HasPrefix(url, "http://") {
        return errors.New("HTTP not allowed in production")
    }
    return nil
}
```

**2. Validate user-provided URLs:**
```go
// Never scrape user input without validation
userURL := getUserInput()

// Validate before scraping
config.URLValidator = func(url string) error {
    if !isValidDomain(url) {
        return errors.New("invalid domain")
    }
    return nil
}

products, _ := gtmlp.ScrapeURL[Product](ctx, userURL, config)
```

**3. Set timeouts:**
```go
config := &gtmlp.Config{
    Timeout: 30 * time.Second, // Prevent hanging
    // ...
}
```

**4. Limit pagination:**
```go
config := &gtmlp.Config{
    Pagination: &gtmlp.PaginationConfig{
        MaxPages: 50,              // Limit pages
        Timeout:  10 * time.Minute, // Total timeout
    },
}
```

**5. Monitor for SSRF attempts:**
```go
products, err := gtmlp.ScrapeURL[Product](ctx, url, config)
if err != nil && strings.Contains(err.Error(), "SSRF protection") {
    log.Warn("SSRF attempt detected", "url", url, "error", err)
    metrics.IncrementSSRFAttempts()
}
```

For more security information, see [SECURITY.md](../SECURITY.md).

## Fallback XPath Chains

GTMLP supports fallback XPath chains using `altXpath` and `altContainer` to handle varying HTML structures gracefully. This is useful when:

- Different pages have different HTML structures
- Elements can appear in multiple locations
- Handling A/B tests or template variations
- Dealing with legacy and modern page layouts

### How It Works

**Field-Level Fallback (`altXpath`):**
1. Try the primary `xpath`
2. Apply `pipes` to the result
3. If result is empty after pipes, try each `altXpath` in order
4. Return first non-empty result
5. Return empty string if all fail

**Container-Level Fallback (`altContainer`):**
1. Try the primary `container` XPath
2. If no elements found, try each `altContainer` in order
3. Use first non-empty container list
4. Return empty array if all fail

### Configuration

**JSON:**
```json
{
  "container": "//div[@class='product']",
  "altContainer": [
    "//article[@class='product']",
    "//div[@class='item']"
  ],
  "fields": {
    "name": {
      "xpath": ".//h2[@class='title']/text()",
      "altXpath": [
        ".//h3[@class='product-name']/text()",
        ".//h1/text()"
      ],
      "pipes": ["trim"]
    },
    "price": {
      "xpath": ".//span[@class='price']/text()",
      "altXpath": [".//div[@class='price']/text()"],
      "pipes": ["trim", "tofloat"]
    }
  }
}
```

**YAML:**
```yaml
container: //div[@class='product']
altContainer:
  - //article[@class='product']
  - //div[@class='item']

fields:
  name:
    xpath: .//h2[@class='title']/text()
    altXpath:
      - .//h3[@class='product-name']/text()
      - .//h1/text()
    pipes: [trim]

  price:
    xpath: .//span[@class='price']/text()
    altXpath:
      - .//div[@class='price']/text()
    pipes: [trim, tofloat]
```

### Behavior with Pipes

Pipes are applied to **each XPath attempt** before checking if the result is empty:

```go
// Example: field with whitespace-only content
// HTML: <h2>   </h2><h3>Product Name</h3>

config := &Config{
    Container: "//div",
    Fields: map[string]FieldConfig{
        "name": {
            XPath:    ".//h2/text()",          // Returns "   "
            AltXPath: []string{".//h3/text()"}, // Returns "Product Name"
            Pipes:    []string{"trim"},
        },
    },
}

// Process:
// 1. Extract ".//h2/text()" → "   "
// 2. Apply "trim" pipe → ""
// 3. Result is empty, try altXpath
// 4. Extract ".//h3/text()" → "Product Name"
// 5. Apply "trim" pipe → "Product Name"
// 6. Result is non-empty, return "Product Name"
```

This allows you to:
- Filter out whitespace-only elements
- Reject placeholder text (e.g., "N/A", "TBD")
- Validate data before accepting it as the final result

### Empty vs Error

**Important:** Fallback is triggered by **empty results**, not XPath syntax errors.

- ✅ **Empty result** → Try next XPath in fallback chain
- ❌ **XPath syntax error** → Fail immediately (configuration bug)

This design ensures:
- Syntax errors are caught during validation
- Debugging is easier (errors fail fast)
- Fallback is only used for legitimate structural variations

### Examples

**Use Case 1: Different Title Tags**
```json
{
  "fields": {
    "title": {
      "xpath": ".//h2[@class='title']/text()",
      "altXpath": [".//h1[@class='heading']/text()", ".//div[@class='title']/text()"]
    }
  }
}
```

**Use Case 2: Multiple Price Formats**
```json
{
  "fields": {
    "price": {
      "xpath": ".//span[@class='sale-price']/text()",
      "altXpath": [".//span[@class='regular-price']/text()"],
      "pipes": ["trim", "tofloat"]
    }
  }
}
```

**Use Case 3: Container Variations**
```json
{
  "container": "//div[@class='product-grid']//div[@class='product']",
  "altContainer": [
    "//ul[@class='products']//li[@class='product']",
    "//section[@class='catalog']//article"
  ]
}
```

## Pagination

GTMLP supports automatic pagination to scrape multi-page listings.

### Pagination Types

**Next-Link Pagination** - Follow "Next" links sequentially:
```json
{
  "pagination": {
    "type": "next-link",
    "nextSelector": "//a[@rel='next']/@href",
    "altSelectors": ["//a[contains(@class, 'next')]/@href"],
    "maxPages": 50
  }
}
```

**Numbered Pagination** - Extract all page links:
```json
{
  "pagination": {
    "type": "numbered",
    "pageSelector": "//div[@class='pagination']//a/@href",
    "maxPages": 20
  }
}
```

### Usage Modes

**Auto-Follow** (combined results):
```go
config, _ := gtmlp.LoadConfig("selectors.json", nil)

// Returns all items from all pages in a single array
products, err := gtmlp.ScrapeURL[Product](ctx, "https://example.com/products", config)
```

**Page-Separated** (with metadata):
```go
// Returns results grouped by page
results, err := gtmlp.ScrapeURLWithPages[Product](ctx, url, config)

for _, page := range results.Pages {
    fmt.Printf("Page %d: %d items at %s\n", page.PageNum, len(page.Items), page.URL)
}
```

**Extract-Only** (manual control):
```go
// Get pagination URLs without scraping
info, err := gtmlp.ExtractPaginationURLs(ctx, url, config)

for _, pageURL := range info.URLs {
    // Scrape each page manually
    items, _ := gtmlp.ScrapeURL[Product](ctx, pageURL, config)
}
```

### Pagination Features

- **Fallback selectors** - `altSelectors` like `altXpath`
- **Pipe support** - Transform pagination URLs
- **Duplicate detection** - Prevents circular references with warnings
- **Relative URL resolution** - Auto-convert relative → absolute URLs
- **Safety limits** - `maxPages` (default: 100), `timeout` (default: 10m)
- **Progress logging** - Use `SetLogLevel(slog.LevelInfo)` to see pagination progress
- **Error handling** - Returns partial results on failure

### Configuration Reference

```go
type PaginationConfig struct {
    Type          string        // "next-link" or "numbered"
    NextSelector  string        // XPath for next link (next-link)
    AltSelectors  []string      // Fallback selectors
    PageSelector  string        // XPath for all pages (numbered)
    Pipes         []string      // URL transformation pipes
    MaxPages      int           // Max pages (default: 100)
    Timeout       time.Duration // Total timeout (default: 10m)
}
```

**Logging Pagination Progress:**

```go
import "log/slog"

// Enable pagination logging
gtmlp.SetLogLevel(slog.LevelInfo)

// Logs will show:
// - "pagination starting" - When pagination begins
// - "pagination page scraped" - After each page
// - "pagination completed" - When finished
// - "pagination duplicate url" - If duplicate detected

config, _ := gtmlp.LoadConfig("selectors.json", nil)
products, _ := gtmlp.ScrapeURL[Product](ctx, url, config)
```

### Examples

See working examples:
- **[pagination_next_json](../examples/v2/pagination_next_json)** - Next-link pagination
- **[pagination_numbered_yaml](../examples/v2/pagination_numbered_yaml)** - Numbered pagination

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
    Container    string                    // Repeating element selector
    AltContainer []string                  // Alternative container selectors (fallback)
    Fields       map[string]FieldConfig    // Field name → Field configuration

    // HTTP options
    Timeout    time.Duration
    UserAgent  string
    RandomUA   bool
    MaxRetries int
    Proxy      string
    Headers    map[string]string

    // Security options
    URLValidator    func(string) error    // Custom URL validator
    AllowPrivateIPs bool                  // Allow private IPs (default: false)

    // Pagination options
    Pagination *PaginationConfig          // Pagination configuration
}
```

### FieldConfig

Field configuration with XPath and optional pipes.

```go
type FieldConfig struct {
    XPath    string   // XPath expression
    AltXPath []string // Alternative XPath expressions (fallback)
    Pipes    []string // Pipe chain (e.g., ["trim", "tofloat"])
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
