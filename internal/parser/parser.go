package parser

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// ReturnType defines how to extract content from nodes.
type ReturnType string

const (
	// ReturnTypeText returns plain text content.
	ReturnTypeText ReturnType = "text"
	// ReturnTypeHTML returns HTML content.
	ReturnTypeHTML ReturnType = "html"
)

// Parser represents an HTML document that can be queried with XPath.
type Parser struct {
	root           *html.Node
	url            string
	html           string
	suppressErrors bool
}

// Selection represents a selected HTML node.
type Selection struct {
	node *html.Node
}

// New creates a new Parser from an HTML string.
func New(htmlContent string) (*Parser, error) {
	htmlContent = strings.TrimSpace(htmlContent)
	if htmlContent == "" {
		return nil, fmt.Errorf("HTML content cannot be empty")
	}

	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return &Parser{
		root: doc,
		html: htmlContent,
	}, nil
}

// LoadFromURL creates a Parser by fetching HTML from a URL.
// Deprecated: Use gtmlp.ParseURL() instead, which includes HTTP client functionality.
func LoadFromURL(url string) (*Parser, error) {
	return nil, fmt.Errorf("use gtmlp.ParseURL() instead, which includes HTTP client functionality")
}

// Root returns the root HTML node.
func (p *Parser) Root() *html.Node {
	return p.root
}

// HTML returns the original HTML string.
func (p *Parser) HTML() string {
	return p.html
}

// URL returns the URL if the HTML was loaded from a URL.
func (p *Parser) URL() string {
	return p.url
}

// String returns a string representation of the parser.
func (p *Parser) String() string {
	if p.url != "" {
		return fmt.Sprintf("Parser{url: %s}", p.url)
	}
	return fmt.Sprintf("Parser{length: %d}", len(p.html))
}

// WithSuppressErrors enables error suppression for XPath queries.
// When enabled, XPath errors return nil instead of error values.
func (p *Parser) WithSuppressErrors() *Parser {
	p.suppressErrors = true
	return p
}

// SelectionsFromNodes creates a slice of Selections from a slice of html.Nodes.
func SelectionsFromNodes(nodes []*html.Node) []*Selection {
	selections := make([]*Selection, len(nodes))
	for i, node := range nodes {
		selections[i] = &Selection{node: node}
	}
	return selections
}

// SelectionFromNode creates a Selection from an html.Node.
func SelectionFromNode(node *html.Node) *Selection {
	return &Selection{node: node}
}

// XPath executes an XPath expression and returns the first matching node.
// Returns nil if no match is found.
func (p *Parser) XPath(expr string) (*Selection, error) {
	if expr == "" {
		if p.suppressErrors {
			return nil, nil
		}
		return nil, fmt.Errorf("XPath expression cannot be empty")
	}

	node, err := htmlquery.Query(p.root, expr)
	if err != nil {
		if p.suppressErrors {
			return nil, nil
		}
		return nil, fmt.Errorf("XPath query error: %w", err)
	}

	if node == nil {
		return nil, nil
	}

	return &Selection{node: node}, nil
}

// XPathAll executes an XPath expression and returns all matching nodes.
func (p *Parser) XPathAll(expr string) ([]*Selection, error) {
	if expr == "" {
		if p.suppressErrors {
			return nil, nil
		}
		return nil, fmt.Errorf("XPath expression cannot be empty")
	}

	nodes, err := htmlquery.QueryAll(p.root, expr)
	if err != nil {
		if p.suppressErrors {
			return nil, nil
		}
		return nil, fmt.Errorf("XPath query error: %w", err)
	}

	selections := make([]*Selection, len(nodes))
	for i, node := range nodes {
		selections[i] = &Selection{node: node}
	}

	return selections, nil
}
