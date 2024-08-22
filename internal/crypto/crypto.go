package crypto

import (
	"cryptotracker/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
)

const baseURL = "http://api.coinlayer.com/"

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
