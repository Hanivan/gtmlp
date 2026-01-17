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
	case "01", "xpath":
		RunXPathExample()
	case "02", "json":
		RunJSONExample()
	case "03", "builder":
		RunBuilderExample()
	case "04", "url", "url-fetch":
		RunURLFetchExample()
	case "05", "random-ua":
		RunRandomUAExample()
	case "06", "pattern":
		RunPatternExample()
	case "07", "chainable", "chainable-pattern":
		RunChainablePatternExample()
	case "all":
		runAllExamples()
	default:
		fmt.Printf("Unknown example type: %s\n", *exampleType)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("GTMLP Basic Examples")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run examples/basic/main.go -type=<example_type>")
	fmt.Println()
	fmt.Println("Available examples:")
	fmt.Println("  01, xpath              - XPath selector examples")
	fmt.Println("  02, json               - HTML to JSON conversion examples")
	fmt.Println("  03, builder            - Fluent API builder examples")
	fmt.Println("  04, url, url-fetch     - URL fetch and parse examples")
	fmt.Println("  05, random-ua          - Random User-Agent examples")
	fmt.Println("  06, pattern            - Pattern-based extraction examples")
	fmt.Println("  07, chainable          - Chainable pattern API examples")
	fmt.Println("  all                    - Run all examples in sequence")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run examples/basic/main.go -type=01")
	fmt.Println("  go run examples/basic/main.go -type=xpath")
	fmt.Println("  go run examples/basic/main.go -type=all")
}

func runAllExamples() {
	fmt.Println("Running all basic examples in sequence...")
	fmt.Println(strings.Repeat("=", 70))

	examples := []struct {
		name string
		fn   func()
	}{
		{"01 - XPath Selectors", RunXPathExample},
		{"02 - HTML to JSON Conversion", RunJSONExample},
		{"03 - Fluent API Builder", RunBuilderExample},
		{"04 - URL Fetch and Parse", RunURLFetchExample},
		{"05 - Random User-Agent", RunRandomUAExample},
		{"06 - Pattern-Based Extraction", RunPatternExample},
		{"07 - Chainable Pattern API", RunChainablePatternExample},
	}

	for i, example := range examples {
		fmt.Printf("\n[%d/%d] Running: %s\n", i+1, len(examples), example.name)
		fmt.Println(strings.Repeat("-", 70))
		example.fn()
		fmt.Println()
	}

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("\\(^o^)/ All basic examples completed!")
}

func init() {
	// Add import for strings package used in runAllExamples
	// This is a hack but works for examples
	_ = fmt.Sprintf("") // ensure fmt is used
}
