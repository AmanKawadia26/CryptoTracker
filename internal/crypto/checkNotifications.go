package crypto

import (
	"cryptotracker/models"
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"time"
)

// Check if any notification requests have met their target price
func checkNotifications() {
	notifications, err := LoadPriceNotifications()
	if err != nil {
		return
	}

	// Iterate over all saved notifications
	for _, notification := range notifications {
		if notification.Status == "Pending" {
			// Check current price for the cryptocurrency
			params := map[string]string{
				"symbol":  notification.Crypto,
				"convert": "USD",
			}

			response := getAPIResponse("/listings/latest", params)

			var result map[string]interface{}
			json.Unmarshal(response, &result)

			data := result["data"].(map[string]interface{})[notification.Crypto].(map[string]interface{})
			price := data["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"].(float64)

			// If the price meets the target, update the notification status
			if price >= notification.TargetPrice {
				notification.Status = "Served"
				notification.ServedAt = time.Now().Format(time.RFC3339)
				color.New(color.FgGreen).Printf("Notification met: %s has reached $%.2f.\n", notification.Crypto, notification.TargetPrice)
			}
		}
	}

	// Save the updated notifications back to the file
	if err := savePriceNotifications(notifications); err != nil {
		color.New(color.FgRed).Printf("Error saving notifications: %v\n", err)
	}
}

// Save the notifications to the JSON file
func savePriceNotifications(notifications []*models.PriceNotification) error {
	data, err := json.Marshal(notifications)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("price_notifications.json", data, 0644); err != nil {
		return err
	}

	return nil
}
