# GTMLP

Type-safe HTML scraping with XPath selectors and external configuration.

## Features

- **Type-safe** with Go generics
- **External config** (JSON/YAML)
- **XPath validation** before scraping
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
    "name": ".//h2/text()",
    "price": ".//span[@class='price']/text()"
  }
}
```

**main.go:**
```go
type Product struct {
    Name  string `json:"name"`
    Price string `json:"price"`
}

config, _ := gtmlp.LoadConfig("selectors.json", nil)
products, _ := gtmlp.ScrapeURL[Product]("https://example.com", config)

for _, p := range products {
    fmt.Printf("%s: %s\n", p.Name, p.Price)
}
```

**Or embed config with `go:embed`:**
```go
//go:embed selectors.yaml
var configYAML string

config, _ := gtmlp.ParseConfig(configYAML, gtmlp.FormatYAML, nil)
products, _ := gtmlp.ScrapeURL[Product]("https://example.com", config)
```

## Usage

```go
// Load config from file
config, _ := gtmlp.LoadConfig("selectors.yaml", nil)

// Or embed with go:embed
//go:embed selectors.yaml
var yaml string
config, _ := gtmlp.ParseConfig(yaml, gtmlp.FormatYAML, nil)

// Scrape
products, _ := gtmlp.ScrapeURL[Product](url, config)
results, _ := gtmlp.ScrapeURLUntyped(url, config) // returns []map[string]any
```

## Environment Variables

```bash
export GTMLP_TIMEOUT=30s
export GTMLP_USER_AGENT=MyBot/1.0
export GTMLP_PROXY=http://proxy:8080
```

## Documentation & Examples

- **[API_V2.md](docs/API_V2.md)** - Complete API reference
- **[examples/v2/](examples/v2/)** - 8 working examples (JSON/YAML, embed, tables, ecommerce)

## License

MIT
