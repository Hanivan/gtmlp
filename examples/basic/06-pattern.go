package main

import (
	"fmt"

	"github.com/Hanivan/gtmlp"
)

func RunPatternExample() {
	fmt.Println("=== Pattern-Based Extraction Examples ===")
	fmt.Println()

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Product Catalog</title>
	</head>
	<body>
		<div class="products">
			<div class="product">
				<h2 class="name">Laptop</h2>
				<span class="price">$999.99</span>
				<span class="category">Electronics</span>
			</div>
			<div class="product">
				<h2 class="name">Mouse</h2>
				<span class="price">$29.99</span>
				<span class="category">Electronics</span>
			</div>
			<div class="product">
				<h2 class="name">Desk Chair</h2>
				<span class="price">$199.99</span>
				<span class="category">Furniture</span>
			</div>
		</div>
	</body>
	</html>
	`

	// Example 1: Extract single field
	fmt.Println("Example 1: Extract single field")
	p, _ := gtmlp.Parse(html)
	field := gtmlp.NewPatternField("title", "//title")
	result, _ := gtmlp.ExtractSingle(p, field)
	fmt.Printf("  Title: %v\n\n", result)

	// Example 2: Extract multiple values as array
	fmt.Println("Example 2: Extract multiple values as array")
	pricesField := gtmlp.NewPatternFieldWithMultiple("prices", "//span[@class='price']", gtmlp.MultipleArray)
	prices, _ := gtmlp.ExtractSingle(p, pricesField)
	if arr, ok := prices.([]string); ok {
		fmt.Printf("  Prices: %v\n\n", arr)
	}

	// Example 3: Extract with pipes
	fmt.Println("Example 3: Extract with transformation pipes")
	meta := gtmlp.DefaultPatternMeta()
	meta.Multiple = gtmlp.MultipleArray
	namesField := gtmlp.PatternField{
		Key:        "names",
		Patterns:   []string{"//h2[@class='name']"},
		ReturnType: gtmlp.ReturnTypeText,
		Meta:       meta,
		Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe(), gtmlp.NewUpperCasePipe()},
	}
	names, _ := gtmlp.ExtractSingle(p, namesField)
	if arr, ok := names.([]string); ok {
		fmt.Printf("  Names (uppercased): %v\n\n", arr)
	}

	// Example 4: Extract with container pattern
	fmt.Println("Example 4: Extract structured data with containers")
	patterns := []gtmlp.PatternField{
		gtmlp.NewContainerPattern("products", "//div[@class='product']"),
		gtmlp.NewPatternField("name", ".//h2[@class='name']"),
		gtmlp.NewPatternField("price", ".//span[@class='price']"),
		gtmlp.NewPatternField("category", ".//span[@class='category']"),
	}
	products, _ := gtmlp.ExtractWithPatterns(p, patterns)
	fmt.Printf("  Extracted %d products:\n", len(products))
	for i, product := range products {
		fmt.Printf("    %d. %s - %s (%s)\n", i+1, product["name"], product["price"], product["category"])
	}
	fmt.Println()

	// Example 5: Extract with alternative patterns
	fmt.Println("Example 5: Extract with alternative/fallback patterns")
	altHTML := `<html><body><h3 class="title">Fallback Title</h3></body></html>`
	p2, _ := gtmlp.Parse(altHTML)
	altField := gtmlp.PatternField{
		Key:          "heading",
		Patterns:     []string{"//h1"},        // Won't match
		AlterPattern: []string{"//h2", "//h3"}, // Will match h3
		ReturnType:   gtmlp.ReturnTypeText,
		Meta:         gtmlp.DefaultPatternMeta(),
	}
	altResult, _ := gtmlp.ExtractSingle(p2, altField)
	fmt.Printf("  Found: %v (using alternative pattern)\n\n", altResult)

	// Example 6: Multiple with space separator
	fmt.Println("Example 6: Extract multiple values joined with spaces")
	catsField := gtmlp.NewPatternFieldWithMultiple("categories", "//span[@class='category']", gtmlp.MultipleSpace)
	cats, _ := gtmlp.ExtractSingle(p, catsField)
	fmt.Printf("  Categories: %v\n\n", cats)

	// Example 7: Extract HTML content
	fmt.Println("Example 7: Extract HTML content instead of text")
	htmlField := gtmlp.NewPatternFieldWithHTML("product_html", "//div[@class='product'][1]")
	htmlContent, _ := gtmlp.ExtractSingle(p, htmlField)
	if str, ok := htmlContent.(string); ok {
		fmt.Printf("  HTML length: %d characters\n", len(str))
	}

	fmt.Println("\n(^o^) Pattern-based extraction examples completed successfully!")
}
