package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Hanivan/gtmlp"
)

// Product represents a scraped product
type Product struct {
	Name  string
	Price string
	Image string
}

func RunBasicProductScraping() {
	fmt.Println("(>_<) Basic Product Scraping Demo")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\n(>_<) Fetching product data from scrapingcourse.com...")

	// Define extraction patterns
	patterns := []gtmlp.PatternField{
		// Container pattern - defines each product
		gtmlp.NewContainerPattern("container", "//li[contains(@class, 'product')]"),

		// Field patterns - extracted from each container
		{
			Key:        "name",
			Patterns:   []string{".//h2/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
				gtmlp.NewReplacePipe(`\s+`, " "), // Replace multiple spaces with single space
			},
		},
		{
			Key:        "price",
			Patterns:   []string{".//span[@class='price']//bdi/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "image",
			Patterns:   []string{".//img/@src"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
		},
	}

	// Execute scraping
	p, err := gtmlp.ParseURL("https://www.scrapingcourse.com/ecommerce/",
		gtmlp.WithTimeout(30*1000000000),
	)
	if err != nil {
		log.Fatalf("(x_x) Error during scraping: %v", err)
	}

	results, err := gtmlp.ExtractWithPatterns(p, patterns)
	if err != nil {
		log.Fatalf("(x_x) Error extracting patterns: %v", err)
	}

	fmt.Printf("(^_^) Successfully extracted %d products\n\n", len(results))

	// Display first 5 products
	fmt.Println("(>_<) First 5 Products:")
	for i, product := range results {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, product["name"])
		fmt.Printf("   (._.) Price: %s\n", product["price"])
		fmt.Printf("   (._.) Image: %s\n", product["image"])
		fmt.Println("")
	}

	// Analytics
	fmt.Println("(._.) Analytics:")
	fmt.Printf("   Total Products: %d\n", len(results))

	productsWithImage := 0
	totalNameLength := 0
	for _, p := range results {
		if img, ok := p["image"].(string); ok && img != "" {
			productsWithImage++
		}
		if name, ok := p["name"].(string); ok {
			totalNameLength += len(name)
		}
	}

	fmt.Printf("   Products with Images: %d\n", productsWithImage)
	if len(results) > 0 {
		avgNameLength := totalNameLength / len(results)
		fmt.Printf("   Average Name Length: %d characters\n", avgNameLength)
	}

	fmt.Println("\n\\(^o^)/ Basic product scraping demo completed!")
}
