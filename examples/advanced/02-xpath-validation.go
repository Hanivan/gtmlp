package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Hanivan/gtmlp"
)

func RunXpathValidation() {
	fmt.Println("(o_o) XPath Validation Demo")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\n(>_<) Fetching HTML from scrapingcourse.com...")

	// Fetch HTML first
	p, err := gtmlp.ParseURL("https://www.scrapingcourse.com/ecommerce/",
		gtmlp.WithTimeout(30*1000000000),
	)
	if err != nil {
		log.Fatalf("(x_x) Error during scraping: %v", err)
	}

	fmt.Println("(^_^) HTML fetched successfully")

	// Define XPath patterns to test
	xpathPatterns := []string{
		// Valid patterns
		"//title/text()",
		"//li[contains(@class, 'product')]",
		"//h2/text()",
		"//span[@class='price']//bdi/text()",
		"//img/@src",

		// Patterns that might not match
		"//div[@class='non-existent']",
		"//meta[@name='author']/@content",

		// Invalid XPath syntax
		"//invalid[@xpath[syntax",
		"//unclosed[@bracket",
	}

	fmt.Println("(o_o) Testing XPath Patterns:")
	fmt.Println(strings.Repeat("─", 80))

	// Validate all patterns
	validation := gtmlp.ValidateXPathWithParser(p, xpathPatterns)

	// Display results
	for i, result := range validation {
		status := "(^_^)"
		if !result.Valid {
			status = "(x_x)"
		}
		matchInfo := ""
		if result.Valid {
			matchCount := "match"
			if result.MatchCount != 1 {
				matchCount = "matches"
			}
			matchInfo = fmt.Sprintf(" (%d %s)", result.MatchCount, matchCount)
		}

		fmt.Printf("%d. %s %s%s\n", i+1, status, result.XPath, matchInfo)

		if result.Valid && result.Sample != "" {
			// Show sample value (truncate if too long)
			sample := result.Sample
			if len(sample) > 60 {
				sample = sample[:60] + "..."
			}
			fmt.Printf("   (._.) Sample: \"%s\"\n", sample)
		}

		if !result.Valid && result.Error != "" {
			fmt.Printf("   (o_o) Error: %s\n", result.Error)
		}

		fmt.Println("")
	}

	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("\n(._.) Validation Summary:")

	validCount := 0
	for _, r := range validation {
		if r.Valid {
			validCount++
		}
	}
	invalidCount := len(validation) - validCount

	totalMatches := 0
	for _, r := range validation {
		totalMatches += r.MatchCount
	}

	fmt.Printf("   Total Patterns Tested: %d\n", len(validation))
	fmt.Printf("   (^_^) Valid: %d\n", validCount)
	fmt.Printf("   (x_x) Invalid: %d\n", invalidCount)
	fmt.Printf("   (._.) Total Matches Found: %d\n", totalMatches)

	overallStatus := "PASSED"
	for _, r := range validation {
		if !r.Valid {
			overallStatus = "FAILED"
			break
		}
	}
	fmt.Printf("   (☆^O^☆) Overall Status: %s\n", overallStatus)

	// Demonstrate pattern refinement
	fmt.Println("\n\n(・_・) Pattern Refinement Demo")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\nTesting different approaches for product titles:")

	titlePatterns := []string{
		"//h2/text()", // Direct approach
		"//li[contains(@class, 'product')]//h2/text()",           // More specific
		".//h2[@class='woocommerce-loop-product__title']/text()", // Full class
		"//a/h2/text()", // Through link
	}

	for i, pattern := range titlePatterns {
		results := gtmlp.ValidateXPathWithParser(p, []string{pattern})
		matchCount := results[0].MatchCount
		sample := results[0].Sample
		if sample == "" {
			sample = "N/A"
		}

		fmt.Printf("%d. Pattern: %s\n", i+1, pattern)
		fmt.Printf("   Matches: %d\n", matchCount)
		if len(sample) > 50 {
			sample = sample[:50] + "..."
		}
		fmt.Printf("   Sample: %s\n", sample)
		fmt.Println("")
	}

	fmt.Println("\\(^o^)/ XPath validation demo completed!")
}
