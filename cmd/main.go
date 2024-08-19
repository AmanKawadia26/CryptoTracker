package main

import (
	"cryptotracker/internal/admin"
	"cryptotracker/internal/auth"
	"cryptotracker/internal/crypto"
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"cryptotracker/pkg/storage"
	"cryptotracker/pkg/ui"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

func main() {
	// Load the configuration
	config.LoadConfig()

	// Display welcome banner
	ui.DisplayWelcomeBanner()

	// Start login/signup process
	user, isAdmin := auth.AuthenticateUser()

	// If user is admin, show admin panel
	if isAdmin {
		admin.ShowAdminPanel()
		return
	}

	// Display main user menu
	mainMenu(user)
}

// mainMenu displays the main menu for a regular user
func mainMenu(user *models.User) {
	for {
		ui.ClearScreen()
		ui.DisplayMainMenu()

		var choice int
		color.New(color.FgYellow).Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			crypto.DisplayTopCryptocurrencies()
		case 2:
			crypto.SearchCryptocurrency()
		case 3:
			crypto.SetPriceAlert(user)
		case 4:
			if user.IsAdmin {
				admin.ShowAdminPanel()
			} else {
				color.New(color.FgRed).Println("This user is not and admin.\n")
			}
		case 5:
			UserProfile(user.Username)
		case 6:
			color.New(color.FgCyan).Println("Logging out...")
			log.Println("Logging out...")
			os.Exit(0)
		default:
			color.New(color.FgRed).Println("Invalid choice, please try again.")
		}
	}
}

func UserProfile(username string) {
	// Fetch user profile
	user, err := storage.GetUserProfile(username)
	if err != nil {
		color.New(color.FgRed).Println("Error fetching user profile:", err)
		return
	}

	// Display user profile
	fmt.Println()
	color.New(color.FgGreen).Println("User Profile:")
	color.New(color.FgCyan).Printf("Username: %s\n", user.Username)
	color.New(color.FgCyan).Printf("Email: %s\n", user.Email)
	color.New(color.FgCyan).Printf("Mobile: %d\n", user.Mobile)
	color.New(color.FgCyan).Printf("Is Admin: %t\n", user.IsAdmin)
	fmt.Println()
}
