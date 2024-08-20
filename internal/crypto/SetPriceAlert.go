package crypto

import (
	"cryptotracker/models"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
)

func SetPriceAlert(user *models.User) {
	var symbol string
	var targetPrice float64

	color.New(color.FgCyan).Print("Enter the symbol of the cryptocurrency: ")
	fmt.Scan(&symbol)
	color.New(color.FgCyan).Print("Enter your target price in USD: ")
	fmt.Scan(&targetPrice)

	params := map[string]string{
		"symbol":  symbol,
		"convert": "USD",
	}

	response := getAPIResponse("/quotes/latest", params)

	var result map[string]interface{}
	json.Unmarshal(response, &result)

	data := result["data"].(map[string]interface{})[symbol].(map[string]interface{})
	price := data["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"].(float64)

	if price >= targetPrice {
		color.New(color.FgGreen).Printf("Alert: %s has reached your target price of $%.2f. Current price: $%.2f\n", symbol, targetPrice, price)
	} else {
		color.New(color.FgYellow).Printf("%s is still below your target price. Current price: $%.2f\n", symbol, price)
	}
}
