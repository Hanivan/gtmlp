package main

import (
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

func RunJSONExample() {
	fmt.Println("=== JSON Conversion Examples ===")
	fmt.Println()

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>JSON Example</title>
	</head>
	<body>
		<div id="main" class="content">
			<h1>HTML to JSON</h1>
			<p>Convert HTML documents to structured JSON format</p>
			<div class="metadata">
				<span data-id="123" data-category="tech">Tech</span>
			</div>
		</div>
	</body>
	</html>
	`

	// Example 1: Basic HTML to JSON
	fmt.Println("Example 1: Basic HTML to JSON (default options)")
	json, err := gtmlp.ToJSON(html)
	if err != nil {
		log.Fatalf("Failed to convert to JSON: %v", err)
	}
	fmt.Printf("%s\n\n", string(json))

	// Example 2: JSON with attributes
	fmt.Println("Example 2: JSON with attributes included")
	json, err = gtmlp.ToJSONWithOptions(html, gtmlp.JSONOptions{
		IncludeAttributes:  true,
		IncludeTextContent: true,
		PrettyPrint:        true,
		TrimWhitespace:     true,
	})
	if err != nil {
		log.Fatalf("Failed to convert to JSON: %v", err)
	}
	fmt.Printf("%s\n\n", string(json))

	// Example 3: Convert specific element
	fmt.Println("Example 3: Convert specific element to JSON")
	p, _ := gtmlp.Parse(html)
	div, _ := p.XPath("//div[@id='main']")
	if div != nil {
		json, _ := div.ToJSON(gtmlp.JSONOptions{
			IncludeAttributes: true,
			PrettyPrint:       true,
		})
		fmt.Printf("%s\n", string(json))
	}

	fmt.Println("\n(^o^) JSON examples completed successfully!")
}
