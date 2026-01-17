package main

import (
	"fmt"
	"log"

	"github.com/Hanivan/gtmlp"
)

func RunChainablePatternExample() {
	fmt.Println("=== Chainable Pattern API Examples ===")
	fmt.Println()

	// Sample HTML with product data
	html := `
		<html>
			<body>
				<div class="products">
					<div class="product">
						<h2>  Wireless Mouse  </h2>
						<span class="price">$  29.99  </span>
						<span class="rating">4.5</span>
						<img src="/images/mouse.jpg" />
					</div>
					<div class="product">
						<h2>  Mechanical Keyboard  </h2>
						<span class="price">$  89.99  </span>
						<span class="rating">4.8</span>
						<img src="/images/keyboard.jpg" />
					</div>
					<div class="product">
						<h2>  USB-C Cable  </h2>
						<span class="price">$  12.99  </span>
						<span class="rating">4.2</span>
						<img src="/images/cable.jpg" />
					</div>
				</div>
			</body>
		</html>
	`

	// Define extraction patterns
	patterns := []gtmlp.PatternField{
		// Container pattern - defines each product
		gtmlp.NewContainerPattern("container", "//div[@class='product']"),

		// Field patterns - extracted from each container
		{
			Key:        "name",
			Patterns:   []string{".//h2/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
			},
		},
		{
			Key:        "price",
			Patterns:   []string{".//span[@class='price']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
			},
		},
		{
			Key:        "rating",
			Patterns:   []string{".//span[@class='rating']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
		},
		{
			Key:        "image",
			Patterns:   []string{".//img/@src"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
		},
	}

	// Example 1: Chainable API with FromHTML
	fmt.Println("Example 1: Using FromHTML with patterns (chainable)")
	fmt.Println("-----------------------------------------------------")

	results, err := gtmlp.FromHTML(html).
		WithPatterns(patterns).
		Extract()

	if err != nil {
		log.Fatalf("(x_x) Error: %v", err)
	}

	fmt.Printf("(^_^) Extracted %d products:\n\n", len(results))
	for i, product := range results {
		fmt.Printf("%d. %s\n", i+1, product["name"])
		fmt.Printf("   (._.) Price:  %s\n", product["price"])
		fmt.Printf("   (._.) Rating: %s\n", product["rating"])
		fmt.Printf("   (._.) Image:  %s\n", product["image"])
		fmt.Println()
	}

	// Example 2: Builder pattern with New()
	fmt.Println("Example 2: Using New() builder")
	fmt.Println("--------------------------------")

	builderResults, err := gtmlp.New().
		FromHTML(html).
		WithPatterns(patterns).
		Extract()

	if err != nil {
		log.Fatalf("(x_x) Error: %v", err)
	}

	fmt.Printf("(^_^) Extracted %d products using builder\n", len(builderResults))

	// Example 3: Alternative patterns with chainable API
	fmt.Println("\nExample 3: Alternative patterns (chainable)")
	fmt.Println("--------------------------------------------")

	alternativeHTML := `
		<html>
			<body>
				<div class="article">
					<meta property="og:title" content="Advanced Go Patterns" />
					<h1>Fallback Title</h1>
					<meta name="description" content="Learn advanced patterns" />
					<p>Article content here</p>
				</div>
			</body>
		</html>
	`

	altPatterns := []gtmlp.PatternField{
		{
			Key: "title",
			Patterns: []string{
				"//meta[@property='og:title']/@content", // Try Open Graph first
				"//h1/text()",                           // Fallback to h1
				"//title/text()",                        // Finally try title tag
			},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
		},
		{
			Key: "description",
			Patterns: []string{
				"//meta[@property='og:description']/@content",
				"//meta[@name='description']/@content",
				"//p[1]/text()",
			},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
		},
	}

	altResults, err := gtmlp.FromHTML(alternativeHTML).
		WithPatterns(altPatterns).
		Extract()

	if err != nil {
		log.Fatalf("(x_x) Error: %v", err)
	}

	if len(altResults) > 0 {
		fmt.Printf("(^_^) Title: %s\n", altResults[0]["title"])
		fmt.Printf("(^_^) Description: %s\n", altResults[0]["description"])
	}

	fmt.Println("\n(^o^) Chainable pattern API examples completed successfully!")
}
