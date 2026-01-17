package main

import (
	"fmt"
	"strings"

	"github.com/Hanivan/gtmlp"
)

// SitemapUrl represents a URL from a sitemap
type SitemapUrl struct {
	Loc        string
	Lastmod    string
	Changefreq string
	Priority   string
}

// RssItem represents an item from an RSS feed
type RssItem struct {
	Title       string
	Link        string
	Description string
	PubDate     string
}

func RunXMLParsing() {
	fmt.Println("(._.) XML Parsing Demo")
	fmt.Println(strings.Repeat("=", 50))

	// Sample sitemap XML
	sitemapXML := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2024-01-15</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>https://example.com/about</loc>
    <lastmod>2024-01-10</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>https://example.com/products</loc>
    <lastmod>2024-01-20</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.9</priority>
  </url>
  <url>
    <loc>https://example.com/blog</loc>
    <lastmod>2024-01-22</lastmod>
    <changefreq>daily</changefreq>
    <priority>0.7</priority>
  </url>
</urlset>`

	// Sample RSS feed XML
	rssXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Example Blog</title>
    <link>https://example.com/blog</link>
    <description>Latest blog posts</description>
    <item>
      <title>Getting Started with Web Scraping</title>
      <link>https://example.com/blog/web-scraping-101</link>
      <description>Learn the basics of web scraping with practical examples</description>
      <pubDate>Mon, 15 Jan 2024 10:00:00 GMT</pubDate>
    </item>
    <item>
      <title>Advanced XPath Techniques</title>
      <link>https://example.com/blog/xpath-advanced</link>
      <description>Master XPath for complex data extraction scenarios</description>
      <pubDate>Wed, 17 Jan 2024 14:30:00 GMT</pubDate>
    </item>
    <item>
      <title>Building Robust Scrapers</title>
      <link>https://example.com/blog/robust-scrapers</link>
      <description>Best practices for creating maintainable web scrapers</description>
      <pubDate>Fri, 19 Jan 2024 09:15:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

	fmt.Println("\n(o_o) Parsing Sitemap XML:")
	fmt.Println(strings.Repeat("─", 70))

	// Define sitemap extraction patterns
	sitemapPatterns := []gtmlp.PatternField{
		gtmlp.NewContainerPattern("container", "//url"),
		{
			Key:        "loc",
			Patterns:   []string{".//loc/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "lastmod",
			Patterns:   []string{".//lastmod/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "changefreq",
			Patterns:   []string{".//changefreq/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "priority",
			Patterns:   []string{".//priority/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
	}

	// Extract sitemap URLs
	p1, _ := gtmlp.Parse(sitemapXML)
	sitemapResult, _ := gtmlp.ExtractWithPatterns(p1, sitemapPatterns)

	fmt.Printf("Found %d URLs in sitemap:\n\n", len(sitemapResult))
	for i, url := range sitemapResult {
		fmt.Printf("%d. %v\n", i+1, url["loc"])
		fmt.Printf("   Last Modified: %v\n", url["lastmod"])
		fmt.Printf("   Change Freq:   %v\n", url["changefreq"])
		fmt.Printf("   Priority:      %v\n", url["priority"])
		fmt.Println("")
	}

	fmt.Println("\n(o_o) Parsing RSS Feed XML:")
	fmt.Println(strings.Repeat("─", 70))

	// Define RSS extraction patterns
	rssPatterns := []gtmlp.PatternField{
		gtmlp.NewContainerPattern("container", "//item"),
		{
			Key:        "title",
			Patterns:   []string{".//title/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "link",
			Patterns:   []string{".//link/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "description",
			Patterns:   []string{".//description/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe(), gtmlp.NewDecodePipe()},
		},
		{
			Key:        "pubDate",
			Patterns:   []string{".//pubDate/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
	}

	// Extract RSS items
	p2, _ := gtmlp.Parse(rssXML)
	rssResult, _ := gtmlp.ExtractWithPatterns(p2, rssPatterns)

	fmt.Printf("Found %d blog posts in RSS feed:\n\n", len(rssResult))
	for i, item := range rssResult {
		fmt.Printf("%d. %v\n", i+1, item["title"])
		fmt.Printf("   Link:        %v\n", item["link"])
		fmt.Printf("   Description: %v\n", item["description"])
		fmt.Printf("   Published:   %v\n", item["pubDate"])
		fmt.Println("")
	}

	fmt.Println("\n(o_o) Extracting Channel Metadata from RSS:")
	fmt.Println(strings.Repeat("─", 70))

	// Extract channel-level data (non-container pattern)
	channelPatterns := []gtmlp.PatternField{
		{
			Key:        "title",
			Patterns:   []string{"//channel/title/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "link",
			Patterns:   []string{"//channel/link/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
		{
			Key:        "description",
			Patterns:   []string{"//channel/description/text()"},
			ReturnType: gtmlp.ReturnTypeText,
			Meta:       gtmlp.DefaultPatternMeta(),
			Pipes:      []gtmlp.Pipe{gtmlp.NewTrimPipe()},
		},
	}

	channelResult, _ := gtmlp.ExtractWithPatterns(p2, channelPatterns)

	if len(channelResult) > 0 {
		channel := channelResult[0]
		fmt.Println("Channel Information:")
		fmt.Printf("  Title:       %v\n", channel["title"])
		fmt.Printf("  Link:        %v\n", channel["link"])
		fmt.Printf("  Description: %v\n", channel["description"])
	}

	fmt.Println("\n\n(._.) XML Parsing Summary:")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Println("(^_^) XPath is perfect for navigating XML structure")
	fmt.Println("(^_^) Container patterns work the same as HTML")
	fmt.Println("(^_^) All pipe transformations work with XML data")
	fmt.Println("(^_^) Great for sitemaps, RSS feeds, SOAP responses, etc.")

	fmt.Println("\n\n(._.) Common XML Use Cases:")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Println("(._.) Sitemap parsing (sitemap.xml)")
	fmt.Println("(._.) RSS/Atom feed parsing")
	fmt.Println("(>_<) SOAP API response parsing")
	fmt.Println("(・_・) Configuration file parsing")
	fmt.Println("(._.) Data export file parsing")
	fmt.Println("(._.) XML database exports")

	fmt.Println("\n\\(^o^)/ XML parsing demo completed!")
}
