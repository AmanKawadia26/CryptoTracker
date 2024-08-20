package auth

import (
	"cryptotracker/models"
	"cryptotracker/pkg/ui"
	"fmt"
	"github.com/fatih/color"
)

// AuthenticateUser handles the login/signup process
func AuthenticateUser() (*models.User, bool) {
	for {
		ui.ClearScreen()
		ui.DisplayAuthMenu()

		var choice int
		color.New(color.FgCyan).Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if user, isAdmin, err := Login(); err == nil {
				return user, isAdmin
			} else {
				color.New(color.FgRed).Println("Login failed:", err)
			}
		case 2:
			if err := Signup(); err != nil {
				color.New(color.FgRed).Println("Signup failed:", err)
			} else {
				color.New(color.FgGreen).Println("Signup successful. Please login.")
			}
		case 3:
			color.New(color.FgYellow).Println("Exiting...")
			return nil, false
		default:
			color.New(color.FgRed).Println("Invalid choice, please try again.")
		}
	}
}
