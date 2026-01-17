package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Hanivan/gtmlp"
	app "github.com/lib4u/fake-useragent"
)

func RunRandomUAExample() {
	fmt.Println("=== Random User-Agent Demonstration ===")

	// Show different user-agent types
	ua, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize fake-useragent: %v", err)
	}

	fmt.Println("Random User-Agents:")
	fmt.Println("====================")

	// Get 5 random user agents
	for i := 1; i <= 5; i++ {
		uaString := ua.GetRandom()
		fmt.Printf("%d. %s\n\n", i, uaString)
	}

	fmt.Println("\nFiltered by Browser:")
	fmt.Println("====================")

	// Chrome user agents
	fmt.Println("Chrome Desktop:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("  %s\n", ua.Filter().Chrome().Get())
	}

	fmt.Println("\nFirefox Mobile:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("  %s\n", ua.Filter().Firefox().Platform(app.Mobile).Get())
	}

	fmt.Println("\n--- Using Random UAs with gtmlp ---")

	// Now demonstrate using it with gtmlp
	fmt.Println("Fetching URL with random user-agent:")
	url := "https://httpbin.org/user-agent"

	// Create parser with random UA (default behavior)
	parser, err := gtmlp.ParseURL(url,
		gtmlp.WithTimeout(5*time.Second),
		// Random UA is enabled by default, no option needed!
	)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	// Show what we got
	title, _ := parser.XPath("//title")
	if title != nil {
		fmt.Printf("Page Title: %s\n\n", title.Text())
	}

	fmt.Println("Request was made with a realistic browser user-agent!")
	fmt.Println("This helps avoid detection during web scraping.")

	fmt.Println("\n(^o^) Random User-Agent examples completed successfully!")
}
