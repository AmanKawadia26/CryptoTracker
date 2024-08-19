package crypto

import (
	"cryptotracker/models"
	"cryptotracker/pkg/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
)

const baseURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency"

func getAPIResponse(endpoint string, params map[string]string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+endpoint, nil)

	// Add API key to the request header
	req.Header.Add("X-CMC_PRO_API_KEY", config.AppConfig.APIKey)

	// Add parameters to the request URL
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making API request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading API response: %v", err)
	}

	return body
}

func DisplayTopCryptocurrencies() {
	params := map[string]string{
		"start":   "1",
		"limit":   "10",
		"convert": "USD",
	}

	response := getAPIResponse("/listings/latest", params)

	var result map[string]interface{}
	json.Unmarshal(response, &result)

	data := result["data"].([]interface{})

	color.New(color.FgGreen).Println("Top 10 Cryptocurrencies:")
	for i, crypto := range data {
		cryptoMap := crypto.(map[string]interface{})
		name := cryptoMap["name"].(string)
		symbol := cryptoMap["symbol"].(string)
		price := cryptoMap["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"].(float64)

		fmt.Printf("%d. %s (%s): $%.2f\n", i+1, name, symbol, price)
	}
	fmt.Println()
}

func SearchCryptocurrency() {
	var symbol string
	color.New(color.FgCyan).Print("Enter the symbol of the cryptocurrency: ")
	fmt.Scan(&symbol)

	params := map[string]string{
		"symbol":  symbol,
		"convert": "USD",
	}

	response := getAPIResponse("/quotes/latest", params)

	var result map[string]interface{}
	err := json.Unmarshal(response, &result)
	if err != nil {
		color.New(color.FgRed).Printf("Error unmarshalling API response: %v\n", err)
		return
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		color.New(color.FgRed).Println("Data not found in the response.")
		return
	}

	cryptoData, ok := data[symbol].(map[string]interface{})
	if !ok {
		color.New(color.FgYellow).Printf("Cryptocurrency data not found for symbol: %s\n", symbol)
		color.New(color.FgMagenta).Println("Please request the addition of this cryptocurrency to our app.")

		// Create a new request for the unavailable cryptocurrency
		request := &models.UnavailableCryptoRequest{
			CryptoSymbol:   symbol,
			RequestMessage: "Please add this cryptocurrency.",
			Status:         "Pending",
		}

		// Save the request to the file
		if err := saveUnavailableCryptoRequest(request); err != nil {
			color.New(color.FgRed).Printf("Error saving unavailable crypto request: %v\n", err)
			return
		}
		color.New(color.FgGreen).Println("Request to add the cryptocurrency has been submitted.")
		return
	}

	// Proceed with normal processing if the cryptocurrency is found
	name, ok := cryptoData["name"].(string)
	if !ok {
		color.New(color.FgRed).Println("Could not retrieve the name of the cryptocurrency.")
		return
	}

	priceObj, ok := cryptoData["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"]
	if !ok {
		color.New(color.FgRed).Println("Could not retrieve the price of the cryptocurrency.")
		return
	}

	price, ok := priceObj.(float64)
	if !ok {
		color.New(color.FgRed).Println("Failed to convert price to float64.")
		return
	}

	color.New(color.FgGreen).Printf("%s (%s): $%.2f\n", name, symbol, price)
	fmt.Println()
}

func SetPriceAlert(user *models.User) {
	var symbol string
	var targetPrice float64

	color.New(color.FgCyan).Print("Enter the symbol of the cryptocurrency: ")
	fmt.Scan(&symbol)
	color.New(color.FgCyan).Print("Enter your target price in USD: ")
	fmt.Scan(&targetPrice)

	params := map[string]string{
		"symbol":  symbol,
		"convert": "USD",
	}

	response := getAPIResponse("/quotes/latest", params)

	var result map[string]interface{}
	json.Unmarshal(response, &result)

	data := result["data"].(map[string]interface{})[symbol].(map[string]interface{})
	price := data["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"].(float64)

	if price >= targetPrice {
		color.New(color.FgGreen).Printf("Alert: %s has reached your target price of $%.2f. Current price: $%.2f\n", symbol, targetPrice, price)
	} else {
		color.New(color.FgYellow).Printf("%s is still below your target price. Current price: $%.2f\n", symbol, price)
	}
}

func saveUnavailableCryptoRequest(request *models.UnavailableCryptoRequest) error {
	requests, err := loadUnavailableCryptoRequests()
	if err != nil {
		return err
	}

	requests = append(requests, request)
	data, err := json.Marshal(requests)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("unavailable_cryptos.json", data, 0644); err != nil {
		return err
	}

	return nil
}

func loadUnavailableCryptoRequests() ([]*models.UnavailableCryptoRequest, error) {
	if _, err := os.Stat("unavailable_cryptos.json"); os.IsNotExist(err) {
		return []*models.UnavailableCryptoRequest{}, nil
	}

	data, err := ioutil.ReadFile("unavailable_cryptos.json")
	if err != nil {
		return nil, err
	}

	var requests []*models.UnavailableCryptoRequest
	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}
