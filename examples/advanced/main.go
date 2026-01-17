package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	// Define flags
	exampleType := flag.String("type", "", "Example type to run:")
	flag.Parse()

	if *exampleType == "" {
		printUsage()
		return
	}

	switch *exampleType {
	case "01", "basic-product-scraping":
		RunBasicProductScraping()
	case "02", "xpath-validation":
		RunXpathValidation()
	case "03", "data-cleaning-pipes":
		RunDataCleaningPipes()
	case "04", "alternative-patterns":
		RunAlternativePatterns()
	case "05", "xml-parsing":
		RunXMLParsing()
	case "06", "real-world-ecommerce":
		RunRealWorldEcommerce()
	case "07", "url-health-check":
		RunURLHealthCheck()
	case "08", "configuration-options":
		RunConfigurationOptions()
	case "09", "custom-pipes":
		RunCustomPipes()
	case "all":
		runAllExamples()
	default:
		fmt.Printf("Unknown example type: %s\n", *exampleType)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("GTMLP Advanced Examples (1:1 from nestjs-xpath-parser)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run examples/advanced/main.go -type=<example_type>")
	fmt.Println()
	fmt.Println("Available examples:")
	fmt.Println("  01, basic-product-scraping    - Container-based extraction with pipes")
	fmt.Println("  02, xpath-validation          - Validate XPath patterns before scraping")
	fmt.Println("  03, data-cleaning-pipes       - Trim, decode, replace, case transformations")
	fmt.Println("  04, alternative-patterns       - Fallback patterns for robust extraction")
	fmt.Println("  05, xml-parsing               - Sitemap and RSS feed parsing")
	fmt.Println("  06, real-world-ecommerce      - Complete product scraping with analytics")
	fmt.Println("  07, url-health-check          - Check URL availability concurrently")
	fmt.Println("  08, configuration-options     - Suppress errors, timeouts, retries")
	fmt.Println("  09, custom-pipes              - Create and use custom transformation pipes")
	fmt.Println("  all                          - Run all examples in sequence")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run examples/advanced/main.go -type=01")
	fmt.Println("  go run examples/advanced/main.go -type=basic-product-scraping")
	fmt.Println("  go run examples/advanced/main.go -type=all")
}

func runAllExamples() {
	fmt.Println("Running all advanced examples in sequence...")
	fmt.Println(strings.Repeat("=", 70))

	examples := []struct {
		name string
		fn   func()
	}{
		{"01 - Basic Product Scraping", RunBasicProductScraping},
		{"02 - XPath Validation", RunXpathValidation},
		{"03 - Data Cleaning Pipes", RunDataCleaningPipes},
		{"04 - Alternative Patterns", RunAlternativePatterns},
		{"05 - XML Parsing", RunXMLParsing},
		{"06 - Real-World E-commerce", RunRealWorldEcommerce},
		{"07 - URL Health Check", RunURLHealthCheck},
		{"08 - Configuration Options", RunConfigurationOptions},
		{"09 - Custom Pipes", RunCustomPipes},
	}

	for i, example := range examples {
		fmt.Printf("\n[%d/%d] Running: %s\n", i+1, len(examples), example.name)
		fmt.Println(strings.Repeat("-", 70))
		example.fn()
		fmt.Println()
	}

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("\\(^o^)/ All advanced examples completed!")
}

func init() {
	// Add import for strings package used in runAllExamples
	// This is a hack but works for examples
	_ = fmt.Sprintf("") // ensure fmt is used
}
