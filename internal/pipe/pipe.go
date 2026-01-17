package pipe

import (
	"html"
	"regexp"
	"strings"
)

// Pipe defines the interface for data transformation pipes.
type Pipe interface {
	Process(string) string
}

// TrimPipe removes leading and trailing whitespace.
type TrimPipe struct{}

// Process removes leading and trailing whitespace from the input.
func (p *TrimPipe) Process(s string) string {
	return strings.TrimSpace(s)
}

// LowerCasePipe converts text to lowercase.
type LowerCasePipe struct{}

// Process converts the input to lowercase.
func (p *LowerCasePipe) Process(s string) string {
	return strings.ToLower(s)
}

// UpperCasePipe converts text to uppercase.
type UpperCasePipe struct{}

// Process converts the input to uppercase.
func (p *UpperCasePipe) Process(s string) string {
	return strings.ToUpper(s)
}

// DecodePipe decodes HTML entities.
type DecodePipe struct{}

// Process decodes HTML entities in the input.
func (p *DecodePipe) Process(s string) string {
	return html.UnescapeString(s)
}

// ReplacePipe replaces text using a regular expression.
type ReplacePipe struct {
	Pattern string
	With    string
}

// Process applies regex replacement to the input.
func (p *ReplacePipe) Process(s string) string {
	re := regexp.MustCompile(p.Pattern)
	return re.ReplaceAllString(s, p.With)
}

// TrimLeftPipe removes leading whitespace.
type TrimLeftPipe struct{}

// Process removes leading whitespace from the input.
func (p *TrimLeftPipe) Process(s string) string {
	return strings.TrimLeft(s, " \t\n\r")
}

// TrimRightPipe removes trailing whitespace.
type TrimRightPipe struct{}

// Process removes trailing whitespace from the input.
func (p *TrimRightPipe) Process(s string) string {
	return strings.TrimRight(s, " \t\n\r")
}

// StripHTMLPipe removes all HTML tags.
type StripHTMLPipe struct{}

// Process removes HTML tags from the input.
func (p *StripHTMLPipe) Process(s string) string {
	// Simple regex-based HTML tag removal
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

// NewTrimPipe creates a new TrimPipe.
func NewTrimPipe() *TrimPipe {
	return &TrimPipe{}
}

// NewLowerCasePipe creates a new LowerCasePipe.
func NewLowerCasePipe() *LowerCasePipe {
	return &LowerCasePipe{}
}

// NewUpperCasePipe creates a new UpperCasePipe.
func NewUpperCasePipe() *UpperCasePipe {
	return &UpperCasePipe{}
}

// NewDecodePipe creates a new DecodePipe.
func NewDecodePipe() *DecodePipe {
	return &DecodePipe{}
}

// NewReplacePipe creates a new ReplacePipe.
func NewReplacePipe(pattern, with string) *ReplacePipe {
	return &ReplacePipe{Pattern: pattern, With: with}
}

// NewTrimLeftPipe creates a new TrimLeftPipe.
func NewTrimLeftPipe() *TrimLeftPipe {
	return &TrimLeftPipe{}
}

// NewTrimRightPipe creates a new TrimRightPipe.
func NewTrimRightPipe() *TrimRightPipe {
	return &TrimRightPipe{}
}

// NewStripHTMLPipe creates a new StripHTMLPipe.
func NewStripHTMLPipe() *StripHTMLPipe {
	return &StripHTMLPipe{}
}
