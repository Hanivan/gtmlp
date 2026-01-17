package main

import (
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

func RunXPathExample() {
	fmt.Println("=== XPath Examples ===")
	fmt.Println()

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Example Page</title>
	</head>
	<body>
		<div class="container">
			<h1 class="title">Welcome to GTMLP</h1>
			<p class="description">A Go HTML parsing library</p>
			<ul class="features">
				<li>XPath selectors</li>
				<li>HTML to JSON conversion</li>
				<li>Fluent API</li>
			</ul>
		</div>
		<div class="footer">
			<a href="https://github.com">GitHub</a>
		</div>
	</body>
	</html>
	`

	// Parse HTML
	parser, err := gtmlp.Parse(html)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// Example 1: Get title
	fmt.Println("Example 1: Get page title")
	title, _ := parser.XPath("//title")
	if title != nil {
		fmt.Printf("  Title: %s\n\n", title.Text())
	}

	// Example 2: Get element by class
	fmt.Println("Example 2: Get element by class")
	h1, _ := parser.XPath("//h1[@class='title']")
	if h1 != nil {
		fmt.Printf("  H1 Text: %s\n\n", h1.Text())
	}

	// Example 3: Get multiple elements
	fmt.Println("Example 3: Get all list items")
	items, _ := parser.XPathAll("//li")
	for i, item := range items {
		fmt.Printf("  Item %d: %s\n", i+1, item.Text())
	}
	fmt.Println()

	// Example 4: Get attribute value
	fmt.Println("Example 4: Get link href")
	link, _ := parser.XPath("//a[@href]")
	if link != nil {
		fmt.Printf("  Link URL: %s\n", link.Attr("href"))
		fmt.Printf("  Link Text: %s\n\n", link.Text())
	}

	// Example 5: Navigate the DOM
	fmt.Println("Example 5: Navigate DOM structure")
	container, _ := parser.XPath("//div[@class='container']")
	if container != nil {
		firstChild := container.FirstChild()
		if firstChild != nil {
			fmt.Printf("  First child tag: %s\n", firstChild.Attr("class"))
		}
	}

	fmt.Println("\n(^o^) XPath examples completed successfully!")
}
