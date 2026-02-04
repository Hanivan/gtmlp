package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

func main() {
	// Load selector configuration from YAML file
	config, err := gtmlp.LoadConfig("selectors.yaml", nil)
	if err != nil {
		log.Fatal(err)
	}

	products, err := gtmlp.ScrapeURL[Product](context.Background(), "https://www.scrapingcourse.com/ecommerce/", config)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range products {
		fmt.Printf("%s: %s - %s\n", p.Name, p.Price, p.Link)
	}
}
