package pattern

import (
	"fmt"
	"strings"

	"github.com/Hanivan/gtmlp/internal/parser"
	"github.com/Hanivan/gtmlp/internal/pipe"
)

// Extractor handles pattern-based extraction from HTML.
type Extractor struct {
	parser *parser.Parser
}

// NewExtractor creates a new Extractor from a Parser.
func NewExtractor(p *parser.Parser) *Extractor {
	return &Extractor{parser: p}
}

// ExtractWithPatterns extracts data using a list of pattern fields.
// Returns a slice of maps, where each map represents one extracted item.
// If a container pattern is present, returns multiple items.
// Otherwise, returns a single-item slice with one map.
func (e *Extractor) ExtractWithPatterns(patterns []PatternField) ([]map[string]any, error) {
	if len(patterns) == 0 {
		return []map[string]any{}, nil
	}

	// Check if there's a container pattern
	var container *PatternField
	var fields []PatternField

	for i, p := range patterns {
		if p.Meta.IsContainer {
			container = &patterns[i]
		} else {
			fields = append(fields, p)
		}
	}

	// If no container, extract single item
	if container == nil {
		return e.extractWithoutContainers(fields)
	}

	// Extract with containers
	return e.extractWithContainers(*container, fields)
}

// extractWithoutContainers extracts fields without a container pattern.
func (e *Extractor) extractWithoutContainers(fields []PatternField) ([]map[string]any, error) {
	result := make(map[string]any)

	for _, field := range fields {
		value, err := e.extractField(field, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to extract field '%s': %w", field.Key, err)
		}

		if value != nil {
			result[field.Key] = value
		}
	}

	return []map[string]any{result}, nil
}

// extractWithContainers extracts fields using a container pattern.
func (e *Extractor) extractWithContainers(container PatternField, fields []PatternField) ([]map[string]any, error) {
	// Find container selections
	containerSels, err := e.findSelectionsByPattern(container)
	if err != nil {
		return nil, fmt.Errorf("failed to find containers: %w", err)
	}

	if len(containerSels) == 0 {
		return []map[string]any{}, nil
	}

	results := make([]map[string]any, 0, len(containerSels))

	// For each container, extract fields
	for _, containerSel := range containerSels {
		item := make(map[string]any)

		// Add container key if specified
		if container.Meta.ContainerKey != "" {
			item[container.Meta.ContainerKey] = containerSel.Content(container.ReturnType)
		}

		// Extract fields within this container
		for _, field := range fields {
			value, err := e.extractFieldWithin(field, containerSel)
			if err != nil {
				continue // Skip fields that fail to extract
			}

			if value != nil {
				item[field.Key] = value
			}
		}

		// Only add non-empty items
		if len(item) > 0 {
			results = append(results, item)
		}
	}

	return results, nil
}

// ExtractSingle extracts a single field value.
func (e *Extractor) ExtractSingle(field PatternField) (any, error) {
	sels, err := e.findSelectionsByPattern(field)
	if err != nil {
		return nil, err
	}

	if len(sels) == 0 {
		return nil, nil
	}

	return e.processSelections(sels, field), nil
}

// extractField extracts a single field from root.
func (e *Extractor) extractField(field PatternField, _ any) (any, error) {
	sels, err := e.findSelectionsByPattern(field)
	if err != nil {
		return nil, err
	}

	if len(sels) == 0 {
		return nil, nil
	}

	return e.processSelections(sels, field), nil
}

// patternQueryFunc defines a function that queries patterns and returns selections.
type patternQueryFunc func(pattern string) ([]*parser.Selection, error)

// tryPatterns attempts patterns in order, returning the first successful match.
func (e *Extractor) tryPatterns(patterns []string, queryFn patternQueryFunc) []*parser.Selection {
	for _, pattern := range patterns {
		sels, err := queryFn(pattern)
		if err != nil {
			continue
		}
		if len(sels) > 0 {
			return sels
		}
	}
	return nil
}

// extractFieldWithin extracts a single field within a context selection.
func (e *Extractor) extractFieldWithin(field PatternField, contextSel *parser.Selection) (any, error) {
	queryFn := func(pattern string) ([]*parser.Selection, error) {
		return contextSel.FindAll(pattern)
	}

	// Try primary patterns, then alternatives
	if sels := e.tryPatterns(field.Patterns, queryFn); len(sels) > 0 {
		return e.processSelections(sels, field), nil
	}
	if sels := e.tryPatterns(field.AlterPattern, queryFn); len(sels) > 0 {
		return e.processSelections(sels, field), nil
	}

	return nil, nil
}

// findSelectionsByPattern finds selections matching a pattern field.
func (e *Extractor) findSelectionsByPattern(field PatternField) ([]*parser.Selection, error) {
	queryFn := func(pattern string) ([]*parser.Selection, error) {
		return e.parser.XPathAll(pattern)
	}

	// Try primary patterns, then alternatives
	if sels := e.tryPatterns(field.Patterns, queryFn); len(sels) > 0 {
		return sels, nil
	}
	if sels := e.tryPatterns(field.AlterPattern, queryFn); len(sels) > 0 {
		return sels, nil
	}

	return []*parser.Selection{}, nil
}

// processSelections processes extracted selections according to field settings.
func (e *Extractor) processSelections(sels []*parser.Selection, field PatternField) any {
	if len(sels) == 0 {
		return nil
	}

	// Handle multiple values based on Meta.Multiple
	switch field.Meta.Multiple {
	case MultipleArray:
		values := e.extractMultiple(sels, field)
		return e.applyPipes(values, field.Pipes)
	case MultipleSpace:
		values := e.extractMultiple(sels, field)
		transformed := e.applyPipes(values, field.Pipes)
		if arr, ok := transformed.([]string); ok {
			return strings.Join(arr, " ")
		}
		return transformed
	case MultipleComma:
		values := e.extractMultiple(sels, field)
		transformed := e.applyPipes(values, field.Pipes)
		if arr, ok := transformed.([]string); ok {
			return strings.Join(arr, ", ")
		}
		return transformed
	default: // MultipleNone
		value := sels[0].Content(field.ReturnType)
		return e.applyPipe(value, field.Pipes)
	}
}

// extractMultiple extracts multiple values from selections.
func (e *Extractor) extractMultiple(sels []*parser.Selection, field PatternField) []string {
	values := make([]string, 0, len(sels))

	for _, sel := range sels {
		value := sel.Content(field.ReturnType)
		if value != "" {
			values = append(values, value)
		}
	}

	return values
}

// applyPipe applies pipes to a single string value.
func (e *Extractor) applyPipe(value string, pipes []pipe.Pipe) string {
	result := value
	for _, p := range pipes {
		result = p.Process(result)
	}
	return result
}

// applyPipes applies pipes to an array of string values.
func (e *Extractor) applyPipes(values []string, pipes []pipe.Pipe) any {
	if len(pipes) == 0 {
		return values
	}

	result := make([]string, len(values))
	for i, value := range values {
		result[i] = e.applyPipe(value, pipes)
	}
	return result
}
