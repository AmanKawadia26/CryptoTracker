package crypto

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"sort"
)

func DisplayTopCryptocurrencies() {
	url := fmt.Sprintf("%slive?access_key=%s", baseURL, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		color.New(color.FgRed).Printf("Error fetching data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.New(color.FgRed).Printf("Error reading response body: %v\n", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		color.New(color.FgRed).Printf("Error unmarshalling API response: %v\n", err)
		return
	}

	rates, ok := result["rates"].(map[string]interface{})
	if !ok {
		color.New(color.FgRed).Println("Error: Unexpected data structure in API response")
		return
	}

	// Convert rates to a slice for sorting
	type cryptoData struct {
		Symbol string
		Price  float64
	}
	var cryptoList []cryptoData

	for symbol, price := range rates {
		priceFloat, ok := price.(float64)
		if !ok {
			color.New(color.FgRed).Printf("Error: Unable to parse price for %s\n", symbol)
			continue
		}
		cryptoList = append(cryptoList, cryptoData{Symbol: symbol, Price: priceFloat})
	}

	// Sort cryptocurrencies by price (descending order)
	sort.Slice(cryptoList, func(i, j int) bool {
		return cryptoList[i].Price > cryptoList[j].Price
	})

	color.New(color.FgGreen).Println("Top 10 Cryptocurrencies:")
	for i, crypto := range cryptoList[:10] {
		fmt.Printf("%d. %s: $%.2f\n", i+1, crypto.Symbol, crypto.Price)
	}
	fmt.Println()
}
