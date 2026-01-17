# API Reference

## Package Functions

### Parsing

- `Parse(html string) (*Parser, error)` - Parse HTML string
- `ParseURL(url string, opts ...Option) (*Parser, error)` - Fetch and parse URL
- `ToJSON(html string) ([]byte, error)` - Convert HTML to JSON
- `ToJSONWithOptions(html string, opts JSONOptions) ([]byte, error)` - Convert with options

### Chainable API

- `New() *Builder` - Create new builder instance
- `FromHTML(html string) *Builder` - Create builder from HTML
- `FromURL(url string) *Builder` - Create builder from URL

## Parser Methods

- `XPath(expr string) (*Selection, error)` - Execute XPath, return first match
- `XPathAll(expr string) ([]*Selection, error)` - Execute XPath, return all matches
- `ToJSON() ([]byte, error)` - Convert to JSON
- `ToJSONWithOptions(opts JSONOptions) ([]byte, error)` - Convert with options

## Builder Methods

### Configuration

- `FromHTML(html string) *Builder` - Set HTML content to parse
- `FromURL(url string) *Builder` - Set URL to fetch and parse
- `WithTimeout(d time.Duration) *Builder` - Set HTTP timeout
- `WithUserAgent(ua string) *Builder` - Set User-Agent header
- `WithHeaders(h map[string]string) *Builder` - Set custom headers
- `WithProxy(proxyURL string) *Builder` - Set proxy URL
- `WithJSONOptions(opts JSONOptions) *Builder` - Set JSON conversion options
- `WithPatterns(patterns []PatternField) *Builder` - Set extraction patterns

### Execution

- `Parse() (*Parser, error)` - Parse HTML and return Parser
- `XPath(expr string) (*Selection, error)` - Execute XPath, return first match
- `XPathAll(expr string) ([]*Selection, error)` - Execute XPath, return all matches
- `ToJSON() ([]byte, error)` - Convert to JSON
- `Text(expr string) (string, error)` - Get text content of element
- `HTML(expr string) (string, error)` - Get HTML content of element
- `Attr(expr, name string) (string, error)` - Get attribute value
- `Extract() ([]map[string]any, error)` - Execute pattern-based extraction

## Selection Methods

- `Text() string` - Get text content
- `TextTrimmed() string` - Get trimmed text content
- `HTML() string` - Get outer HTML
- `InnerHTML() string` - Get inner HTML
- `Attr(name string) string` - Get attribute value
- `AttrOr(name, defaultValue string) string` - Get attribute or default
- `Find(expr string) (*Selection, error)` - Find child element
- `FindAll(expr string) ([]*Selection, error)` - Find all child elements
- `Each(fn func(int, *Selection))` - Iterate over children
- `Parent() *Selection` - Get parent element
- `Children() []*Selection` - Get all children
- `FirstChild() *Selection` - Get first child element
- `LastChild() *Selection` - Get last child element
- `NextSibling() *Selection` - Get next sibling
- `PrevSibling() *Selection` - Get previous sibling
- `ToJSON(opts JSONOptions) ([]byte, error)` - Convert to JSON
- `ToMap(opts JSONOptions) map[string]any` - Convert to map

## Options

- `WithTimeout(d time.Duration)` - Set HTTP timeout
- `WithUserAgent(ua string)` - Set User-Agent header
- `WithHeaders(h map[string]string)` - Set custom headers
- `WithProxy(proxyURL string)` - Set proxy URL
- `WithMaxRetries(n int)` - Set max retries for HTTP requests
- `WithSuppressErrors()` - Suppress XPath errors and handle gracefully

## Pattern-Based Extraction

- `ExtractWithPatterns(parser *Parser, patterns []PatternField) ([]map[string]any, error)` - Extract structured data using patterns
- `NewContainerPattern(key, pattern string) PatternField` - Create container pattern for grouped extraction
- `NewTrimPipe()` - Create trim pipe
- `NewReplacePipe(pattern, replacement string)` - Create replace pipe with regex support
- `NewToLowerCasePipe()` / `NewToUpperCasePipe()` - Create case transformation pipes
- `NewDecodePipe()` - Create HTML entity decode pipe
- `NewStripHTMLPipe()` - Create HTML tag stripping pipe
- `NewNumNormalizePipe()` - Create number normalization pipe
- `NewExtractEmailPipe()` / `NewValidateEmailPipe()` - Email extraction and validation pipes
- `NewValidateURLPipe()` - URL validation pipe

## XPath Validation

- `ValidateXPath(html string, patterns []string) []ValidationResult` - Validate multiple XPath patterns
- `ValidateXPathSingle(html string, pattern string) ValidationResult` - Validate single XPath pattern

## URL Health Check

- `CheckURLHealth(url string, opts ...Option) (*HealthCheckResult, error)` - Check single URL availability
- `CheckMultipleURLs(urls []string, opts ...Option) ([]*HealthCheckResult, error)` - Check multiple URLs concurrently

## JSON Conversion Options

```go
type JSONOptions struct {
    IncludeAttributes bool  // Include element attributes
    IncludeTextContent bool // Include text content
    PrettyPrint bool        // Pretty-print JSON output
    TrimWhitespace bool     // Trim whitespace in text
}
```
