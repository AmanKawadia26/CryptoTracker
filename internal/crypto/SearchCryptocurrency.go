package crypto

import (
	"cryptotracker/models"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	apiKey = "885b7aa609bbe129f6df31ea09a03a4e"
)

func SearchCryptocurrency() {
	var input string
	color.New(color.FgCyan).Print("Enter the symbol or name of the cryptocurrency: ")
	fmt.Scan(&input)

	// Normalize the input to lowercase for case-insensitive comparison
	input = strings.ToLower(input)

	// Make an API call to get the list of cryptocurrencies
	cryptoList := getCryptoList()

	for symbol, name := range cryptoList {
		if strings.ToLower(symbol) == input || strings.ToLower(name) == input {
			// Cryptocurrency found, fetch current price
			currentPrice := getCurrentPrice(symbol)

			// Display the result
			color.New(color.FgGreen).Printf("%s (%s): $%.2f\n", name, symbol, currentPrice)
			fmt.Println()

			// Fetch and display the graph
			displayCryptoGraph(symbol, name)
			return
		}
	}

	// If no match is found
	color.New(color.FgYellow).Printf("Cryptocurrency not found for input: %s\n", input)
	color.New(color.FgMagenta).Println("Please request the addition of this cryptocurrency to our app.")

	// Create a new request for the unavailable cryptocurrency
	request := &models.UnavailableCryptoRequest{
		CryptoSymbol:   input,
		RequestMessage: "Please add this cryptocurrency.",
		Status:         "Pending",
	}

	// Save the request to the file
	if err := saveUnavailableCryptoRequest(request); err != nil {
		color.New(color.FgRed).Printf("Error saving unavailable crypto request: %v\n", err)
		return
	}
	color.New(color.FgGreen).Println("Request to add the cryptocurrency has been submitted.")
}

func getCryptoList() map[string]string {
	url := fmt.Sprintf("%slist?access_key=%s", baseURL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		color.New(color.FgRed).Printf("Error fetching crypto list: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.New(color.FgRed).Printf("Error reading response body: %v\n", err)
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		color.New(color.FgRed).Printf("Error unmarshalling API response: %v\n", err)
		return nil
	}

	cryptoList := make(map[string]string)
	for symbol, data := range result["crypto"].(map[string]interface{}) {
		cryptoData := data.(map[string]interface{})
		cryptoList[symbol] = cryptoData["name"].(string)
	}

	return cryptoList
}

func getCurrentPrice(symbol string) float64 {
	url := fmt.Sprintf("%slive?access_key=%s&symbols=%s", baseURL, apiKey, symbol)
	resp, err := http.Get(url)
	if err != nil {
		color.New(color.FgRed).Printf("Error fetching current price: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.New(color.FgRed).Printf("Error reading response body: %v\n", err)
		return 0
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		color.New(color.FgRed).Printf("Error unmarshalling API response: %v\n", err)
		return 0
	}

	rates := result["rates"].(map[string]interface{})
	return rates[symbol].(float64)
}

func displayCryptoGraph(symbol, name string) {
	prices := getHistoricalPrices(symbol)

	color.New(color.FgCyan).Printf("30-day price graph for %s (%s):\n\n", name, symbol)

	maxPrice := prices[0]
	minPrice := prices[0]
	for _, price := range prices {
		if price > maxPrice {
			maxPrice = price
		}
		if price < minPrice {
			minPrice = price
		}
	}

	graphHeight := 20
	for i := 0; i < graphHeight; i++ {
		price := maxPrice - (float64(i) * (maxPrice - minPrice) / float64(graphHeight-1))
		fmt.Printf("%8.2f |", price)

		for _, p := range prices {
			if p >= price {
				color.New(color.FgGreen).Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Print("         ")
	fmt.Println(strings.Repeat("-", len(prices)))
	fmt.Print("         ")
	for i := 0; i < len(prices); i += len(prices) / 5 {
		fmt.Printf("%-6d", 30-i)
	}
	fmt.Println("\n         Days ago")
}

func getHistoricalPrices(symbol string) []float64 {
	prices := make([]float64, 30)
	today := time.Now()

	for i := 29; i >= 0; i-- {
		date := today.AddDate(0, 0, -i)
		url := fmt.Sprintf("%s%s?access_key=%s&symbols=%s", baseURL, date.Format("2006-01-02"), apiKey, symbol)
		resp, err := http.Get(url)
		if err != nil {
			color.New(color.FgRed).Printf("Error fetching historical price for %s: %v\n", date.Format("2006-01-02"), err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			color.New(color.FgRed).Printf("Error reading response body: %v\n", err)
			continue
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			color.New(color.FgRed).Printf("Error unmarshalling API response: %v\n", err)
			continue
		}

		rates := result["rates"].(map[string]interface{})
		prices[i] = rates[symbol].(float64)
	}

	return prices
}
