package auth

import (
	"crypto/sha256"
	"cryptotracker/models"
	"cryptotracker/pkg/storage"
	"cryptotracker/pkg/ui"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"unicode"

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

	hashedPassword := hashPassword(password)
	if user.Password != hashedPassword {
		return nil, false, errors.New("invalid username or password")
	}

	color.New(color.FgGreen).Println("Login successful.")
	return user, user.IsAdmin, nil
}

// Signup handles the signup process
func Signup() error {
	var username, password, email, pan string
	var mobile int

	color.New(color.FgCyan).Print("Enter username: ")
	fmt.Scan(&username)
	if !isValidUsername(username) {
		return errors.New("invalid username: must be one word, alphanumeric, and can contain underscores")
	}

	password = ui.GetHiddenInput("Enter password: ")
	if !isValidPassword(password) {
		return errors.New("invalid password: must be at least 8 characters, include an uppercase letter, a number, and a special character")
	}

	color.New(color.FgCyan).Print("Enter email: ")
	fmt.Scan(&email)
	if !isValidEmail(email) {
		return errors.New("invalid email: must be a valid email address")
	}

	color.New(color.FgCyan).Print("Enter mobile (10 digits): ")
	fmt.Scan(&mobile)
	if !isValidMobile(mobile) {
		return errors.New("invalid mobile number: must be 10 digits")
	}

	color.New(color.FgCyan).Print("Enter PAN: ")
	fmt.Scan(&pan)

	hashedPassword := hashPassword(password)

	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Mobile:   mobile,
		IsAdmin:  false,
	}

	if err := storage.SaveUser(user); err != nil {
		return err
	}

	color.New(color.FgGreen).Println("Signup successful.")
	return nil
}

// Utility functions

// isValidUsername validates the username format
func isValidUsername(username string) bool {
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validUsername.MatchString(username)
}

// isValidPassword validates the password strength
func isValidPassword(password string) bool {
	var (
		hasMinLen      = len(password) >= 8
		hasUpper       = false
		hasLower       = false
		hasNumber      = false
		hasSpecialChar = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecialChar
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	validEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return validEmail.MatchString(email)
}

// isValidMobile validates the mobile number format
func isValidMobile(mobile int) bool {
	return len(fmt.Sprintf("%d", mobile)) == 10
}

// hashPassword hashes the password using SHA-256
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
