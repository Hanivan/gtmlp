# GTMLP

Type-safe HTML scraping with XPath selectors and external configuration.

## Features

- **Type-safe** with Go generics
- **External config** (JSON/YAML)
- **XPath validation** before scraping
- **Fallback XPath chains** (`altXpath`, `altContainer`) for handling varying HTML structures
- **Data transformation pipes** (trim, int/float conversion, regex, URL parsing, etc.)
- **Custom pipe registration** for domain-specific transformations
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
- **[examples/v2/](examples/v2/)** - 8 working examples (JSON/YAML, embed, tables, ecommerce)

## License

MIT
