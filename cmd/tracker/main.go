package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MedBou33/gold-tracker/internal/provider"
)

func main() {
	client := provider.NewClient("https://api.gold-api.com")

	quote, err := client.GetQuote(context.Background(), "XAU", "EUR")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s: %.2f %s (updated %s)\n", quote.Name, quote.Price, quote.Currency, quote.UpdatedAt)
}
