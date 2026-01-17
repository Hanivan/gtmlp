package main

import (
	"fmt"
	"strings"

	"github.com/Hanivan/gtmlp"
)

func RunConfigurationOptions() {
	fmt.Println("\n=== Example: Configuration Options ===")

	// Example 1: WithSuppressErrors
	fmt.Println("=== Example 1: Suppress XPath Errors ===")

	html1 := `
		<html>
			<body>
				<h1>Article Title</h1>
				<p class="content">Article content goes here.</p>
			</body>
		</html>
	`

	// Without error suppression
	fmt.Println("Without error suppression:")
	p1, _ := gtmlp.Parse(html1)
	_, err := p1.XPath("//meta[@property='og:title']/@content")
	if err != nil {
		fmt.Printf("   XPath error: %v\n", err)
	}

	// With error suppression (using WithSuppressErrors option on ParseURL)
	// For direct parsing, we can use the WithSuppressErrors option
	fmt.Println("\nWith error suppression:")
	// This would be used with ParseURL, but for parsed HTML we handle errors manually

	// Using alternative patterns instead of relying on error suppression
	patterns1 := []gtmlp.PatternField{
		{
			Key:          "title",
			Patterns:     []string{"//h1/text()"},
			ReturnType:   gtmlp.ReturnTypeText,
			AlterPattern: []string{"//title/text()"},
			Meta:         gtmlp.DefaultPatternMeta(),
			Pipes:        []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "content",
			Patterns:   []string{".//p[@class='content']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
	}

	result1, _ := gtmlp.ExtractWithPatterns(p1, patterns1)
	fmt.Printf("Title: %v\n", result1[0]["title"])
	fmt.Printf("Content: %v\n", result1[0]["content"])
	fmt.Println("\n✓ No XPath errors - handled gracefully")

	// Example 2: XPath validation with proper error handling
	fmt.Println("\n=== Example 2: XPath Validation ===")

	html2 := `
		<html>
			<body>
				<h1>Test Page</h1>
				<p>Some content</p>
			</body>
		</html>
	`

	p2, _ := gtmlp.Parse(html2)

	// Validate XPath patterns
	validation := gtmlp.ValidateXPathWithParser(p2, []string{
		"//h1/text()",       // Valid
		"//p/text()",        // Valid
		"//invalid[[[xpath", // Invalid - will show error
	})

	fmt.Println("Validation results:")
	for _, result := range validation {
		status := "✓"
		if !result.Valid {
			status = "✗"
		}
		fmt.Printf("  %s %s\n", status, result.XPath)
		if result.Error != "" {
			fmt.Printf("    Error: %s\n", result.Error)
		}
	}

	fmt.Printf("\nOverall valid: %t\n", validation[0].Valid && validation[1].Valid)
	fmt.Println("✓ Validation completed with proper error handling")

	// Example 3: Multiple field extraction with error handling
	fmt.Println("\n=== Example 3: Complete Configuration ===")

	html3 := `
		<html>
			<body>
				<div class="item">
					<span class="name">Item 1</span>
					<span class="value">100</span>
				</div>
				<div class="item">
					<span class="name">Item 2</span>
					<span class="value">200</span>
				</div>
			</body>
		</html>
	`

	patterns3 := []gtmlp.PatternField{
		gtmlp.NewContainerPattern("container", "//div[@class='item']"),
		{
			Key:        "name",
			Patterns:   []string{".//span[@class='name']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "value",
			Patterns:   []string{".//span[@class='value']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
	}

	p3, _ := gtmlp.Parse(html3)
	result3, _ := gtmlp.ExtractWithPatterns(p3, patterns3)

	fmt.Printf("Items: %v\n", len(result3))
	if len(result3) > 0 {
		fmt.Printf("First item: name=%v, value=%v\n", result3[0]["name"], result3[0]["value"])
	}
	fmt.Println("✓ Complete configuration applied")

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("(☆^O^☆) GTMLP Configuration Features:")
	fmt.Println("   (._.) WithSuppressErrors() - Gracefully handle XPath errors")
	fmt.Println("   (._.) WithTimeout() - Configure request timeouts")
	fmt.Println("   (._.) WithMaxRetries() - Retry failed HTTP requests")
	fmt.Println("   (._.) WithUserAgent() - Set custom user agent")
	fmt.Println("   (._.) WithHeaders() - Add custom HTTP headers")
	fmt.Println("   (._.) WithProxy() - Configure proxy for requests")
	fmt.Println("   (._.) ValidateXPath* - Test patterns before scraping")
	fmt.Println("   (._.) ExtractWithPatterns - Robust extraction with fallbacks")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("")

	fmt.Println("\\(^o^)/ Configuration options demo completed!")
}
