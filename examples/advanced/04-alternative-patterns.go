package main

import (
	"fmt"
	"strings"

	"github.com/Hanivan/gtmlp"
)

// Article represents an article with metadata
type Article struct {
	Title       string
	Description string
	Author      string
	PublishDate string
	Image       string
}

func RunAlternativePatterns() {
	fmt.Println("(>_<) Alternative Patterns Demo")
	fmt.Println(strings.Repeat("=", 50))

	// Sample HTML with different meta tag formats
	sampleHTML1 := `
		<html>
			<head>
				<title>Advanced TypeScript Patterns</title>
				<meta name="description" content="Learn advanced patterns in TypeScript">
				<meta name="author" content="John Doe">
				<meta name="publish-date" content="2024-01-15">
			</head>
			<body>
				<article>
					<h1>Advanced TypeScript Patterns</h1>
					<img src="/main-image.jpg" alt="Article Image">
				</article>
			</body>
		</html>
	`

	// Sample HTML with Open Graph format
	sampleHTML2 := `
		<html>
			<head>
				<title>React Best Practices</title>
				<meta property="og:title" content="React Best Practices Guide">
				<meta property="og:description" content="Master React with these best practices">
				<meta property="article:author" content="Jane Smith">
				<meta property="article:published_time" content="2024-01-20">
				<meta property="og:image" content="/og-image.jpg">
			</head>
			<body>
				<article>
					<h1>React Best Practices</h1>
				</article>
			</body>
		</html>
	`

	// Sample HTML with minimal metadata
	sampleHTML3 := `
		<html>
			<head>
				<title>Simple Blog Post</title>
			</head>
			<body>
				<article>
					<h1>Simple Blog Post</h1>
					<div class="author-info">
						<span>Written by: Alice Johnson</span>
					</div>
					<time datetime="2024-01-25">January 25, 2024</time>
					<figure>
						<img src="/article-image.jpg" alt="Post Image">
					</figure>
				</article>
			</body>
		</html>
	`

	fmt.Println("\n(._.) Extraction with Alternative Patterns:")

	// Define patterns with fallbacks
	patterns := []gtmlp.PatternField{
		{
			Key:        "title",
			Patterns:   []string{"//meta[@property='og:title']/@content"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				"//h1/text()",    // Fallback to h1
				"//title/text()", // Fallback to title tag
			},
			Meta:  gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "description",
			Patterns:   []string{"//meta[@property='og:description']/@content"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				"//meta[@name='description']/@content", // Standard meta description
				"//article/p[1]/text()",                // First paragraph as fallback
			},
			Meta:  gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{gtmlp.NewTrimPipe(), gtmlp.NewDecodePipe()},
		},
		{
			Key:        "author",
			Patterns:   []string{"//meta[@property='article:author']/@content"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				"//meta[@name='author']/@content",         // Standard meta author
				"//div[@class='author-info']/span/text()", // Extract from author div
				"//a[@rel='author']/text()",               // Author link
			},
			Meta:  gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "publishDate",
			Patterns:   []string{"//meta[@property='article:published_time']/@content"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				"//meta[@name='publish-date']/@content",
				"//time/@datetime",
				"//time/text()",
			},
			Meta:  gtmlp.DefaultPatternMeta(),
			Pipes: []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "image",
			Patterns:   []string{"//meta[@property='og:image']/@content"},
			ReturnType: gtmlp.ReturnTypeText,
			AlterPattern: []string{
				"//article//img/@src", // First image in article
				"//figure//img/@src",  // Image in figure
				"//img/@src",          // Any image
			},
			Meta: gtmlp.DefaultPatternMeta(),
		},
	}

	// Test with first HTML (standard meta tags)
	fmt.Println("(o_o) Extracting from HTML with Standard Meta Tags:")
	fmt.Println(strings.Repeat("─", 60))

	p1, _ := gtmlp.Parse(sampleHTML1)
	result1, _ := gtmlp.ExtractWithPatterns(p1, patterns)

	if len(result1) > 0 {
		article1 := result1[0]
		fmt.Printf("   Title:        \"%v\"\n", article1["title"])
		fmt.Printf("   Description:  \"%v\"\n", article1["description"])
		fmt.Printf("   Author:       \"%v\"\n", article1["author"])
		fmt.Printf("   Publish Date: \"%v\"\n", article1["publishDate"])
		fmt.Printf("   Image:        \"%v\"\n", article1["image"])
	}
	fmt.Println("")

	// Test with second HTML (Open Graph tags)
	fmt.Println("(o_o) Extracting from HTML with Open Graph Tags:")
	fmt.Println(strings.Repeat("─", 60))

	p2, _ := gtmlp.Parse(sampleHTML2)
	result2, _ := gtmlp.ExtractWithPatterns(p2, patterns)

	if len(result2) > 0 {
		article2 := result2[0]
		fmt.Printf("   Title:        \"%v\"\n", article2["title"])
		fmt.Printf("   Description:  \"%v\"\n", article2["description"])
		fmt.Printf("   Author:       \"%v\"\n", article2["author"])
		fmt.Printf("   Publish Date: \"%v\"\n", article2["publishDate"])
		fmt.Printf("   Image:        \"%v\"\n", article2["image"])
	}
	fmt.Println("")

	// Test with third HTML (minimal metadata)
	fmt.Println("(o_o) Extracting from HTML with Minimal Metadata (Using Fallbacks):")
	fmt.Println(strings.Repeat("─", 60))

	p3, _ := gtmlp.Parse(sampleHTML3)
	result3, _ := gtmlp.ExtractWithPatterns(p3, patterns)

	if len(result3) > 0 {
		article3 := result3[0]
		fmt.Printf("   Title:        \"%v\"\n", article3["title"])
		desc := "N/A"
		if article3["description"] != nil {
			desc = fmt.Sprintf("%v", article3["description"])
		}
		fmt.Printf("   Description:  \"%s\"\n", desc)
		fmt.Printf("   Author:       \"%v\"\n", article3["author"])
		fmt.Printf("   Publish Date: \"%v\"\n", article3["publishDate"])
		fmt.Printf("   Image:        \"%v\"\n", article3["image"])
	}
	fmt.Println("")

	fmt.Println("\n(._.) Pattern Fallback Strategy:")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println("For each field, the scraper tries patterns in this order:")
	fmt.Println("")
	fmt.Println("Title:")
	fmt.Println("  1. og:title meta tag (Open Graph)")
	fmt.Println("  2. <h1> element (page heading)")
	fmt.Println("  3. <title> tag (browser title)")
	fmt.Println("")
	fmt.Println("Description:")
	fmt.Println("  1. og:description meta tag (Open Graph)")
	fmt.Println("  2. description meta tag (standard)")
	fmt.Println("  3. First paragraph in article")
	fmt.Println("")
	fmt.Println("Author:")
	fmt.Println("  1. article:author meta tag (Open Graph)")
	fmt.Println("  2. author meta tag (standard)")
	fmt.Println("  3. Author info div content")
	fmt.Println("  4. Author link with rel=\"author\"")
	fmt.Println("")
	fmt.Println("Publish Date:")
	fmt.Println("  1. article:published_time meta tag")
	fmt.Println("  2. publish-date meta tag")
	fmt.Println("  3. <time> datetime attribute")
	fmt.Println("  4. <time> text content")
	fmt.Println("")
	fmt.Println("Image:")
	fmt.Println("  1. og:image meta tag")
	fmt.Println("  2. First image in <article>")
	fmt.Println("  3. Image in <figure> element")
	fmt.Println("  4. Any <img> tag")

	fmt.Println("\n\n(☆^O^☆) Benefits of Alternative Patterns:")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println("(^_^) Robust extraction across different page structures")
	fmt.Println("(^_^) Graceful degradation when primary selectors fail")
	fmt.Println("(^_^) Support for multiple metadata standards (OG, Schema.org, etc.)")
	fmt.Println("(^_^) Reduced scraping failures due to page changes")
	fmt.Println("(^_^) Better data quality with multiple fallback options")

	fmt.Println("\n\\(^o^)/ Alternative patterns demo completed!")
}
