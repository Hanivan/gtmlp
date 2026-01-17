package main

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/Hanivan/gtmlp"
)

// ProductListing represents a product from e-commerce site
type ProductListing struct {
	Name   string
	Price  string
	Rating string
	Image  string
	Link   string
}

func RunRealWorldEcommerce() {
	fmt.Println("(>_<) Real-World E-commerce Scraping Demo")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n(>_<) Fetching products from ScrapingCourse.com...")

	// Define comprehensive product extraction patterns
	patterns := []gtmlp.PatternField{
		gtmlp.NewContainerPattern("container", "//li[contains(@class, 'product')]"),

		{
			Key:        "name",
			Patterns:   []string{".//h2[contains(@class, 'woocommerce-loop-product__title')]/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				".//h2/text()",
				".//a/@title",
			},
			Meta: gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
				gtmlp.NewReplacePipe(`\s+`, " "),
			},
		},
		{
			Key:        "price",
			Patterns:   []string{".//span[@class='price']//bdi/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				".//span[@class='price']/text()",
				".//span[contains(@class, 'amount')]/text()",
			},
			Meta:  gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "rating",
			Patterns:   []string{".//div[contains(@class, 'star-rating')]/@style"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				".//div[@class='star-rating']/@aria-label",
				".//span[@class='rating']/text()",
			},
			Meta: gtmlp.DefaultPatternMeta(),
		},
		{
			Key:        "image",
			Patterns:   []string{".//img/@src"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				".//img/@data-src",
				".//img/@data-lazy-src",
			},
			Meta: gtmlp.DefaultPatternMeta(),
		},
		{
			Key:        "link",
			Patterns:   []string{".//a[contains(@class, 'woocommerce-LoopProduct-link')]/@href"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				".//a/@href",
			},
			Meta: gtmlp.DefaultPatternMeta(),
		},
	}

	// Execute scraping
	p, err := gtmlp.ParseURL("https://www.scrapingcourse.com/ecommerce/",
		gtmlp.WithTimeout(30*1000000000),
	)
	if err != nil {
		log.Fatalf("(x_x) Error during e-commerce scraping: %v", err)
	}

	products, err := gtmlp.ExtractWithPatterns(p, patterns)
	if err != nil {
		log.Fatalf("(x_x) Error extracting products: %v", err)
	}

	fmt.Printf("(^_^) Successfully extracted %d products\n\n", len(products))
	fmt.Println(strings.Repeat("─", 80))

	// Display products in a formatted way
	fmt.Println("\n(>_<) Product Catalog:")

	displayCount := len(products)
	if displayCount > 10 {
		displayCount = 10
	}

	for i := 0; i < displayCount; i++ {
		product := products[i]
		fmt.Printf("%d. %v\n", i+1, product["name"])
		fmt.Printf("   (._.) Price:  %v\n", product["price"])
		fmt.Printf("   (☆^O^☆) Rating: %v\n", product["rating"])
		fmt.Printf("   (._.) Link:   %v\n", product["link"])
		fmt.Printf("   (._.) Image:  %v\n", product["image"])
		fmt.Println("")
	}

	if len(products) > 10 {
		fmt.Printf("   ... and %d more products\n\n", len(products)-10)
	}

	// Analytics
	fmt.Println("\n(._.) Product Analytics:")
	fmt.Println(strings.Repeat("─", 80))

	// Parse prices for analytics
	var parsedPrices []float64
	for _, p := range products {
		priceStr := ""
		if ps, ok := p["price"].(string); ok {
			priceStr = ps
		}
		match := regexp.MustCompile(`[\d,]+\.?\d*`).FindString(priceStr)
		if match != "" {
			price, _ := strconv.ParseFloat(strings.Replace(match, ",", "", -1), 64)
			if price > 0 {
				parsedPrices = append(parsedPrices, price)
			}
		}
	}

	if len(parsedPrices) > 0 {
		sum := 0.0
		minPrice := parsedPrices[0]
		maxPrice := parsedPrices[0]

		for _, price := range parsedPrices {
			sum += price
			if price < minPrice {
				minPrice = price
			}
			if price > maxPrice {
				maxPrice = price
			}
		}

		avgPrice := sum / float64(len(parsedPrices))

		fmt.Println("(._.) Price Statistics:")
		fmt.Printf("   Total Products:  %d\n", len(products))
		fmt.Printf("   Products w/Price: %d\n", len(parsedPrices))
		fmt.Printf("   Average Price:   $%.2f\n", avgPrice)
		fmt.Printf("   Min Price:       $%.2f\n", minPrice)
		fmt.Printf("   Max Price:       $%.2f\n", maxPrice)
		fmt.Println("")
	}

	// Rating statistics
	productsWithRating := 0
	for _, p := range products {
		if r, ok := p["rating"].(string); ok && r != "" {
			productsWithRating++
		}
	}

	fmt.Printf("(☆^O^☆) Rating Statistics:")
	fmt.Printf("   Products with Rating: %d\n", productsWithRating)
	fmt.Printf("   Products without Rating: %d\n", len(products)-productsWithRating)
	fmt.Println("")

	// Image statistics
	productsWithImage := 0
	for _, p := range products {
		if img, ok := p["image"].(string); ok && img != "" {
			productsWithImage++
		}
	}

	fmt.Printf("(._.) Image Statistics:")
	fmt.Printf("   Products with Images: %d\n", productsWithImage)
	fmt.Printf("   Products without Images: %d\n", len(products)-productsWithImage)
	fmt.Println("")

	// Name length statistics
	var nameLengths []int
	for _, p := range products {
		if name, ok := p["name"].(string); ok {
			nameLengths = append(nameLengths, len(name))
		}
	}

	if len(nameLengths) > 0 {
		sum := 0
		maxLen := nameLengths[0]
		minLen := nameLengths[0]

		for _, length := range nameLengths {
			sum += length
			if length > maxLen {
				maxLen = length
			}
			if length < minLen && length > 0 {
				minLen = length
			}
		}

		avgNameLength := float64(sum) / float64(len(nameLengths))

		fmt.Printf("(._.) Product Name Statistics:")
		fmt.Printf("   Average Length: %.0f characters\n", avgNameLength)
		fmt.Printf("   Shortest Name:  %d characters\n", minLen)
		fmt.Printf("   Longest Name:   %d characters\n", maxLen)
		fmt.Println("")
	}

	// Price range categories
	if len(parsedPrices) > 0 {
		priceRanges := map[string]int{
			"Under $50":   0,
			"$50 - $100":  0,
			"$100 - $200": 0,
			"Over $200":   0,
		}

		for _, price := range parsedPrices {
			switch {
			case price < 50:
				priceRanges["Under $50"]++
			case price >= 50 && price < 100:
				priceRanges["$50 - $100"]++
			case price >= 100 && price < 200:
				priceRanges["$100 - $200"]++
			case price >= 200:
				priceRanges["Over $200"]++
			}
		}

		fmt.Printf("(._.) Price Range Distribution:")
		for _, rangeName := range []string{"Under $50", "$50 - $100", "$100 - $200", "Over $200"} {
			count := priceRanges[rangeName]
			percentage := float64(count) / float64(len(parsedPrices)) * 100
			bar := strings.Repeat("█", int(math.Round(float64(count)/2)))
			paddedName := fmt.Sprintf("%-15s", rangeName)
			fmt.Printf("   %s %3d (%.1f%%) %s\n", paddedName, count, percentage, bar)
		}
	}

	fmt.Println("\n\n(☆^O^☆) Key Scraping Insights:")
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("(^_^) Container-based extraction works perfectly for product listings")
	fmt.Println("(^_^) Alternative patterns ensure data extraction even with varied HTML")
	fmt.Println("(^_^) Pipes clean and normalize data automatically")
	fmt.Println("(^_^) Type-safe results enable powerful analytics")
	fmt.Println("(^_^) User-agent rotation helps avoid detection (automatic)")

	fmt.Println("\n\n(._.) Production Tips:")
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("(>_<) Add pagination support to scrape all pages")
	fmt.Println("(._.) Store results in database for historical analysis")
	fmt.Println("(._.) Schedule regular scraping for price monitoring")
	fmt.Println("(._.) Set up alerts for price drops or new products")
	fmt.Println("(._.) Track price trends over time")
	fmt.Println("(☆^O^☆) Filter products by category, price range, or rating")

	fmt.Println("\n\\(^o^)/ Real-world e-commerce scraping demo completed!")
}
