package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// Define color variables using Fatih's color package
var (
	colorRed    = color.New(color.FgRed).SprintFunc()
	colorGreen  = color.New(color.FgGreen).SprintFunc()
	colorYellow = color.New(color.FgYellow).SprintFunc()
	colorBlue   = color.New(color.FgBlue).SprintFunc()
	colorPurple = color.New(color.FgMagenta).SprintFunc()
	colorCyan   = color.New(color.FgCyan).SprintFunc()
)

// ClearScreen clears the terminal screen
func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// DisplayWelcomeBanner shows the welcome message
func DisplayWelcomeBanner() {
	ClearScreen()
	fmt.Println(colorCyan("=============================="))
	fmt.Println(colorCyan("  Welcome to CryptoTracker! 🚀"))
	fmt.Println(colorCyan("=============================="))
	fmt.Println()
}

// DisplayAuthMenu shows the authentication menu
func DisplayAuthMenu() {
	fmt.Println(colorYellow("Authentication Menu:"))
	fmt.Println(colorBlue("1. Login"))
	fmt.Println(colorBlue("2. SignUp"))
	fmt.Println(colorBlue("3. Exit"))
	fmt.Println()
}

// DisplayMainMenu shows the main menu for regular users
func DisplayMainMenu() {
	fmt.Println(colorYellow("Main Menu:"))
	fmt.Println(colorGreen("1. View Top 10 Cryptocurrencies"))
	fmt.Println(colorGreen("2. Search for a Cryptocurrency"))
	fmt.Println(colorGreen("3. Set Price Alert"))
	fmt.Println(colorGreen("4. Check if user is Admin"))
	fmt.Println(colorGreen("5. User Profile"))
	fmt.Println(colorRed("6. Logout"))
	fmt.Println()
}

// GetHiddenInput securely gets user input for passwords
func GetHiddenInput(prompt string) string {
	fmt.Print(colorPurple(prompt))
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		PrintError("Error reading password: " + err.Error())
		return ""
	}
	fmt.Println() // Print a newline after input
	return strings.TrimSpace(string(bytePassword))
}

// PrintError prints an error message in red
func PrintError(message string) {
	fmt.Println(colorRed(message))
}
