package main

import (
    "log/slog"
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
    // Set log level for development
    gtmlp.SetLogLevel(slog.LevelInfo)
	// Load YAML configuration with numbered pagination
	config, err := gtmlp.LoadConfig("selectors.yaml", nil)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Extract all page URLs first
	fmt.Println("=== Extracting Pagination URLs ===")
	info, err := gtmlp.ExtractPaginationURLs(
		context.Background(),
		"https://example.com/catalog",
		config,
	)
	if err != nil {
		log.Fatalf("URL extraction failed: %v", err)
	}

	fmt.Printf("Found %d pages\n", len(info.URLs))
	fmt.Println("Page URLs:")
	for i, url := range info.URLs {
		fmt.Printf("  %d. %s\n", i+1, url)
	}

	// Scrape specific pages (manual control)
	fmt.Println("\n=== Scraping First 3 Pages Manually ===")
	var allProducts []Product
	for i, pageURL := range info.URLs {
		if i >= 3 {
			break
		}

		// Scrape individual page (disable pagination for single page)
		pageConfig := *config
		pageConfig.Pagination = nil

		products, err := gtmlp.ScrapeURL[Product](
			context.Background(),
			pageURL,
			&pageConfig,
		)
		if err != nil {
			log.Printf("Failed to scrape page %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Page %d: %d products\n", i+1, len(products))
		allProducts = append(allProducts, products...)
	}

	fmt.Printf("\nTotal products from 3 pages: %d\n", len(allProducts))
}
