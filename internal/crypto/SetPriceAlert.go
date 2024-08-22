package crypto

import (
	"cryptotracker/models"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"time"
)

func SetPriceAlert(user *models.User) {
	var symbol string
	var targetPrice float64

	color.New(color.FgCyan).Print("Enter the symbol of the cryptocurrency: ")
	fmt.Scan(&symbol)
	color.New(color.FgCyan).Print("Enter your target price in USD: ")
	fmt.Scan(&targetPrice)

	url := fmt.Sprintf("%slive?access_key=%s&symbols=%s", baseURL, apiKey, symbol)

	resp, err := http.Get(url)
	if err != nil {
		color.New(color.FgRed).Printf("Error fetching data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.New(color.FgRed).Printf("Error reading response body: %v\n", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		color.New(color.FgRed).Printf("Error unmarshalling API response: %v\n", err)
		return
	}

	rates, ok := result["rates"].(map[string]interface{})
	if !ok {
		color.New(color.FgRed).Println("Error: Unexpected data structure in API response")
		return
	}

	currentPrice, ok := rates[symbol].(float64)
	if !ok {
		color.New(color.FgRed).Printf("Error: Price data not available for %s\n", symbol)
		return
	}

	// Check if the target price has been met
	if currentPrice >= targetPrice {
		color.New(color.FgGreen).Printf("Alert: %s has reached your target price of $%.2f. Current price: $%.2f\n", symbol, targetPrice, currentPrice)
	} else {
		// Create a notification request if the target price is not met
		notification := &models.PriceNotification{
			Crypto:      symbol,
			TargetPrice: targetPrice,
			Username:    user.Username,
			AskedAt:     time.Now().Format(time.RFC3339),
			Status:      "Pending",
		}

		// Save the notification
		err = savePriceNotification(notification)
		if err != nil {
			color.New(color.FgRed).Printf("Error saving notification: %v\n", err)
			return
		}

		color.New(color.FgYellow).Printf("%s is still below your target price. Current price: $%.2f. Notification created.\n", symbol, currentPrice)
	}
}
