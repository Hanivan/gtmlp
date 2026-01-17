package main

import (
	"fmt"
	"strings"

	"github.com/Hanivan/gtmlp"
)

// CustomPipe for reversing text
type ReversePipe struct{}

func (p *ReversePipe) Process(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func NewReversePipe() *ReversePipe {
	return &ReversePipe{}
}

func RunCustomPipes() {
	fmt.Println("(>_<) Custom Pipes and Predefined Pipes Demo")
	fmt.Println(strings.Repeat("=", 60))

	// Register custom pipe
	gtmlp.RegisterPipe("reverse", func() gtmlp.Pipe {
		return NewReversePipe()
	})

	html := `
		<html>
			<body>
				<div class="article">
					<span class="date">2024-01-15</span>
					<span class="views">1.5K views</span>
					<a href="/anime/episode-1">Watch Episode 1</a>
					<p>Contact: support@example.com</p>
					<h1>Hello World</h1>
				</div>
			</body>
		</html>
	`

	// Example 1: Predefined pipes
	fmt.Println("\n=== Example 1: Predefined Pipes ===")

	patterns1 := []gtmlp.PatternField{
		{
			Key:        "timestamp",
			Patterns:   []string{".//span[@class='date']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewDateFormatPipe("2006-01-02")},
		},
		{
			Key:        "viewCount",
			Patterns:   []string{".//span[@class='views']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewNumberNormalizePipe()},
		},
		{
			Key:        "email",
			Patterns:   []string{".//p[contains(text(), 'Contact:')]/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewExtractEmailPipe()},
		},
	}

	p, _ := gtmlp.Parse(html)
	results1, _ := gtmlp.ExtractWithPatterns(p, patterns1)

	fmt.Println("Predefined Pipes Example Results:")
	if len(results1) > 0 {
		result := results1[0]
		fmt.Printf("  Timestamp: %v (Unix timestamp from 2024-01-15)\n", result["timestamp"])
		fmt.Printf("  View Count: %v (normalized from 1.5K)\n", result["viewCount"])
		fmt.Printf("  Email: %v (extracted from text)\n", result["email"])
	}

	// Example 2: Regex pipe for pattern-based replacements
	fmt.Println("\n=== Example 2: RegexPipe ===")

	html2 := `
		<html>
			<body>
				<div class="anime-info">
					<h1>Judul : Maou no Musume wa Yasashisugiru!!</h1>
					<span class="price">$25.5K</span>
				</div>
			</body>
		</html>
	`

	patterns2 := []gtmlp.PatternField{
		{
			Key:        "title",
			Patterns:   []string{".//h1/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewReplacePipe(`^Judul : `, ""),
			},
		},
		{
			Key:        "price",
			Patterns:   []string{".//span[@class='price']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewReplacePipe(`[$]`, ""),
				gtmlp.NewReplacePipe(`K$`, "000"),
				gtmlp.NewTrimPipe(),
			},
		},
	}

	p2, _ := gtmlp.Parse(html2)
	results2, _ := gtmlp.ExtractWithPatterns(p2, patterns2)

	fmt.Println("\nRegexPipe Example:")
	if len(results2) > 0 {
		result := results2[0]
		fmt.Printf("  Title: %v\n", result["title"])
		fmt.Printf("  Price: %v\n", result["price"])
	}

	// Example 3: URL resolution pipe
	fmt.Println("\n=== Example 3: URL Resolution Pipe ===")

	html3 := `
		<html>
			<body>
				<div class="links">
					<a href="/anime/episode-1">Episode 1</a>
					<a href="episode-2">Episode 2</a>
					<a href="https://other.com/page">External</a>
				</div>
			</body>
		</html>
	`

	patterns3 := []gtmlp.PatternField{
		{
			Key:        "link",
			Patterns:   []string{".//a/@href"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewURLResolvePipe("https://example.com/blabla/blabla")},
		},
	}

	p3, _ := gtmlp.Parse(html3)
	results3, _ := gtmlp.ExtractWithPatterns(p3, patterns3)

	fmt.Println("\nURL Resolution Pipe Example:")
	if len(results3) > 0 {
		for i, link := range results3 {
			if i >= 3 {
				break
			}
			fmt.Printf("  Link %d: %v (resolved relative to base URL)\n", i+1, link["link"])
		}
	}

	// Example 4: Chaining multiple pipes
	fmt.Println("\n=== Example 4: Chain Pipes ===")

	html4 := `
		<html>
			<body>
				<div class="product">
					<span class="price">    $25.5K    </span>
				</div>
			</body>
		</html>
	`

	patterns4 := []gtmlp.PatternField{
		{
			Key:        "price",
			Patterns:   []string{".//span[@class='price']/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{
				gtmlp.NewTrimPipe(),
				gtmlp.NewReplacePipe(`^\$`, ""),
				gtmlp.NewNumberNormalizePipe(),
			},
		},
	}

	p4, _ := gtmlp.Parse(html4)
	results4, _ := gtmlp.ExtractWithPatterns(p4, patterns4)

	fmt.Println("\nChain Pipes Example:")
	if len(results4) > 0 {
		result := results4[0]
		fmt.Printf("  Price: %v (trimmed -> $ removed -> 25.5K normalized to 25500)\n", result["price"])
	}

	// Example 5: Custom reverse pipe
	fmt.Println("\n=== Example 5: Custom Pipe ===")

	patterns5 := []gtmlp.PatternField{
		{
			Key:        "reversed",
			Patterns:   []string{".//h1/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{NewReversePipe()},
		},
	}

	p5, _ := gtmlp.Parse(html)
	results5, _ := gtmlp.ExtractWithPatterns(p5, patterns5)

	fmt.Println("\nCustom Pipe Example:")
	if len(results5) > 0 {
		result := results5[0]
		fmt.Printf("  Original: \"Hello World\"")
		fmt.Printf("  Reversed: %v\n", result["reversed"])
	}

	// Example 6: List all available pipes
	fmt.Println("\n=== Example 6: Available Pipes ===")

	allPipes := gtmlp.ListPipes()
	fmt.Println("Built-in pipes:")
	for _, pipeName := range allPipes {
		fmt.Printf("  - %s\n", pipeName)
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("(☆^O^☆) Custom Pipe Features:")
	fmt.Println("   (._.) Reuse transformation logic across patterns")
	fmt.Println("   (._.) Encapsulate complex data processing")
	fmt.Println("   (._.) Extend library with domain-specific operations")
	fmt.Println("   (._.) Combine custom pipes with built-in pipes")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("")

	fmt.Println("\\(^o^)/ Custom pipes demo completed!")
}
