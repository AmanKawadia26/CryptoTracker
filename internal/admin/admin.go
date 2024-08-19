package admin

import (
	"cryptotracker/pkg/storage"
	"fmt"

	"github.com/fatih/color"
)

func ShowAdminPanel() {
	for {
		color.New(color.FgGreen).Println("Admin Panel")
		fmt.Println("1. Manage Users")
		fmt.Println("2. View User Profiles")
		fmt.Println("3. Manage User Requests")
		fmt.Println("4. Logout")

		var choice int
		color.New(color.FgYellow).Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			ManageUsers()
		case 2:
			ViewUserProfiles()
		case 3:
			ManageUserRequests()
		case 4:
			color.New(color.FgCyan).Println("Logging out...")
			return
		default:
			color.New(color.FgRed).Println("Invalid choice, please try again.")
		}
	}
}

func ManageUsers() {
	// Manage users logic here
	fmt.Println()
	color.New(color.FgGreen).Println("Managing users...")
	fmt.Println()
}

func ManageUserRequests() {
	// Add handling for unavailable crypto requests
	unavailableRequests, err := storage.GetAllUnavailableCryptoRequests()
	if err != nil {
		color.New(color.FgRed).Println("Error fetching unavailable crypto requests:", err)
		return
	}

	if len(unavailableRequests) == 0 {
		color.New(color.FgCyan).Println("No pending unavailable crypto requests.")
		return
	}

	for i, request := range unavailableRequests {
		color.New(color.FgYellow).Printf("[%d] - Symbol: %s, User: %s, Message: %s, Status: %s\n", i+1, request.CryptoSymbol, request.UserID, request.RequestMessage, request.Status)
	}

	var choice int
	color.New(color.FgYellow).Print("Select a request to manage: ")
	fmt.Scan(&choice)

	if choice < 1 || choice > len(unavailableRequests) {
		color.New(color.FgRed).Println("Invalid selection.")
		return
	}

	selectedRequest := unavailableRequests[choice-1]
	color.New(color.FgGreen).Println("Selected request:", selectedRequest)

	// Allow admin to approve or reject the request
	var action string
	color.New(color.FgYellow).Print("Enter 'approve' to approve the request or 'reject' to reject it: ")
	fmt.Scan(&action)

	if action == "approve" {
		selectedRequest.Status = "Approved"
	} else if action == "reject" {
		selectedRequest.Status = "Rejected"
	} else {
		color.New(color.FgRed).Println("Invalid action.")
		return
	}

	if err := storage.SaveUnavailableCryptoRequest(selectedRequest); err != nil {
		color.New(color.FgRed).Println("Error updating request status:", err)
	} else {
		color.New(color.FgGreen).Println("Request status updated.")
	}
}

func ViewUserProfiles() {
	users, err := storage.GetAllUsers()
	if err != nil {
		color.New(color.FgRed).Println("Error fetching user profiles:", err)
		return
	}

	if len(users) == 0 {
		color.New(color.FgCyan).Println("No users found.")
		return
	}

	for i, user := range users {
		color.New(color.FgYellow).Printf("[%d] - Username: %s, Email: %s, Mobile: %s, Is Admin: %t\n", i+1, user.Username, user.Email, user.Mobile, user.IsAdmin)
	}
}
