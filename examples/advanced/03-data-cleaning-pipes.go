package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Hanivan/gtmlp"
)

// CleanedProduct represents a product with cleaned data
type CleanedProduct struct {
	Name          string
	NameUppercase string
	NameLowercase string
	Price         string
	PriceNumeric  string
	Description   string
}

func RunDataCleaningPipes() {
	fmt.Println("(・_・) Data Cleaning Pipes Demo")
	fmt.Println(strings.Repeat("=", 50))

	// Sample HTML with messy data
	sampleHTML := `
		<html>
			<body>
				<div class="product">
					<h2>   Wireless  Mouse    with   Extra   Buttons   </h2>
					<p class="description">Premium &amp; Professional  &quot;Gaming&quot;  Mouse</p>
					<span class="price">$  29.99  USD</span>
					<span class="sku">  SKU-12345  </span>
				</div>
			</body>
		</html>
	`

	fmt.Println("\n(._.) Sample HTML (with messy data):")
	fmt.Println(sampleHTML)

	fmt.Println("\n(o_o) Testing Different Pipe Configurations:")

	// Example 1: Basic trim
	fmt.Println("(o_o) Trim Only:")
	p1, _ := gtmlp.Parse(sampleHTML)
	h1Match := regexp.MustCompile(`<h2>(.*?)</h2>`).FindStringSubmatch(sampleHTML)
	originalName := ""
	if len(h1Match) > 1 {
		originalName = h1Match[1]
	}
	fmt.Printf("   Original: \"%s\"\n", originalName)

	field1 := gtmlp.PatternField{
		Key:        "name",
		Patterns:   []string{"//h2/text()"},
		ReturnType: gtmlp.ReturnTypeText,
		Meta:       gtmlp.DefaultPatternMeta(),
		Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
	}
	result1, _ := gtmlp.ExtractSingle(p1, field1)
	fmt.Printf("   Cleaned:  \"%v\"\n\n", result1)

	// Example 2: Trim + Replace multiple spaces
	fmt.Println("(o_o) Trim + Replace Multiple Spaces:")
	field2 := gtmlp.PatternField{
		Key:        "name",
		Patterns:   []string{"//h2/text()"},
		ReturnType: gtmlp.ReturnTypeText,
		Meta:       gtmlp.DefaultPatternMeta(),
		Pipes: []gtmlp.Pipe{
			gtmlp.NewTrimPipe(),
			gtmlp.NewReplacePipe(`\s+`, " "),
		},
	}
	result2, _ := gtmlp.ExtractSingle(p1, field2)
	fmt.Printf("   Cleaned: \"%v\"\n\n", result2)

	// Example 3: Case transformations
	fmt.Println("(o_o) Case Transformations:")
	meta := gtmlp.DefaultPatternMeta()
	meta.Multiple = gtmlp.MultipleArray

	field3a := gtmlp.PatternField{
		Key:        "nameUppercase",
		Patterns:   []string{"//h2/text()"},
		ReturnType: gtmlp.ReturnTypeText,
		Meta:       meta,
		Pipes: []gtmlp.Pipe{
			gtmlp.NewTrimPipe(),
			gtmlp.NewReplacePipe(`\s+`, " "),
			gtmlp.NewUpperCasePipe(),
		},
	}
	field3b := gtmlp.PatternField{
		Key:        "nameLowercase",
		Patterns:   []string{"//h2/text()"},
		ReturnType: gtmlp.ReturnTypeText,
		Meta:       meta,
		Pipes: []gtmlp.Pipe{
			gtmlp.NewTrimPipe(),
			gtmlp.NewReplacePipe(`\s+`, " "),
			gtmlp.NewLowerCasePipe(),
		},
	}
	result3a, _ := gtmlp.ExtractSingle(p1, field3a)
	result3b, _ := gtmlp.ExtractSingle(p1, field3b)
	if arr, ok := result3a.([]string); ok && len(arr) > 0 {
		fmt.Printf("   Uppercase: \"%v\"\n", arr[0])
	}
	if arr, ok := result3b.([]string); ok && len(arr) > 0 {
		fmt.Printf("   Lowercase: \"%v\"\n\n", arr[0])
	}

	// Example 4: HTML entity decoding
	fmt.Println("(o_o) HTML Entity Decoding:")
	pDescMatch := regexp.MustCompile(`<p class="description">(.*?)</p>`).FindStringSubmatch(sampleHTML)
	rawDesc := ""
	if len(pDescMatch) > 1 {
		rawDesc = pDescMatch[1]
	}
	fmt.Printf("   Raw:     \"%s\"\n", rawDesc)

	field4 := gtmlp.PatternField{
		Key:        "description",
		Patterns:   []string{"//p[@class='description']/text()"},
		ReturnType: gtmlp.ReturnTypeText,
		Meta:       gtmlp.DefaultPatternMeta(),
		Pipes: []gtmlp.Pipe{
			gtmlp.NewTrimPipe(),
			gtmlp.NewDecodePipe(),
			gtmlp.NewReplacePipe(`\s+`, " "),
		},
	}
	result4, _ := gtmlp.ExtractSingle(p1, field4)
	fmt.Printf("   Decoded: \"%v\"\n\n", result4)

	// Example 5: Price cleaning
	fmt.Println("(o_o) Price Cleaning (Multiple Replacements):")
	field5a := gtmlp.PatternField{
		Key:        "price",
		Patterns:   []string{"//span[@class='price']/text()"},
		ReturnType: gtmlp.ReturnTypeText,
		Meta:       gtmlp.DefaultPatternMeta(),
		Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
	}
	field5b := gtmlp.PatternField{
		Key:        "priceNumeric",
		Patterns:   []string{"//span[@class='price']/text()"},
		ReturnType: gtmlp.ReturnTypeText,
		Meta:       gtmlp.DefaultPatternMeta(),
		Pipes: []gtmlp.Pipe{
			gtmlp.NewTrimPipe(),
			gtmlp.NewReplacePipe(`\$`, ""),
			gtmlp.NewReplacePipe(`USD`, ""),
			gtmlp.NewReplacePipe(`\s+`, ""),
		},
	}
	result5a, _ := gtmlp.ExtractSingle(p1, field5a)
	result5b, _ := gtmlp.ExtractSingle(p1, field5b)
	fmt.Printf("   Original: \"%v\"\n", result5a)
	fmt.Printf("   Numeric:  \"%v\"\n\n", result5b)

	// Example 6: Complete product with all pipes
	fmt.Println("\n(>_<) Complete Product Extraction with All Pipes:")

	completePatterns := []gtmlp.PatternField{
		{
			Key:        "name",
			Patterns:   []string{"//h2/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
				gtmlp.NewReplacePipe(`\s+`, " "),
			},
		},
		{
			Key:        "nameUppercase",
			Patterns:   []string{"//h2/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
				gtmlp.NewReplacePipe(`\s+`, " "),
				gtmlp.NewUpperCasePipe(),
			},
		},
		{
			Key:        "nameLowercase",
			Patterns:   []string{"//h2/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
				gtmlp.NewReplacePipe(`\s+`, " "),
				gtmlp.NewLowerCasePipe(),
			},
		},
		{
			Key:        "price",
			Patterns:   []string{"//span[@class='price']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "priceNumeric",
			Patterns:   []string{"//span[@class='price']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
				gtmlp.NewReplacePipe(`\$`, ""),
				gtmlp.NewReplacePipe(`USD`, ""),
				gtmlp.NewReplacePipe(`\s+`, ""),
			},
		},
		{
			Key:        "description",
			Patterns:   []string{"//p[@class='description']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
				gtmlp.NewDecodePipe(),
				gtmlp.NewReplacePipe(`\s+`, " "),
			},
		},
	}

	completeResult, _ := gtmlp.ExtractWithPatterns(p1, completePatterns)
	if len(completeResult) > 0 {
		product := completeResult[0]
		fmt.Println("Product Data:")
		fmt.Printf("  Name:            \"%v\"\n", product["name"])
		fmt.Printf("  Name (Upper):    \"%v\"\n", product["nameUppercase"])
		fmt.Printf("  Name (Lower):    \"%v\"\n", product["nameLowercase"])
		fmt.Printf("  Price:           \"%v\"\n", product["price"])
		fmt.Printf("  Price (Numeric): \"%v\"\n", product["priceNumeric"])
		fmt.Printf("  Description:     \"%v\"\n", product["description"])
	}

	fmt.Println("\n(._.) Pipe Usage Summary:")
	fmt.Println("   (^_^) trim: Remove leading/trailing whitespace")
	fmt.Println("   (^_^) toLowerCase: Convert to lowercase")
	fmt.Println("   (^_^) toUpperCase: Convert to uppercase")
	fmt.Println("   (^_^) decode: Decode HTML entities (&amp; → &, &quot; → \", etc.)")
	fmt.Println("   (^_^) replace: Find and replace with regex support")
	fmt.Println("   (^_^) Multiple replace rules can be chained")

	fmt.Println("\n\\(^o^)/ Data cleaning pipes demo completed!")
}
