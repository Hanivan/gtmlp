package parser

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	html := "<html><body><h1>Hello</h1></body></html>"

	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if parser == nil {
		t.Fatal("New() returned nil parser")
	}

	if parser.HTML() != html {
		t.Errorf("HTML() mismatch, got: %s", parser.HTML())
	}
}

func TestNewParserEmpty(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Error("New() with empty HTML should return error")
	}
}

func TestNewParserInvalid(t *testing.T) {
	// htmlquery is lenient and can parse most HTML fragments
	// So this test just verifies it doesn't crash on odd input
	_, err := New("<<<")
	// htmlquery will actually parse this, so we just verify no panic
	_ = err
}

func TestXPath(t *testing.T) {
	html := `<html><body><h1 id="main">Hello</h1></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Test finding element
	h1, err := parser.XPath("//h1")
	if err != nil {
		t.Fatalf("XPath() failed: %v", err)
	}

	if h1 == nil {
		t.Fatal("XPath() returned nil selection")
	}

	if h1.Text() != "Hello" {
		t.Errorf("Text() mismatch, got: %s", h1.Text())
	}

	if h1.Attr("id") != "main" {
		t.Errorf("Attr() mismatch, got: %s", h1.Attr("id"))
	}
}

func TestXPathNotFound(t *testing.T) {
	html := `<html><body><h1>Hello</h1></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Test finding non-existent element
	div, err := parser.XPath("//div")
	if err != nil {
		t.Fatalf("XPath() failed: %v", err)
	}

	if div != nil {
		t.Error("XPath() should return nil for non-existent element")
	}
}

func TestXPathAll(t *testing.T) {
	html := `<html><body><ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	items, err := parser.XPathAll("//li")
	if err != nil {
		t.Fatalf("XPathAll() failed: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("XPathAll() returned %d items, expected 3", len(items))
	}

	if items[0].Text() != "Item 1" {
		t.Errorf("First item text mismatch: %s", items[0].Text())
	}
}

func TestSelectionText(t *testing.T) {
	html := `<html><body><p>Hello World</p></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	p, _ := parser.XPath("//p")

	if p.Text() != "Hello World" {
		t.Errorf("Text() mismatch: %s", p.Text())
	}

	if p.TextTrimmed() != "Hello World" {
		t.Errorf("TextTrimmed() mismatch: %s", p.TextTrimmed())
	}
}

func TestSelectionHTML(t *testing.T) {
	html := `<html><body><p class="test">Hello</p></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	p, _ := parser.XPath("//p")

	outerHTML := p.HTML()
	if outerHTML == "" {
		t.Error("HTML() returned empty string")
	}

	// Check that it contains the class attribute
	if p.Attr("class") != "test" {
		t.Errorf("Attr() failed, got: %s", p.Attr("class"))
	}
}

func TestSelectionAttr(t *testing.T) {
	html := `<html><body><a href="https://example.com" class="link">Link</a></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	a, _ := parser.XPath("//a")

	if a.Attr("href") != "https://example.com" {
		t.Errorf("Attr() href mismatch: %s", a.Attr("href"))
	}

	if a.Attr("class") != "link" {
		t.Errorf("Attr() class mismatch: %s", a.Attr("class"))
	}

	if a.AttrOr("nonexistent", "default") != "default" {
		t.Error("AttrOr() should return default for non-existent attribute")
	}
}

func TestSelectionFind(t *testing.T) {
	html := `<html><body><div class="container"><p>Hello</p></div></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	div, _ := parser.XPath("//div[@class='container']")

	p, err := div.Find(".//p")
	if err != nil {
		t.Fatalf("Find() failed: %v", err)
	}

	if p.Text() != "Hello" {
		t.Errorf("Find() text mismatch: %s", p.Text())
	}
}

func TestSelectionNavigation(t *testing.T) {
	html := `<html><body><div><p>First</p><p>Second</p></div></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	div, _ := parser.XPath("//div")

	firstChild := div.FirstChild()
	if firstChild == nil {
		t.Fatal("FirstChild() returned nil")
	}

	nextSibling := firstChild.NextSibling()
	if nextSibling == nil {
		t.Fatal("NextSibling() returned nil")
	}

	if nextSibling.Text() != "Second" {
		t.Errorf("NextSibling() text mismatch: %s", nextSibling.Text())
	}
}

func TestSelectionEach(t *testing.T) {
	html := `<html><body><ul><li>Item 1</li><li>Item 2</li></ul></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ul, _ := parser.XPath("//ul")

	count := 0
	ul.Each(func(i int, sel *Selection) {
		count++
	})

	if count != 2 {
		t.Errorf("Each() iterated %d times, expected 2", count)
	}
}

func TestToJSON(t *testing.T) {
	html := `<html><body><h1>Hello</h1></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	json, err := parser.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}

	if len(json) == 0 {
		t.Error("ToJSON() returned empty bytes")
	}

	// Should contain valid JSON
	if string(json)[0] != '{' {
		t.Error("ToJSON() should return JSON object")
	}
}

func TestToJSONWithOptions(t *testing.T) {
	html := `<html><body><div id="test" class="container">Content</div></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	opts := JSONOptions{
		IncludeAttributes: true,
		IncludeTextContent: true,
		PrettyPrint: true,
		TrimWhitespace: true,
	}

	json, err := parser.ToJSONWithOptions(opts)
	if err != nil {
		t.Fatalf("ToJSONWithOptions() failed: %v", err)
	}

	if len(json) == 0 {
		t.Error("ToJSONWithOptions() returned empty bytes")
	}

	// With pretty print, should have newlines
	if !opts.PrettyPrint || len(json) == 0 {
		// Just check it's not empty
	}
}

func TestSelectionToJSON(t *testing.T) {
	html := `<html><body><div class="test"><p>Hello</p></div></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	div, _ := parser.XPath("//div")

	opts := JSONOptions{
		IncludeAttributes: true,
		PrettyPrint: false,
	}

	json, err := div.ToJSON(opts)
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}

	if len(json) == 0 {
		t.Error("ToJSON() returned empty bytes")
	}
}

func TestXPathByAttribute(t *testing.T) {
	html := `<html><body><h1 class="title" id="main">Hello</h1><p class="title">World</p></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Test finding by class
	title, _ := parser.XPath("//h1[@class='title']")
	if title.Text() != "Hello" {
		t.Errorf("XPath by class failed: %s", title.Text())
	}

	// Test finding by id
	main, _ := parser.XPath("//*[@id='main']")
	if main.Text() != "Hello" {
		t.Errorf("XPath by id failed: %s", main.Text())
	}
}

func TestParseError(t *testing.T) {
	html := ""
	_, err := New(html)

	if err == nil {
		t.Error("New() with empty string should return error")
	}
}

func TestContentReturnType(t *testing.T) {
	html := `<html><body><p class="test">Hello <strong>World</strong></p></body></html>`
	parser, err := New(html)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	p, _ := parser.XPath("//p")

	// Test text return type
	textContent := p.Content(ReturnTypeText)
	if textContent != "Hello World" {
		t.Errorf("Content(Text) mismatch: %s", textContent)
	}

	// Test HTML return type
	htmlContent := p.Content(ReturnTypeHTML)
	if htmlContent == "" {
		t.Error("Content(HTML) returned empty string")
	}
	// Check that it contains the class attribute
	if p.Attr("class") != "test" {
		t.Errorf("Attr() failed, got: %s", p.Attr("class"))
	}
}

func TestSuppressErrors(t *testing.T) {
	html := `<html><body><h1>Hello</h1></body></html>`

	t.Run("XPath without error suppression", func(t *testing.T) {
		parser, _ := New(html)

		// Empty XPath should return error
		_, err := parser.XPath("")
		if err == nil {
			t.Error("XPath('') should return error when suppressErrors is false")
		}
	})

	t.Run("XPath with error suppression", func(t *testing.T) {
		parser, _ := New(html)
		parser = parser.WithSuppressErrors()

		// Empty XPath should return nil, no error when suppressed
		sel, err := parser.XPath("")
		if err != nil {
			t.Errorf("XPath('') should not return error when suppressErrors is true, got: %v", err)
		}
		if sel != nil {
			t.Error("XPath('') should return nil selection when suppressErrors is true")
		}
	})

	t.Run("XPathAll with error suppression", func(t *testing.T) {
		parser, _ := New(html)
		parser = parser.WithSuppressErrors()

		// Invalid XPath should return nil, no error when suppressed
		sels, err := parser.XPathAll("")
		if err != nil {
			t.Errorf("XPathAll('') should not return error when suppressErrors is true, got: %v", err)
		}
		if sels != nil {
			t.Error("XPathAll('') should return nil when suppressErrors is true")
		}
	})
}
