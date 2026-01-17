package pattern

import (
	"testing"

	"github.com/Hanivan/gtmlp/internal/parser"
)

func TestNewPatternField(t *testing.T) {
	field := NewPatternField("title", "//h1")

	if field.Key != "title" {
		t.Errorf("Key mismatch, got: %s", field.Key)
	}

	if len(field.Patterns) != 1 {
		t.Errorf("Patterns length mismatch, got: %d", len(field.Patterns))
	}

	if field.Patterns[0] != "//h1" {
		t.Errorf("Pattern mismatch, got: %s", field.Patterns[0])
	}

	if field.ReturnType != parser.ReturnTypeText {
		t.Errorf("ReturnType should be Text, got: %s", field.ReturnType)
	}
}

func TestNewPatternFieldWithMultiple(t *testing.T) {
	field := NewPatternFieldWithMultiple("items", "//li", MultipleArray)

	if field.Meta.Multiple != MultipleArray {
		t.Errorf("Multiple type mismatch, got: %s", field.Meta.Multiple)
	}
}

func TestNewPatternFieldWithHTML(t *testing.T) {
	field := NewPatternFieldWithHTML("content", "//div")

	if field.ReturnType != parser.ReturnTypeHTML {
		t.Errorf("ReturnType should be HTML, got: %s", field.ReturnType)
	}
}

func TestNewContainerPattern(t *testing.T) {
	field := NewContainerPattern("items", "//div[@class='item']")

	if !field.Meta.IsContainer {
		t.Error("IsContainer should be true")
	}

	if field.Meta.ContainerKey != "items" {
		t.Errorf("ContainerKey mismatch, got: %s", field.Meta.ContainerKey)
	}
}

func TestExtractSingleField(t *testing.T) {
	html := `<html><body><h1>Hello World</h1></body></html>`
	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)
	field := NewPatternField("title", "//h1")

	result, err := extractor.ExtractSingle(field)
	if err != nil {
		t.Fatalf("ExtractSingle() failed: %v", err)
	}

	if result == nil {
		t.Fatal("ExtractSingle() returned nil")
	}

	title, ok := result.(string)
	if !ok {
		t.Fatal("Result is not a string")
	}

	if title != "Hello World" {
		t.Errorf("Title mismatch, got: %s", title)
	}
}

func TestExtractWithPatterns(t *testing.T) {
	html := `<html><body>
		<h1>Title</h1>
		<p class="desc">Description</p>
	</body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)
	patterns := []PatternField{
		NewPatternField("title", "//h1"),
		NewPatternField("description", "//p[@class='desc']"),
	}

	results, err := extractor.ExtractWithPatterns(patterns)
	if err != nil {
		t.Fatalf("ExtractWithPatterns() failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got: %d", len(results))
	}

	result := results[0]

	if result["title"] != "Title" {
		t.Errorf("Title mismatch, got: %v", result["title"])
	}

	if result["description"] != "Description" {
		t.Errorf("Description mismatch, got: %v", result["description"])
	}
}

func TestExtractMultipleArray(t *testing.T) {
	html := `<html><body>
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
			<li>Item 3</li>
		</ul>
	</body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)
	field := NewPatternFieldWithMultiple("items", "//li", MultipleArray)

	result, err := extractor.ExtractSingle(field)
	if err != nil {
		t.Fatalf("ExtractSingle() failed: %v", err)
	}

	items, ok := result.([]string)
	if !ok {
		t.Fatal("Result is not a string array")
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got: %d", len(items))
	}

	if items[0] != "Item 1" {
		t.Errorf("First item mismatch, got: %s", items[0])
	}
}

func TestExtractMultipleSpace(t *testing.T) {
	html := `<html><body>
		<ul>
			<li>Apple</li>
			<li>Banana</li>
		</ul>
	</body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)
	field := NewPatternFieldWithMultiple("fruits", "//li", MultipleSpace)

	result, err := extractor.ExtractSingle(field)
	if err != nil {
		t.Fatalf("ExtractSingle() failed: %v", err)
	}

	fruits, ok := result.(string)
	if !ok {
		t.Fatal("Result is not a string")
	}

	if fruits != "Apple Banana" {
		t.Errorf("Fruits mismatch, got: %s", fruits)
	}
}

func TestExtractMultipleComma(t *testing.T) {
	html := `<html><body>
		<span>Red</span>
		<span>Green</span>
		<span>Blue</span>
	</body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)
	field := NewPatternFieldWithMultiple("colors", "//span", MultipleComma)

	result, err := extractor.ExtractSingle(field)
	if err != nil {
		t.Fatalf("ExtractSingle() failed: %v", err)
	}

	colors, ok := result.(string)
	if !ok {
		t.Fatal("Result is not a string")
	}

	if colors != "Red, Green, Blue" {
		t.Errorf("Colors mismatch, got: %s", colors)
	}
}

func TestExtractWithContainer(t *testing.T) {
	html := `<html><body>
		<div class="product">
			<h2>Product 1</h2>
			<span class="price">$10</span>
		</div>
		<div class="product">
			<h2>Product 2</h2>
			<span class="price">$20</span>
		</div>
	</body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)
	patterns := []PatternField{
		NewContainerPattern("products", "//div[@class='product']"),
		NewPatternField("name", ".//h2"),
		NewPatternField("price", ".//span[@class='price']"),
	}

	results, err := extractor.ExtractWithPatterns(patterns)
	if err != nil {
		t.Fatalf("ExtractWithPatterns() failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got: %d", len(results))
	}

	if results[0]["name"] != "Product 1" {
		t.Errorf("First product name mismatch, got: %v", results[0]["name"])
	}

	if results[0]["price"] != "$10" {
		t.Errorf("First product price mismatch, got: %v", results[0]["price"])
	}

	if results[1]["name"] != "Product 2" {
		t.Errorf("Second product name mismatch, got: %v", results[1]["name"])
	}
}

func TestExtractWithAlternativePatterns(t *testing.T) {
	html := `<html><body><h3>Fallback</h3></body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)

	field := PatternField{
		Key:        "title",
		Patterns:   []string{"//h1"}, // This won't match
		AlterPattern: []string{"//h3"}, // This will match
		ReturnType: parser.ReturnTypeText,
		Meta:       DefaultPatternMeta(),
	}

	result, err := extractor.ExtractSingle(field)
	if err != nil {
		t.Fatalf("ExtractSingle() failed: %v", err)
	}

	if result == nil {
		t.Fatal("ExtractSingle() returned nil")
	}

	title, ok := result.(string)
	if !ok {
		t.Fatal("Result is not a string")
	}

	if title != "Fallback" {
		t.Errorf("Title mismatch, got: %s", title)
	}
}

func TestExtractHTMLContent(t *testing.T) {
	html := `<html><body><p class="test">Hello <strong>World</strong></p></body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)
	field := NewPatternFieldWithHTML("content", "//p")

	result, err := extractor.ExtractSingle(field)
	if err != nil {
		t.Fatalf("ExtractSingle() failed: %v", err)
	}

	if result == nil {
		t.Fatal("ExtractSingle() returned nil")
	}

	content, ok := result.(string)
	if !ok {
		t.Fatal("Result is not a string")
	}

	if content == "" {
		t.Error("HTML content is empty")
	}

	// Check that it contains the strong tag
	if !contains(content, "<strong>") {
		t.Error("HTML content should contain <strong> tag")
	}
}

func TestExtractEmpty(t *testing.T) {
	html := `<html><body><p>Content</p></body></html>`

	p, err := parser.New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	extractor := NewExtractor(p)
	field := NewPatternField("missing", "//div[@class='nonexistent']")

	result, err := extractor.ExtractSingle(field)
	if err != nil {
		t.Fatalf("ExtractSingle() failed: %v", err)
	}

	if result != nil {
		t.Error("ExtractSingle() should return nil for non-existent field")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:1] == substr[0:1] && contains(s[1:], substr[1:])) ||
			contains(s[1:], substr)))
}
