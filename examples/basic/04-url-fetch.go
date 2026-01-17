package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Hanivan/gtmlp"
)

func RunURLFetchExample() {
	fmt.Println("=== URL Fetch Example ===")
	fmt.Println()

	// Example using a real website (example.com is safe for testing)
	url := "https://example.com"

	// Fetch and parse with custom options
	parser, err := gtmlp.ParseURL(url,
		gtmlp.WithTimeout(10*time.Second),
		gtmlp.WithUserAgent("GTMLP-Bot/1.0"),
		gtmlp.WithMaxRetries(2),
	)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	fmt.Printf("Successfully fetched and parsed: %s\n\n", url)

	// Get the title
	title, _ := parser.XPath("//title")
	if title != nil {
		fmt.Printf("Page Title: %s\n", title.Text())
	}

	// Get the main heading
	h1, _ := parser.XPath("//h1")
	if h1 != nil {
		fmt.Printf("Main Heading: %s\n", h1.Text())
	}

	// Get all paragraphs
	paragraphs, _ := parser.XPathAll("//p")
	fmt.Printf("\nFound %d paragraphs:\n", len(paragraphs))
	for i, p := range paragraphs {
		text := p.TextTrimmed()
		if text != "" {
			fmt.Printf("  %d. %s\n", i+1, text)
		}
	}

	// Get all links
	links, _ := parser.XPathAll("//a[@href]")
	fmt.Printf("\nFound %d links:\n", len(links))
	for i, link := range links {
		href := link.Attr("href")
		text := link.TextTrimmed()
		if text == "" {
			text = href
		}
		fmt.Printf("  %d. %s -> %s\n", i+1, text, href)
	}

	fmt.Println("\n(^o^) URL fetch example completed successfully!")
}
