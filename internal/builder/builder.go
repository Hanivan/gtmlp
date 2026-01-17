package builder

import (
	"fmt"
	"time"

	"github.com/Hanivan/gtmlp/internal/httpclient"
	"github.com/Hanivan/gtmlp/internal/parser"
	"github.com/Hanivan/gtmlp/internal/pattern"
)

// Builder provides a fluent API for parsing HTML.
type Builder struct {
	html        string
	url         string
	clientOpts  []httpclient.ClientOption
	jsonOpts    parser.JSONOptions
	patterns    []pattern.PatternField
	loaded      bool
}

// New creates a new Builder instance.
func New() *Builder {
	return &Builder{
		clientOpts: []httpclient.ClientOption{},
		jsonOpts:   parser.DefaultJSONOptions(),
	}
}

// FromHTML sets the HTML content to parse.
func (b *Builder) FromHTML(html string) *Builder {
	b.html = html
	b.url = ""
	b.loaded = true
	return b
}

// FromURL sets the URL to fetch and parse.
func (b *Builder) FromURL(url string) *Builder {
	b.url = url
	b.html = ""
	b.loaded = false
	return b
}

// WithTimeout sets the HTTP request timeout.
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.clientOpts = append(b.clientOpts, httpclient.WithTimeout(timeout))
	return b
}

// WithUserAgent sets the User-Agent header.
func (b *Builder) WithUserAgent(ua string) *Builder {
	b.clientOpts = append(b.clientOpts, httpclient.WithUserAgent(ua))
	return b
}

// WithHeaders sets custom HTTP headers.
func (b *Builder) WithHeaders(headers map[string]string) *Builder {
	b.clientOpts = append(b.clientOpts, httpclient.WithHeaders(headers))
	return b
}

// WithProxy sets a proxy URL.
func (b *Builder) WithProxy(proxyURL string) *Builder {
	b.clientOpts = append(b.clientOpts, httpclient.WithProxy(proxyURL))
	return b
}

// WithJSONOptions sets the JSON conversion options.
func (b *Builder) WithJSONOptions(opts parser.JSONOptions) *Builder {
	b.jsonOpts = opts
	return b
}

// Parse parses the HTML and returns a Parser instance.
func (b *Builder) Parse() (*parser.Parser, error) {
	if err := b.loadIfNeeded(); err != nil {
		return nil, err
	}

	if b.html == "" {
		return nil, fmt.Errorf("no HTML content to parse")
	}

	p, err := parser.New(b.html)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return p, nil
}

// XPath executes an XPath expression and returns the first match.
func (b *Builder) XPath(expr string) (*parser.Selection, error) {
	p, err := b.Parse()
	if err != nil {
		return nil, err
	}
	return p.XPath(expr)
}

// XPathAll executes an XPath expression and returns all matches.
func (b *Builder) XPathAll(expr string) ([]*parser.Selection, error) {
	p, err := b.Parse()
	if err != nil {
		return nil, err
	}
	return p.XPathAll(expr)
}

// ToJSON converts the HTML to JSON.
func (b *Builder) ToJSON() ([]byte, error) {
	p, err := b.Parse()
	if err != nil {
		return nil, err
	}
	return p.ToJSONWithOptions(b.jsonOpts)
}

// Text returns the text content of the first element matching the XPath expression.
func (b *Builder) Text(expr string) (string, error) {
	sel, err := b.XPath(expr)
	if err != nil {
		return "", err
	}
	if sel == nil {
		return "", nil
	}
	return sel.Text(), nil
}

// HTML returns the outer HTML of the first element matching the XPath expression.
func (b *Builder) HTML(expr string) (string, error) {
	sel, err := b.XPath(expr)
	if err != nil {
		return "", err
	}
	if sel == nil {
		return "", nil
	}
	return sel.HTML(), nil
}

// Attr returns the value of an attribute for the first element matching the XPath expression.
func (b *Builder) Attr(expr, name string) (string, error) {
	sel, err := b.XPath(expr)
	if err != nil {
		return "", err
	}
	if sel == nil {
		return "", nil
	}
	return sel.Attr(name), nil
}

// WithPatterns sets the pattern fields for pattern-based extraction.
func (b *Builder) WithPatterns(patterns []pattern.PatternField) *Builder {
	b.patterns = patterns
	return b
}

// Extract executes pattern-based extraction using the configured patterns.
// Returns a slice of maps containing the extracted data.
func (b *Builder) Extract() ([]map[string]any, error) {
	if len(b.patterns) == 0 {
		return nil, fmt.Errorf("no patterns configured; use WithPatterns() first")
	}

	p, err := b.Parse()
	if err != nil {
		return nil, err
	}

	extractor := pattern.NewExtractor(p)
	return extractor.ExtractWithPatterns(b.patterns)
}

// loadIfNeeded fetches HTML from URL if not already loaded.
func (b *Builder) loadIfNeeded() error {
	if !b.loaded && b.url != "" {
		client := httpclient.NewClient(b.clientOpts...)
		html, err := client.GetHTML(b.url)
		if err != nil {
			return fmt.Errorf("failed to fetch URL: %w", err)
		}
		b.html = html
		b.loaded = true
	}
	return nil
}
