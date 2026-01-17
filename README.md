# GTMLP - Go HTML Parsing Library

A powerful and easy-to-use Go library for parsing HTML documents with XPath selectors and pattern-based extraction.

## Why GTMLP?

- **Clean Chainable API** - Fluent syntax for all operations including pattern-based extraction
- **Pattern-Based Extraction** - Define what to extract once, apply everywhere
- **Data Transformation Pipes** - Built-in pipes for cleaning, normalizing, and transforming data
- **Production Ready** - Random user-agents, retries, timeouts, error handling built-in
- **Zero External Dependencies** - Pure Go implementation with minimal dependencies

## Features

- **XPath Selectors** - Query HTML documents using XPath 1.0 expressions
- **Pattern-Based Extraction** - Container patterns for structured data extraction
- **Data Transformation Pipes** - Built-in pipes (trim, decode, replace, case conversion, etc.)
- **Chainable Fluent API** - Clean, composable syntax for all operations
- **HTML to JSON** - Convert HTML documents to structured JSON format
- **HTTP Client** - Built-in client for fetching and parsing remote URLs
- **Random User-Agents** - Realistic rotating browser user-agents for anti-detection
- **XPath Validation** - Validate XPath patterns before scraping
- **Alternative Patterns** - Fallback patterns for robust extraction
- **URL Health Check** - Check URL availability concurrently
- **XML Parsing** - Parse XML documents (sitemaps, RSS feeds, SOAP responses)

## Installation

```bash
go get github.com/Hanivan/gtmlp
```

## Quick Start

### Basic HTML Parsing

```go
package main

import (
    "fmt"
    "log"
    "github.com/Hanivan/gtmlp"
)

func main() {
    html := `<html><body><h1>Hello, World!</h1></body></html>`

    parser, err := gtmlp.Parse(html)
    if err != nil {
        log.Fatal(err)
    }

    h1, _ := parser.XPath("//h1")
    fmt.Println(h1.Text()) // Output: Hello, World!
}
```

### Pattern-Based Scraping (Recommended)

```go
package main

import (
    "fmt"
    "log"
    "github.com/Hanivan/gtmlp"
)

func main() {
    // Define what to extract
    patterns := []gtmlp.PatternField{
        gtmlp.NewContainerPattern("container", "//div[@class='product']"),
        {
            Key:      "name",
            Patterns: []string{".//h2/text()"},
            Pipes:    []gtmlp.Pipe{gtmlp.NewTrimPipe()},
        },
        {
            Key:      "price",
            Patterns: []string{".//span[@class='price']/text()"},
            Pipes:    []gtmlp.Pipe{gtmlp.NewTrimPipe()},
        },
    }

    // Extract with chainable API
    results, err := gtmlp.FromURL("https://example.com/products").
        WithPatterns(patterns).
        Extract()

    if err != nil {
        log.Fatal(err)
    }

    // Use the extracted data
    for _, product := range results {
        fmt.Printf("%s: %s\n", product["name"], product["price"])
    }
}
```

### Fluent API

```go
// Fetch and parse URL
parser, err := gtmlp.ParseURL("https://example.com",
    gtmlp.WithTimeout(10*time.Second),
    gtmlp.WithUserAgent("MyBot/1.0"),
)

// Chain operations
text, _ := gtmlp.FromHTML(html).
    XPath("//p[@class='content']").
    Text()

// Convert to JSON
json, _ := gtmlp.New().
    FromURL("https://example.com").
    ToJSON()
```

## Documentation

- **[API Reference](docs/API.md)** - Complete API documentation
- **[Pattern Extraction Guide](docs/PATTERNS.md)** - Pattern-based extraction and pipes
- **[Examples](docs/EXAMPLES.md)** - All examples and usage guides

## Running Examples

```bash
# Basic examples
go run examples/basic/*.go -type=all

# Advanced examples
go run examples/advanced/*.go -type=all

# Specific example
make example-chainable
```

See [docs/EXAMPLES.md](docs/EXAMPLES.md) for all available examples and commands.

## Anti-Detection Features

Random User-Agents are enabled by default. All HTTP requests use realistic, rotating user-agent strings from a comprehensive database of real browsers (Chrome, Firefox, Safari, Edge, etc.).

```go
// Random UA is enabled by default (recommended)
parser, err := gtmlp.ParseURL("https://example.com")

// Use custom user-agent instead
parser, err := gtmlp.ParseURL("https://example.com",
    gtmlp.WithUserAgent("CustomBot/1.0"),
)
```

**Benefits:**
- Avoid bot detection - Requests appear to come from real browsers
- Reduced blocking - Websites less likely to block your scraper
- Better success rate - More consistent scraping results
- No extra configuration - Works out of the box

## Project Structure

```
gtmlp/
├── internal/        # Internal packages
│   ├── parser/      # Core parsing functionality
│   ├── httpclient/  # HTTP client
│   ├── builder/     # Fluent API builder
│   ├── pipe/        # Data transformation pipes
│   ├── pattern/     # Pattern-based extraction
│   ├── validator/   # XPath validation
│   └── health/      # URL health checking
├── examples/
│   ├── basic/       # Basic examples (7 files)
│   └── advanced/    # Advanced examples (9 files)
├── docs/            # Documentation
├── types.go         # Public types and errors
├── options.go       # Configuration options
└── gtmlp.go         # Root package (public API)
```

## Dependencies

- `github.com/antchfx/htmlquery` - HTML parsing with XPath
- `github.com/antchfx/xpath` - XPath 1.0 implementation
- `github.com/lib4u/fake-useragent` - Realistic user-agent strings
- `golang.org/x/net` - Additional networking libraries

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Write tests for new features
4. Commit your changes
5. Push to the branch
6. Create a Pull Request

## Acknowledgments

Built with:
- [htmlquery](https://github.com/antchfx/htmlquery) - HTML parsing
- [fake-useragent](https://github.com/lib4u/fake-useragent) - User-agent database
