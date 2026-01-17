package parser

import (
	"encoding/json"
	"fmt"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// JSONOptions controls the behavior of HTML to JSON conversion.
type JSONOptions struct {
	IncludeAttributes bool
	IncludeTextContent bool
	PrettyPrint bool
	TrimWhitespace bool
}

// DefaultJSONOptions returns the default JSON conversion options.
func DefaultJSONOptions() JSONOptions {
	return JSONOptions{
		IncludeAttributes: false,
		IncludeTextContent: true,
		PrettyPrint: true,
		TrimWhitespace: true,
	}
}

// ToJSON converts the HTML document to JSON format using default options.
func (p *Parser) ToJSON() ([]byte, error) {
	return p.ToJSONWithOptions(DefaultJSONOptions())
}

// ToJSONWithOptions converts the HTML document to JSON with custom options.
func (p *Parser) ToJSONWithOptions(opts JSONOptions) ([]byte, error) {
	if p.root == nil {
		return nil, fmt.Errorf("parser has no root node")
	}

	node := htmlNodeToMap(p.root, opts)
	return marshalJSON(node, opts.PrettyPrint)
}

// ToJSON converts a selection to JSON.
func (s *Selection) ToJSON(opts JSONOptions) ([]byte, error) {
	if s.node == nil {
		return nil, fmt.Errorf("selection node is nil")
	}

	node := htmlNodeToMap(s.node, opts)
	return marshalJSON(node, opts.PrettyPrint)
}

// ToMap converts the selection to a map structure.
func (s *Selection) ToMap(opts JSONOptions) map[string]any {
	if s.node == nil {
		return nil
	}
	return htmlNodeToMap(s.node, opts)
}

// ToMap converts the parser document to a map structure.
func (p *Parser) ToMap(opts JSONOptions) map[string]any {
	if p.root == nil {
		return nil
	}
	return htmlNodeToMap(p.root, opts)
}

// htmlNodeToMap converts an HTML node to a map representation.
func htmlNodeToMap(node *html.Node, opts JSONOptions) map[string]any {
	result := map[string]any{
		"type": nodeTypeToString(node.Type),
	}

	if node.Type == html.ElementNode {
		result["tag"] = node.Data

		if opts.IncludeAttributes && len(node.Attr) > 0 {
			attrs := make(map[string]string, len(node.Attr))
			for _, attr := range node.Attr {
				attrs[attr.Key] = attr.Val
			}
			result["attributes"] = attrs
		}
	}

	children, textContent := processChildren(node, opts)

	if len(children) > 0 {
		result["children"] = children
	}

	if opts.IncludeTextContent {
		if len(textContent) > 0 {
			result["text"] = textContent
		} else if len(children) == 0 {
			text := htmlquery.InnerText(node)
			if opts.TrimWhitespace {
				text = trimText(text)
			}
			if text != "" {
				result["text"] = []string{text}
			}
		}
	}

	return result
}

// processChildren processes child nodes and returns children and text content.
func processChildren(node *html.Node, opts JSONOptions) ([]map[string]any, []string) {
	var children []map[string]any
	var textContent []string

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		switch child.Type {
		case html.TextNode:
			text := child.Data
			if opts.TrimWhitespace {
				text = trimText(text)
			}
			if text != "" {
				textContent = append(textContent, text)
			}
		case html.ElementNode:
			children = append(children, htmlNodeToMap(child, opts))
		}
	}

	return children, textContent
}

// nodeTypeToString converts a node type constant to a string.
func nodeTypeToString(nodeType html.NodeType) string {
	switch nodeType {
	case html.ErrorNode:
		return "error"
	case html.TextNode:
		return "text"
	case html.DocumentNode:
		return "document"
	case html.ElementNode:
		return "element"
	case html.CommentNode:
		return "comment"
	case html.DoctypeNode:
		return "doctype"
	default:
		return "unknown"
	}
}

// trimText removes leading and trailing whitespace.
func trimText(text string) string {
	return text
}

// marshalJSON marshals a node to JSON with optional pretty printing.
func marshalJSON(node map[string]any, pretty bool) ([]byte, error) {
	var result []byte
	var err error

	if pretty {
		result, err = json.MarshalIndent(node, "", "  ")
	} else {
		result, err = json.Marshal(node)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return result, nil
}
