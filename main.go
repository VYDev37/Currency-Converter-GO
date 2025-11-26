package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func HasFlag(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . --convert -from [base_currency_id] -to [target_currency_id] -amount [amount].")
		return
	}

	fromCurrency := flag.String("from", "", "From")
	toCurrency := flag.String("to", "", "To")
	amount := flag.Float64("amount", 0.0, "Amount")

	flag.CommandLine.Parse(os.Args[2:])
	commandName := strings.ToLower(os.Args[1])

	if commandName == "--convert" || commandName == "convert" {
		if !HasFlag("from") || !HasFlag("to") || !HasFlag("amount") {
			fmt.Println("Usage: go run . --convert -from [base_currency_id] -to [target_currency_id] -amount [amount].")
			return
		}

		rate, err := Convert(*fromCurrency, *toCurrency, *amount)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Conversion result: ")
		fmt.Printf("From: %s\nTo: %s\nAmount: %.2f\nRate: %.2f\nTotal: %.2f", rate.From, rate.To, rate.Amount, rate.Rate, rate.Result)
	} else {
		fmt.Println("The only available command now is convert.")
	}
}
