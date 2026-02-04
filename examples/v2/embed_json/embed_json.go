package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

//go:embed selectors.json
var selectorsJSON string

// Product represents a scraped product item
type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

func main() {
	// Parse embedded config using ParseConfig (no file I/O needed)
	// This is useful for:
	// - Single-binary deployments (no external config files needed)
	// - Versioned configurations (config travels with code)
	// - Multiple configs in one binary (embed multiple files)
	config, err := gtmlp.ParseConfig(selectorsJSON, gtmlp.FormatJSON, nil)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Scrape the ecommerce test site
	url := "https://www.scrapingcourse.com/ecommerce/"
	products, err := gtmlp.ScrapeURL[Product](context.Background(), url, config)
	if err != nil {
		log.Fatalf("Failed to scrape: %v", err)
	}

	// Print results
	fmt.Printf("Found %d products:\n\n", len(products))
	for i, p := range products {
		fmt.Printf("%d. %s\n", i+1, p.Name)
		fmt.Printf("   Price: %s\n", p.Price)
		fmt.Printf("   Link: %s\n\n", p.Link)
	}
}
