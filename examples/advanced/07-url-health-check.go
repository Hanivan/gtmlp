package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Hanivan/gtmlp"
)

func RunURLHealthCheck() {
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("(o_o) Example 07: URL Health Check")
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("")

	// Example 1: Check single URL
	fmt.Println("(o_o) Example 1: Check Single URL")
	fmt.Println(strings.Repeat("━", 80))

	singleResult := gtmlp.CheckURLHealth([]string{"https://example.com"}, 5*time.Second)

	fmt.Println("Single URL check result:")
	fmt.Printf("   URL: %s\n", singleResult[0].URL)
	fmt.Printf("   Alive: %v\n", singleResult[0].Alive)
	fmt.Printf("   Status Code: %d\n", singleResult[0].StatusCode)
	fmt.Println("")

	// Example 2: Check multiple URLs
	fmt.Println("(o_o) Example 2: Check Multiple URLs")
	fmt.Println(strings.Repeat("━", 80))

	testUrls := []string{
		"https://example.com",
		"https://www.scrapingcourse.com/ecommerce/",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/404",
		"https://httpbin.org/status/500",
		"https://this-domain-does-not-exist-12345.com",
	}

	results := gtmlp.CheckURLHealth(testUrls, 10*time.Second)

	fmt.Println("Multiple URLs check results:")
	for _, result := range results {
		if result.Alive {
			fmt.Printf("   (OK) %s\n", result.URL)
			fmt.Printf("        Status: %d\n", result.StatusCode)
		} else {
			fmt.Printf("   (X) %s\n", result.URL)
			fmt.Printf("        Status: %d\n", result.StatusCode)
			fmt.Printf("        Error: %s\n", result.Error)
		}
	}
	fmt.Println("")

	// Example 3: Filter dead URLs
	fmt.Println("(o_o) Example 3: Filter Dead URLs")
	fmt.Println(strings.Repeat("━", 80))

	deadUrls := 0
	aliveUrls := 0
	for _, r := range results {
		if r.Alive {
			aliveUrls++
		} else {
			deadUrls++
		}
	}

	fmt.Printf("Total URLs checked: %d\n", len(results))
	fmt.Printf("   Alive: %d\n", aliveUrls)
	fmt.Printf("   Dead: %d\n", deadUrls)

	if deadUrls > 0 {
		fmt.Println("")
		fmt.Println("Dead URLs found:")
		for _, r := range results {
			if !r.Alive {
				statusCode := "N/A"
				if r.StatusCode != 0 {
					statusCode = fmt.Sprintf("%d", r.StatusCode)
				}
				fmt.Printf("   - %s (%s)\n", r.URL, statusCode)
			}
		}
	}
	fmt.Println("")

	// Example 4: Combine with scraping to find broken links
	fmt.Println("(o_o) Example 4: Scrape URLs and Check Their Health")
	fmt.Println(strings.Repeat("━", 80))

	p, err := gtmlp.ParseURL("https://www.scrapingcourse.com/ecommerce/",
		gtmlp.WithTimeout(30*1000000000),
	)
	if err != nil {
		log.Fatalf("(x_x) Error: %v", err)
	}

	patterns := []gtmlp.PatternField{
		gtmlp.NewContainerPattern("container", "//li[contains(@class, 'product')]"),
		{
			Key:        "name",
			Patterns:   []string{".//h2[contains(@class, 'product-title')]//text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "link",
			Patterns:   []string{".//a[contains(@class, 'product-link')]/@href"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
		},
	}

	scrapedData, _ := gtmlp.ExtractWithPatterns(p, patterns)

	// Extract first 3 product URLs
	var productUrls []string
	for i, r := range scrapedData {
		if i >= 3 {
			break
		}
		if link, ok := r["link"].(string); ok && link != "" && strings.HasPrefix(link, "http") {
			productUrls = append(productUrls, link)
		}
	}

	fmt.Printf("Found %d products, checking first %d URLs...\n", len(scrapedData), len(productUrls))

	if len(productUrls) > 0 {
		healthResults := gtmlp.CheckURLHealth(productUrls, 10*time.Second)

		fmt.Println("")
		for i, result := range healthResults {
			productName := "N/A"
			if i < len(scrapedData) {
				if name, ok := scrapedData[i]["name"].(string); ok {
					productName = name
				}
			}
			fmt.Printf("   Product: \"%s\"\n", productName)
			fmt.Printf("   URL: %s\n", result.URL)
			status := "✗ Dead"
			if result.Alive {
				status = "✓ Alive"
			}
			fmt.Printf("   Status: %s (%d)\n", status, result.StatusCode)
			fmt.Println("")
		}
	}

	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("(☆^O^☆) Use Cases:")
	fmt.Println("   (._.) Verify scraped URLs are valid before storing them")
	fmt.Println("   (._.) Monitor website availability and uptime")
	fmt.Println("   (._.) Check link health in sitemaps")
	fmt.Println("   (._.) Validate API endpoints before making requests")
	fmt.Println("   (._.) Clean up dead links from databases")
	fmt.Println("   (._.) Batch check URLs from crawled data")
	fmt.Println("   (._.) Use proxy for corporate/restricted networks")
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("")

	fmt.Println("\\(^o^)/ URL health check demo completed!")
}
