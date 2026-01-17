package gtmlp

import (
	"time"

	"github.com/Hanivan/gtmlp/internal/builder"
	"github.com/Hanivan/gtmlp/internal/health"
	"github.com/Hanivan/gtmlp/internal/parser"
	"github.com/Hanivan/gtmlp/internal/pattern"
	"github.com/Hanivan/gtmlp/internal/pipe"
	"github.com/Hanivan/gtmlp/internal/validator"
)

// Parser is an alias for internal parser.Parser.
type Parser = parser.Parser

// Selection is an alias for internal parser.Selection.
type Selection = parser.Selection

// JSONOptions is an alias for internal parser.JSONOptions.
type JSONOptions = parser.JSONOptions

// ReturnType is an alias for internal parser.ReturnType.
type ReturnType = parser.ReturnType

// PatternField is an alias for internal pattern.PatternField.
type PatternField = pattern.PatternField

// PatternMeta is an alias for internal pattern.PatternMeta.
type PatternMeta = pattern.PatternMeta

// MultipleType is an alias for internal pattern.MultipleType.
type MultipleType = pattern.MultipleType

// Extractor is an alias for internal pattern.Extractor.
type Extractor = pattern.Extractor

// Pipe is an alias for internal pipe.Pipe.
type Pipe = pipe.Pipe

// ValidationResult is an alias for internal validator.ValidationResult.
type ValidationResult = validator.ValidationResult

// HealthResult is an alias for internal health.HealthResult.
type HealthResult = health.HealthResult

const (
	// ReturnTypeText returns plain text content.
	ReturnTypeText = parser.ReturnTypeText
	// ReturnTypeHTML returns HTML content.
	ReturnTypeHTML = parser.ReturnTypeHTML
)

// Multiple type constants.
const (
	// MultipleNone returns only the first match.
	MultipleNone = pattern.MultipleNone
	// MultipleArray returns all matches as an array.
	MultipleArray = pattern.MultipleArray
	// MultipleSpace returns all matches joined with spaces.
	MultipleSpace = pattern.MultipleSpace
	// MultipleComma returns all matches joined with commas.
	MultipleComma = pattern.MultipleComma
)

// Parse creates a new Parser from an HTML string.
func Parse(html string) (*Parser, error) {
	return parser.New(html)
}

// ParseURL fetches HTML from a URL and creates a Parser.
func ParseURL(url string, opts ...Option) (*Parser, error) {
	cfg, client := applyOptions(opts...)

	html, err := client.GetHTML(url)
	if err != nil {
		return nil, NewParseError("failed to fetch URL", err)
	}

	p, err := parser.New(html)
	if err != nil {
		return nil, err
	}

	// Apply suppressErrors if configured
	if cfg.suppressErrors {
		p = p.WithSuppressErrors()
	}

	return p, nil
}

// ToJSON converts HTML string to JSON.
func ToJSON(html string) ([]byte, error) {
	p, err := Parse(html)
	if err != nil {
		return nil, err
	}
	return p.ToJSON()
}

// ToJSONWithOptions converts HTML string to JSON with custom options.
func ToJSONWithOptions(html string, opts JSONOptions) ([]byte, error) {
	p, err := Parse(html)
	if err != nil {
		return nil, err
	}
	return p.ToJSONWithOptions(opts)
}

// URLToJSON fetches HTML from a URL and converts it to JSON.
func URLToJSON(url string, parseOpts []Option, jsonOpts JSONOptions) ([]byte, error) {
	p, err := ParseURL(url, parseOpts...)
	if err != nil {
		return nil, err
	}
	return p.ToJSONWithOptions(jsonOpts)
}

// New creates a new Builder for fluent API usage.
func New() *builder.Builder {
	return builder.New()
}

// FromHTML creates a Builder from HTML content.
func FromHTML(html string) *builder.Builder {
	return builder.New().FromHTML(html)
}

// FromURL creates a Builder from a URL.
func FromURL(url string) *builder.Builder {
	return builder.New().FromURL(url)
}

// DefaultJSONOptions returns the default JSON conversion options.
func DefaultJSONOptions() JSONOptions {
	return parser.DefaultJSONOptions()
}

// Pattern helper functions

// NewPatternField creates a new PatternField.
func NewPatternField(key string, xpath string) PatternField {
	return pattern.NewPatternField(key, xpath)
}

// NewPatternFieldWithMultiple creates a new PatternField with multiple value handling.
func NewPatternFieldWithMultiple(key string, xpath string, multiple MultipleType) PatternField {
	return pattern.NewPatternFieldWithMultiple(key, xpath, multiple)
}

// NewPatternFieldWithHTML creates a new PatternField that returns HTML content.
func NewPatternFieldWithHTML(key string, xpath string) PatternField {
	return pattern.NewPatternFieldWithHTML(key, xpath)
}

// NewContainerPattern creates a new container PatternField.
func NewContainerPattern(key string, xpath string) PatternField {
	return pattern.NewContainerPattern(key, xpath)
}

// DefaultPatternMeta returns default pattern metadata.
func DefaultPatternMeta() *PatternMeta {
	return pattern.DefaultPatternMeta()
}

// NewExtractor creates a new Extractor from a Parser.
func NewExtractor(p *Parser) *Extractor {
	return pattern.NewExtractor(p)
}

// ExtractWithPatterns extracts data using pattern fields.
func ExtractWithPatterns(p *Parser, patterns []PatternField) ([]map[string]any, error) {
	extractor := pattern.NewExtractor(p)
	return extractor.ExtractWithPatterns(patterns)
}

// ExtractSingle extracts a single field using a pattern.
func ExtractSingle(p *Parser, field PatternField) (any, error) {
	extractor := pattern.NewExtractor(p)
	return extractor.ExtractSingle(field)
}

// Pipe helper functions

// NewTrimPipe creates a new TrimPipe.
func NewTrimPipe() Pipe {
	return pipe.NewTrimPipe()
}

// NewLowerCasePipe creates a new LowerCasePipe.
func NewLowerCasePipe() Pipe {
	return pipe.NewLowerCasePipe()
}

// NewUpperCasePipe creates a new UpperCasePipe.
func NewUpperCasePipe() Pipe {
	return pipe.NewUpperCasePipe()
}

// NewDecodePipe creates a new DecodePipe.
func NewDecodePipe() Pipe {
	return pipe.NewDecodePipe()
}

// NewReplacePipe creates a new ReplacePipe.
func NewReplacePipe(pattern, with string) Pipe {
	return pipe.NewReplacePipe(pattern, with)
}

// NewNumberNormalizePipe creates a new NumberNormalizePipe.
func NewNumberNormalizePipe() Pipe {
	return pipe.NewNumberNormalizePipe()
}

// NewURLResolvePipe creates a new URLResolvePipe.
func NewURLResolvePipe(baseURL string) Pipe {
	return pipe.NewURLResolvePipe(baseURL)
}

// NewExtractEmailPipe creates a new ExtractEmailPipe.
func NewExtractEmailPipe() Pipe {
	return pipe.NewExtractEmailPipe()
}

// NewDateFormatPipe creates a new DateFormatPipe.
func NewDateFormatPipe(format string) Pipe {
	return pipe.NewDateFormatPipe(format)
}

// RegisterPipe registers a custom pipe factory.
func RegisterPipe(name string, factory func() Pipe) {
	pipe.RegisterPipe(name, factory)
}

// CreatePipe creates a pipe by name.
func CreatePipe(name string) (Pipe, error) {
	return pipe.CreatePipe(name)
}

// ListPipes returns all registered pipe names.
func ListPipes() []string {
	return pipe.ListPipes()
}

// Validation functions

// ValidateXPath validates XPath expressions against HTML.
func ValidateXPath(html string, xpaths []string, suppressErrors bool) []ValidationResult {
	return validator.ValidateXPath(html, xpaths, suppressErrors)
}

// ValidateXPathWithParser validates XPath using an existing Parser.
func ValidateXPathWithParser(p *Parser, xpaths []string) []ValidationResult {
	return validator.ValidateXPathWithParser(p, xpaths)
}

// Health check functions

// CheckURLHealth checks the health of URLs concurrently.
func CheckURLHealth(urls []string, timeout time.Duration) []HealthResult {
	return health.CheckURLHealth(urls, timeout)
}

// CheckURLHealthSequential checks URLs sequentially.
func CheckURLHealthSequential(urls []string, timeout time.Duration) []HealthResult {
	return health.CheckURLHealthSequential(urls, timeout)
}

// CheckURLHealthWithGet checks URL health using GET requests.
func CheckURLHealthWithGet(urls []string, timeout time.Duration) []HealthResult {
	return health.CheckURLHealthWithGet(urls, timeout)
}
