package auth

import (
	"cryptotracker/internal/crypto"
	"cryptotracker/models"
	"cryptotracker/pkg/storage"
	"cryptotracker/pkg/ui"
	"cryptotracker/pkg/utils"
	"errors"
	"fmt"
	"github.com/fatih/color"
)

// Login handles the login process
func Login() (*models.User, string, error) {
	var username, password string

	color.New(color.FgCyan).Print("Enter username: ")
	fmt.Scan(&username)
	password = ui.GetHiddenInput("Enter password: ")

	user, err := storage.GetUserByUsername(username)
	if err != nil {
		return nil, "", err
	}

	hashedPassword := utils.HashPassword(password)
	if user.Password != hashedPassword {
		return nil, "", errors.New("invalid username or password")
	}

	color.New(color.FgGreen).Println("Login successful.")

	// Check if any price alerts have been met for the user
	notifications, err := crypto.LoadPriceNotifications()
	if err == nil {
		for _, notification := range notifications {
			if notification.Username == user.Username && notification.Status == "Served" {
				color.New(color.FgGreen).Printf("Price alert met: %s has reached $%.2f.\n", notification.Crypto, notification.TargetPrice)
			}
		}
	}

	return user, user.Role, nil
}
