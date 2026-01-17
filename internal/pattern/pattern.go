package pattern

import (
	"github.com/Hanivan/gtmlp/internal/parser"
	"github.com/Hanivan/gtmlp/internal/pipe"
)

// MultipleType defines how to handle multiple matches.
type MultipleType string

const (
	// MultipleNone returns only the first match.
	MultipleNone MultipleType = ""
	// MultipleArray returns all matches as an array.
	MultipleArray MultipleType = "array"
	// MultipleSpace returns all matches joined with spaces.
	MultipleSpace MultipleType = "with space"
	// MultipleComma returns all matches joined with commas.
	MultipleComma MultipleType = "with comma"
)

// PatternMeta contains metadata about how to extract and process a pattern.
type PatternMeta struct {
	Multiple     MultipleType      // How to handle multiple matches
	Multiline    bool              // Whether to preserve multiline text
	IsContainer  bool              // Whether this pattern defines a container for other fields
	ContainerKey string            // Key to use for container-based extraction
}

// PatternField defines a field to extract from HTML using XPath patterns.
type PatternField struct {
	Key          string              // The key to use in the result map
	Patterns     []string            // Primary XPath patterns (tried in order)
	AlterPattern []string            // Alternative patterns (tried if primary fails)
	ReturnType   parser.ReturnType   // Whether to return text or HTML (default: text)
	Meta         *PatternMeta        // Metadata about extraction behavior
	Pipes        []pipe.Pipe         // Transformation pipes to apply
}

// DefaultPatternMeta returns default pattern metadata.
func DefaultPatternMeta() *PatternMeta {
	return &PatternMeta{
		Multiple:    MultipleNone,
		Multiline:   false,
		IsContainer: false,
	}
}

// NewPatternField creates a new PatternField with the given key and pattern.
func NewPatternField(key string, pattern string) PatternField {
	return PatternField{
		Key:        key,
		Patterns:   []string{pattern},
		ReturnType: parser.ReturnTypeText,
		Meta:       DefaultPatternMeta(),
	}
}

// NewPatternFieldWithMultiple creates a new PatternField that returns multiple values.
func NewPatternFieldWithMultiple(key string, pattern string, multiple MultipleType) PatternField {
	meta := DefaultPatternMeta()
	meta.Multiple = multiple
	return PatternField{
		Key:        key,
		Patterns:   []string{pattern},
		ReturnType: parser.ReturnTypeText,
		Meta:       meta,
	}
}

// NewPatternFieldWithHTML creates a new PatternField that returns HTML content.
func NewPatternFieldWithHTML(key string, pattern string) PatternField {
	return PatternField{
		Key:        key,
		Patterns:   []string{pattern},
		ReturnType: parser.ReturnTypeHTML,
		Meta:       DefaultPatternMeta(),
	}
}

// NewContainerPattern creates a new PatternField that acts as a container.
func NewContainerPattern(key string, pattern string) PatternField {
	meta := DefaultPatternMeta()
	meta.IsContainer = true
	meta.ContainerKey = key
	return PatternField{
		Key:        key,
		Patterns:   []string{pattern},
		ReturnType: parser.ReturnTypeText,
		Meta:       meta,
	}
}
