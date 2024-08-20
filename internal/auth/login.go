package auth

import (
	"cryptotracker/models"
	"cryptotracker/pkg/storage"
	"cryptotracker/pkg/ui"
	"cryptotracker/pkg/utils"
	"errors"
	"fmt"
	"github.com/fatih/color"
)

// Login handles the login process
func Login() (*models.User, bool, error) {
	var username, password string

	color.New(color.FgCyan).Print("Enter username: ")
	fmt.Scan(&username)
	password = ui.GetHiddenInput("Enter password: ")

	user, err := storage.GetUserByUsername(username)
	if err != nil {
		return nil, false, err
	}

	hashedPassword := utils.HashPassword(password)
	if user.Password != hashedPassword {
		return nil, false, errors.New("invalid username or password")
	}

	color.New(color.FgGreen).Println("Login successful.")
	return user, user.IsAdmin, nil
}
