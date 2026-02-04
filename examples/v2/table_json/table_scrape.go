package main

import (
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

// TableProduct represents a product from the table parsing challenge
type TableProduct struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    string `json:"price"`
	Stock    string `json:"stock"`
}

func main() {
	// Load selector configuration
	config, err := gtmlp.LoadConfig("table_selectors.json", nil)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Scrape the table parsing challenge page
	url := "https://www.scrapingcourse.com/table-parsing"
	products, err := gtmlp.ScrapeURL[TableProduct](url, config)
	if err != nil {
		log.Fatalf("Failed to scrape: %v", err)
	}

	// Print results
	fmt.Printf("Found %d products:\n\n", len(products))
	for i, p := range products {
		fmt.Printf("%d. ID: %s | %s | %s | %s | Stock: %s\n",
			i+1, p.ID, p.Name, p.Category, p.Price, p.Stock)
	}
}
