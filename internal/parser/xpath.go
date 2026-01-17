package parser

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

// Text returns the text content of the selected node.
func (s *Selection) Text() string {
	if s.node == nil {
		return ""
	}
	return htmlquery.InnerText(s.node)
}

// TextTrimmed returns the trimmed text content of the selected node.
func (s *Selection) TextTrimmed() string {
	return strings.TrimSpace(s.Text())
}

// HTML returns the outer HTML of the selected node.
func (s *Selection) HTML() string {
	if s.node == nil {
		return ""
	}
	return htmlquery.OutputHTML(s.node, true)
}

// InnerHTML returns the inner HTML of the selected node (excluding the outer tag).
func (s *Selection) InnerHTML() string {
	if s.node == nil {
		return ""
	}
	var sb strings.Builder
	for child := s.node.FirstChild; child != nil; child = child.NextSibling {
		html.Render(&sb, child)
	}
	return sb.String()
}

// Attr returns the value of an attribute on the selected node.
func (s *Selection) Attr(name string) string {
	if s.node == nil {
		return ""
	}
	return htmlquery.SelectAttr(s.node, name)
}

// AttrOr returns the value of an attribute, or a default value if not found.
func (s *Selection) AttrOr(name, defaultValue string) string {
	if val := s.Attr(name); val != "" {
		return val
	}
	return defaultValue
}

// Find executes an XPath expression relative to the current selection.
// Returns nil if no match is found.
func (s *Selection) Find(expr string) (*Selection, error) {
	if s.node == nil {
		return nil, fmt.Errorf("selection node is nil")
	}
	if expr == "" {
		return nil, fmt.Errorf("XPath expression cannot be empty")
	}

	node, err := htmlquery.Query(s.node, expr)
	if err != nil {
		return nil, fmt.Errorf("XPath query error: %w", err)
	}

	if node == nil {
		return nil, nil
	}

	return &Selection{node: node}, nil
}

// FindAll executes an XPath expression relative to the current selection.
func (s *Selection) FindAll(expr string) ([]*Selection, error) {
	if s.node == nil {
		return nil, fmt.Errorf("selection node is nil")
	}
	if expr == "" {
		return nil, fmt.Errorf("XPath expression cannot be empty")
	}

	nodes, err := htmlquery.QueryAll(s.node, expr)
	if err != nil {
		return nil, fmt.Errorf("XPath query error: %w", err)
	}

	selections := make([]*Selection, len(nodes))
	for i, node := range nodes {
		selections[i] = &Selection{node: node}
	}

	return selections, nil
}

// Each iterates over all child element nodes and calls the given function.
func (s *Selection) Each(fn func(int, *Selection)) {
	if s.node == nil {
		return
	}

	for i, child := range htmlquery.Find(s.node, "./*") {
		fn(i, &Selection{node: child})
	}
}

// Parent returns the parent node of the current selection.
func (s *Selection) Parent() *Selection {
	if s.node == nil || s.node.Parent == nil {
		return nil
	}
	return &Selection{node: s.node.Parent}
}

// Children returns all direct children of the current selection.
func (s *Selection) Children() []*Selection {
	if s.node == nil {
		return nil
	}

	var children []*Selection
	for child := s.node.FirstChild; child != nil; child = child.NextSibling {
		children = append(children, &Selection{node: child})
	}
	return children
}

// FirstChild returns the first child element (not text node).
func (s *Selection) FirstChild() *Selection {
	return s.findFirstChildElement()
}

// LastChild returns the last child element (not text node).
func (s *Selection) LastChild() *Selection {
	if s.node == nil {
		return nil
	}

	var lastElement *html.Node
	for child := s.node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			lastElement = child
		}
	}

	if lastElement != nil {
		return &Selection{node: lastElement}
	}
	return nil
}

// NextSibling returns the next sibling element.
func (s *Selection) NextSibling() *Selection {
	return s.findSibling(true)
}

// PrevSibling returns the previous sibling element.
func (s *Selection) PrevSibling() *Selection {
	return s.findSibling(false)
}

// EvaluateXPath evaluates a raw XPath expression and returns the result.
func (s *Selection) EvaluateXPath(expr string) (any, error) {
	if s.node == nil {
		return nil, fmt.Errorf("selection node is nil")
	}

	exprNode, err := xpath.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to compile XPath expression: %w", err)
	}

	result := exprNode.Evaluate(htmlquery.CreateXPathNavigator(s.node))
	return result, nil
}

// findFirstChildElement finds the first child element (not text node).
func (s *Selection) findFirstChildElement() *Selection {
	if s.node == nil {
		return nil
	}

	for child := s.node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			return &Selection{node: child}
		}
	}
	return nil
}

// findSibling finds a sibling element (next or previous).
func (s *Selection) findSibling(next bool) *Selection {
	if s.node == nil {
		return nil
	}

	var sibling *html.Node
	if next {
		sibling = s.node.NextSibling
	} else {
		sibling = s.node.PrevSibling
	}

	for sibling != nil {
		if sibling.Type == html.ElementNode {
			return &Selection{node: sibling}
		}
		if next {
			sibling = sibling.NextSibling
		} else {
			sibling = sibling.PrevSibling
		}
	}
	return nil
}

// Content returns the content based on the specified return type.
func (s *Selection) Content(returnType ReturnType) string {
	if s.node == nil {
		return ""
	}

	switch returnType {
	case ReturnTypeHTML:
		return s.HTML()
	default:
		return s.Text()
	}
}
