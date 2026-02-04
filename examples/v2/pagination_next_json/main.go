package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	// Load configuration with pagination
	config, err := gtmlp.LoadConfig("selectors.json", nil)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Example 1: Auto-follow pagination (combined results)
	fmt.Println("=== Example 1: Auto-follow Pagination ===")
	products, err := gtmlp.ScrapeURL[Product](
		context.Background(),
		"https://example.com/products",
		config,
	)
	if err != nil {
		log.Fatalf("Scraping failed: %v", err)
	}

	fmt.Printf("Total products scraped: %d\n", len(products))
	for i, p := range products {
		if i >= 5 {
			break // Show first 5
		}
		fmt.Printf("  %d. %s - $%.2f\n", i+1, p.Name, p.Price)
	}

	// Example 2: Page-separated results
	fmt.Println("\n=== Example 2: Page-Separated Results ===")
	results, err := gtmlp.ScrapeURLWithPages[Product](
		context.Background(),
		"https://example.com/products",
		config,
	)
	if err != nil {
		log.Fatalf("Scraping failed: %v", err)
	}

	fmt.Printf("Total pages: %d\n", results.TotalPages)
	fmt.Printf("Total items: %d\n", results.TotalItems)
	for _, page := range results.Pages {
		fmt.Printf("  Page %d (%s): %d items scraped at %s\n",
			page.PageNum, page.URL, len(page.Items), page.ScrapedAt.Format("15:04:05"))
	}

	// Example 3: Extract URLs only (manual control)
	fmt.Println("\n=== Example 3: Extract Pagination URLs ===")
	info, err := gtmlp.ExtractPaginationURLs(
		context.Background(),
		"https://example.com/products",
		config,
	)
	if err != nil {
		log.Fatalf("URL extraction failed: %v", err)
	}

	fmt.Printf("Pagination type: %s\n", info.Type)
	fmt.Printf("Found %d page URLs:\n", len(info.URLs))
	for i, url := range info.URLs {
		if i >= 3 {
			break // Show first 3
		}
		fmt.Printf("  %d. %s\n", i+1, url)
	}
}
