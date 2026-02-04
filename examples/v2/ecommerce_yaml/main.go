package main

import (
    "log/slog"
	"context"
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

// EcommerceProduct represents an item from the ecommerce test site
type EcommerceProduct struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Link  string  `json:"link"`
}

func main() {
    // Set log level for development
    gtmlp.SetLogLevel(slog.LevelInfo)
	// Load selector configuration from YAML file
	config, err := gtmlp.LoadConfig("selectors.yaml", nil)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Scrape the ecommerce test site
	url := "https://www.scrapingcourse.com/ecommerce/"
	products, err := gtmlp.ScrapeURL[EcommerceProduct](context.Background(), url, config)
	if err != nil {
		log.Fatalf("Failed to scrape: %v", err)
	}

	// Print results
	fmt.Printf("Found %d products:\n\n", len(products))
	for i, p := range products {
		fmt.Printf("%d. %s\n", i+1, p.Name)
		fmt.Printf("   Price: $%.2f\n", p.Price)
		fmt.Printf("   Link: %s\n\n", p.Link)
	}
}
