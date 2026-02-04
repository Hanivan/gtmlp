package main

import (
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

// Product represents an item from the ecommerce test site
type EcommerceProduct struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

func main() {
	// Load selector configuration
	config, err := gtmlp.LoadConfig("ecommerce_selectors.json", nil)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Scrape the ecommerce test site
	url := "https://www.scrapingcourse.com/ecommerce/"
	products, err := gtmlp.ScrapeURL[EcommerceProduct](url, config)
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
