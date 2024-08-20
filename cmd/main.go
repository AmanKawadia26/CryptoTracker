package main

import (
	"cryptotracker/internal/admin"
	"cryptotracker/internal/auth"
	"cryptotracker/pkg/config"
	"cryptotracker/pkg/ui"
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
