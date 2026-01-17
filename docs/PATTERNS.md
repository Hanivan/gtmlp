# Pattern-Based Extraction Guide

Extract structured data using container patterns and field patterns with data transformation pipes.

## Traditional Approach

```go
// Define extraction patterns
patterns := []gtmlp.PatternField{
    // Container pattern - defines each product
    gtmlp.NewContainerPattern("container", "//li[contains(@class, 'product')]"),

    // Field patterns - extracted from each container
    {
        Key:        "name",
        Patterns:   []string{".//h2/text()"},
        ReturnType: gtmlp.ReturnTypeText,
        Meta:       gtmlp.DefaultPatternMeta(),
        Pipes: []gtmlp.Pipe{
            gtmlp.NewTrimPipe(),
            gtmlp.NewReplacePipe(`\s+`, " "), // Replace multiple spaces
        },
    },
    {
        Key:        "price",
        Patterns:   []string{".//span[@class='price']//text()"},
        ReturnType: gtmlp.ReturnTypeText,
        Meta:       gtmlp.DefaultPatternMeta(),
        Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
    },
}

// Execute extraction
parser, _ := gtmlp.ParseURL("https://example.com/products")
results, _ := gtmlp.ExtractWithPatterns(parser, patterns)

// Results is []map[string]any with extracted data
for _, product := range results {
    fmt.Printf("Name: %s, Price: %s\n", product["name"], product["price"])
}
```

## Chainable API (Recommended)

```go
// Same patterns as above, but with clean chainable syntax
results, err := gtmlp.FromURL("https://example.com/products").
    WithPatterns(patterns).
    Extract()

// Or from HTML string
results, err := gtmlp.FromHTML(htmlString).
    WithPatterns(patterns).
    Extract()

// With additional options
results, err := gtmlp.New().
    FromURL("https://example.com/products").
    WithTimeout(10 * time.Second).
    WithUserAgent("MyBot/1.0").
    WithPatterns(patterns).
    Extract()
```

## Built-in Pipes

Transform and clean extracted data with built-in pipes:

- `trim` - Remove leading/trailing whitespace
- `toLowerCase` / `toUpperCase` - Case transformations
- `decode` - Decode HTML entities
- `replace` - Find and replace with regex support
- `stripHTML` - Remove HTML tags
- `numNormalize` - Normalize numbers (e.g., "1.5K" → "1500")
- `extractEmail` / `validateEmail` - Email extraction and validation
- `validateURL` - URL validation
- Custom pipes can be created by implementing the `Pipe` interface

## Alternative Patterns

Use multiple patterns with automatic fallback for robust extraction:

```go
{
    Key: "title",
    Patterns: []string{
        "//meta[@property='og:title']/@content",  // Try Open Graph first
        "//h1/text()",                             // Fall back to h1
        "//title/text()",                          // Finally try title tag
    },
    ReturnType: gtmlp.ReturnTypeText,
    Meta:       gtmlp.DefaultPatternMeta(),
}
```

## Common Use Cases

### E-commerce Product Scraping
Use container patterns to extract product listings with names, prices, images, and ratings. Pipes handle data cleaning and normalization automatically.

### News Article Extraction
Use alternative patterns to extract article metadata (title, author, date) from various site structures. Fallback patterns ensure extraction works across different CMS platforms.

### Sitemap and RSS Parsing
Parse XML documents like sitemaps and RSS feeds using the same XPath syntax. Extract URLs, metadata, and change frequencies for crawlers.

### Link Validation
Check scraped URLs for availability before processing. Identify dead links, track response codes, and validate URL health concurrently.

### Data Quality Assurance
Validate XPath patterns before production deployment. Test patterns against sample HTML to ensure they match expected elements.

### Price Monitoring
Track product prices over time with pattern-based extraction. Use pipes to normalize price formats and extract numeric values for comparison.

### Content Aggregation
Scrape content from multiple sources with different HTML structures. Alternative patterns provide robust extraction across varied layouts.
