# GTMLP

Type-safe HTML scraping with XPath selectors and external configuration.

## Features

- **Type-safe** with Go generics
- **External config** (JSON/YAML)
- **XPath validation** before scraping
- **Pagination support** - Auto-follow next-link or numbered pagination
- **Fallback XPath chains** (`altXpath`, `altContainer`) for handling varying HTML structures
- **Data transformation pipes** (trim, int/float conversion, regex, URL parsing, etc.)
- **Custom pipe registration** for domain-specific transformations
- **Structured logging** with configurable levels (slog-based)
- **SSRF protection** - Blocks private IPs by default
- **Health checks** for URLs
- **Production ready** (retries, proxy, timeouts)

## Installation

```bash
go get github.com/Hanivan/gtmlp
```

## Quick Start

**selectors.json:**
```json
{
  "container": "//div[@class='product']",
  "fields": {
    "name": {"xpath": ".//h2/text()"},
    "price": {"xpath": ".//span[@class='price']/text()", "pipes": ["trim", "tofloat"]}
  }
}
```

**main.go:**
```go
type Product struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

config, _ := gtmlp.LoadConfig("selectors.json", nil)
products, _ := gtmlp.ScrapeURL[Product](context.Background(), "https://example.com", config)

for _, p := range products {
    fmt.Printf("%s: %.2f\n", p.Name, p.Price)
}
```

**Or embed config with `go:embed`:**
```go
//go:embed selectors.yaml
var configYAML string

config, _ := gtmlp.ParseConfig(configYAML, gtmlp.FormatYAML, nil)
products, _ := gtmlp.ScrapeURL[Product](context.Background(), "https://example.com", config)
```

## Logging

Configure log levels for different environments:

```go
import "log/slog"

// Development: see HTTP requests and scraping details
gtmlp.SetLogLevel(slog.LevelInfo)

// Troubleshooting: see XPath evaluation and fallbacks
gtmlp.SetLogLevel(slog.LevelDebug)

// Production: warnings and errors only (default)
gtmlp.SetLogLevel(slog.LevelWarn)

// Custom handler (JSON format, custom writer, etc.)
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
gtmlp.SetLogger(slog.New(handler))
```

**Log levels:**
- `Debug` - XPath evaluation, fallback usage, field extraction
- `Info` - HTTP requests, scraping progress, pagination
- `Warn` - Fallback usage, HTTP warnings, duplicate URLs (default)
- `Error` - HTTP failures, parsing errors, validation failures

## Security

GTMLP includes built-in SSRF (Server-Side Request Forgery) protection:

```go
config := &gtmlp.Config{
    Container: "//div[@class='product']",
    Fields:    fields,

    // SSRF protection (default: enabled)
    // Blocks: localhost, 127.0.0.1, 10.x.x.x, 192.168.x.x, 169.254.169.254
    AllowPrivateIPs: false, // set true to allow private IPs

    // Custom URL validator
    URLValidator: func(url string) error {
        if !strings.Contains(url, "example.com") {
            return errors.New("domain not allowed")
        }
        return nil
    },
}
```

See **[SECURITY.md](SECURITY.md)** for security best practices.

## Usage

```go
// Load config from file
config, _ := gtmlp.LoadConfig("selectors.yaml", nil)

// Or embed with go:embed
//go:embed selectors.yaml
var yaml string
config, _ := gtmlp.ParseConfig(yaml, gtmlp.FormatYAML, nil)

// Scrape
products, _ := gtmlp.ScrapeURL[Product](context.Background(), url, config)
results, _ := gtmlp.ScrapeURLUntyped(context.Background(), url, config) // returns []map[string]any
```

## Environment Variables

```bash
export GTMLP_TIMEOUT=30s
export GTMLP_USER_AGENT=MyBot/1.0
export GTMLP_PROXY=http://proxy:8080
```

## Fallback XPath Chains

Handle varying HTML structures with `altXpath` and `altContainer`:

```json
{
  "container": "//div[@class='product']",
  "altContainer": ["//article[@class='product']", "//div[@class='item']"],
  "fields": {
    "name": {
      "xpath": ".//h2/text()",
      "altXpath": [".//h3/text()", ".//h1/text()"]
    },
    "price": {
      "xpath": ".//span[@class='price']/text()",
      "altXpath": [".//div[@class='price']/text()"],
      "pipes": ["trim", "tofloat"]
    }
  }
}
```

**How it works:**
- Tries primary XPath first
- If empty (after pipes), tries each `altXpath` in order
- Returns first non-empty result
- Container fallback works the same way with `altContainer`

## Pagination

Auto-follow pagination or extract URLs for manual control:

**Next-Link Pagination** (follow "Next" buttons):
```json
{
  "container": "//div[@class='product']",
  "fields": {
    "name": {"xpath": ".//h2/text()"}
  },
  "pagination": {
    "type": "next-link",
    "nextSelector": "//a[@rel='next']/@href",
    "altSelectors": ["//a[contains(text(), 'Next')]/@href"],
    "maxPages": 50
  }
}
```

**Numbered Pagination** (extract all page links):
```json
{
  "pagination": {
    "type": "numbered",
    "pageSelector": "//div[@class='pagination']//a/@href",
    "maxPages": 20
  }
}
```

**Usage:**
```go
// Auto-follow: returns combined results from all pages
products, _ := gtmlp.ScrapeURL[Product](ctx, url, config)

// Page-separated: get results per page with metadata
results, _ := gtmlp.ScrapeURLWithPages[Product](ctx, url, config)

// Extract-only: get URLs for manual control
info, _ := gtmlp.ExtractPaginationURLs(ctx, url, config)
```

## Data Transformation Pipes

Transform extracted data using pipes:

```json
{
  "container": "//div[@class='product']",
  "fields": {
    "name": {"xpath": ".//h2/text()", "pipes": ["trim"]},
    "price": {"xpath": ".//span[@class='price']/text()", "pipes": ["trim", "tofloat"]},
    "url": {"xpath": ".//a/@href", "pipes": ["parseurl"]}
  }
}
```

**Built-in pipes:**
- `trim` - Remove whitespace
- `toint` - Convert to integer (strips `$`, `,`)
- `tofloat` - Convert to float (strips `$`, `,`)
- `parseurl` - Convert relative URLs to absolute
- `parsetime:layout:timezone` - Parse datetime
- `regexreplace:pattern:replacement:flags` - Regex substitution
- `humanduration` - Convert seconds to "X minutes ago"

**Custom pipes:**
```go
gtmlp.RegisterPipe("uppercase", func(ctx context.Context, input string, params []string) (any, error) {
    return strings.ToUpper(input), nil
})
```

See **[docs/API_V2.md](docs/API_V2.md)** for complete pipe documentation.

## Documentation & Examples

- **[API_V2.md](docs/API_V2.md)** - Complete API reference
- **[examples/v2/](examples/v2/)** - 10 working examples:
  - Basic scraping (JSON/YAML, embed)
  - E-commerce and tables
  - **Pagination** (next-link, numbered)

## License

MIT
