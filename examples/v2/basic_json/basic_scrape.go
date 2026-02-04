package main

import (
    "fmt"
    "log"
    "github.com/Hanivan/gtmlp"
)

type Product struct {
    Name  string `json:"name"`
    Price string `json:"price"`
    Link  string `json:"link"`
}

func main() {
    config, err := gtmlp.LoadConfig("selectors.json", nil)
    if err != nil {
        log.Fatal(err)
    }

    products, err := gtmlp.ScrapeURL[Product]("https://example.com/products", config)
    if err != nil {
        log.Fatal(err)
    }

    for _, p := range products {
        fmt.Printf("%s: %s - %s\n", p.Name, p.Price, p.Link)
    }
}
