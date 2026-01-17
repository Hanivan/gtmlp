package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Hanivan/gtmlp"
)

func RunBuilderExample() {
	fmt.Println("=== Builder/Fluent API Examples ===")
	fmt.Println()

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Builder API Example</title>
	</head>
	<body>
		<div class="article">
			<h1 id="title">Fluent API for HTML Parsing</h1>
			<p class="author">By John Doe</p>
			<div class="content">
				<p>This library provides a clean, fluent API for parsing HTML.</p>
			</div>
		</div>
	</body>
	</html>
	`

	// Example 1: Basic builder usage
	fmt.Println("Example 1: Get text using builder")
	result, err := gtmlp.FromHTML(html).XPath("//title")
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	if result != nil {
		fmt.Printf("  Title: %s\n\n", result.Text())
	}

	// Example 2: Chain operations
	fmt.Println("Example 2: Chain operations")
	text, err := gtmlp.FromHTML(html).Text("//p[@class='author']")
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Printf("  Author: %s\n\n", text)

	// Example 3: Convert to JSON
	fmt.Println("Example 3: Convert to JSON using builder")
	json, err := gtmlp.New().
		FromHTML(html).
		ToJSON()
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Printf("  JSON: %s\n\n", string(json))

	// Example 4: Get multiple elements
	fmt.Println("Example 4: Get all paragraphs using builder")
	items, err := gtmlp.FromHTML(html).XPathAll("//p")
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	for i, item := range items {
		fmt.Printf("  Paragraph %d: %s\n", i+1, item.Text())
	}

	// Example 5: With options (simulated)
	fmt.Println("\nExample 5: Builder with options (simulated)")
	builder := gtmlp.New().
		WithTimeout(10 * time.Second).
		WithUserAgent("CustomBot/1.0").
		WithHeaders(map[string]string{
			"Accept-Language": "en-US",
		})
	_ = builder // In real usage, you would call .FromURL(url) here
	fmt.Println("  Builder configured with custom options")

	fmt.Println("\n(^o^) Builder examples completed successfully!")
}
