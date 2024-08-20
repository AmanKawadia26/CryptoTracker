package crypto

import (
	"cryptotracker/models"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
)

func SearchCryptocurrency() {
	var symbol string
	color.New(color.FgCyan).Print("Enter the symbol of the cryptocurrency: ")
	fmt.Scan(&symbol)

	params := map[string]string{
		"symbol":  symbol,
		"convert": "USD",
	}

	response := getAPIResponse("/quotes/latest", params)

	var result map[string]interface{}
	err := json.Unmarshal(response, &result)
	if err != nil {
		color.New(color.FgRed).Printf("Error unmarshalling API response: %v\n", err)
		return
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		color.New(color.FgRed).Println("Data not found in the response.")
		return
	}

	cryptoData, ok := data[symbol].(map[string]interface{})
	if !ok {
		color.New(color.FgYellow).Printf("Cryptocurrency data not found for symbol: %s\n", symbol)
		color.New(color.FgMagenta).Println("Please request the addition of this cryptocurrency to our app.")

		// Create a new request for the unavailable cryptocurrency
		request := &models.UnavailableCryptoRequest{
			CryptoSymbol:   symbol,
			RequestMessage: "Please add this cryptocurrency.",
			Status:         "Pending",
		}

		// Save the request to the file
		if err := saveUnavailableCryptoRequest(request); err != nil {
			color.New(color.FgRed).Printf("Error saving unavailable crypto request: %v\n", err)
			return
		}
		color.New(color.FgGreen).Println("Request to add the cryptocurrency has been submitted.")
		return
	}

	// Proceed with normal processing if the cryptocurrency is found
	name, ok := cryptoData["name"].(string)
	if !ok {
		color.New(color.FgRed).Println("Could not retrieve the name of the cryptocurrency.")
		return
	}

	priceObj, ok := cryptoData["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"]
	if !ok {
		color.New(color.FgRed).Println("Could not retrieve the price of the cryptocurrency.")
		return
	}

	price, ok := priceObj.(float64)
	if !ok {
		color.New(color.FgRed).Println("Failed to convert price to float64.")
		return
	}

	color.New(color.FgGreen).Printf("%s (%s): $%.2f\n", name, symbol, price)
	fmt.Println()
}
