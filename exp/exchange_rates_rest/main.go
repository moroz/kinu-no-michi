package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/shopspring/decimal"
)

var COINAPI_API_KEY = config.MustGetenv("COINAPI_API_KEY")

const baseURL = "https://api-realtime.exrates.coinapi.io/v1/exchangerate"

func fetchExchangRate(base, quote string) (*http.Response, error) {
	url := fmt.Sprintf(`%s/%s/%s`, baseURL, base, quote)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", COINAPI_API_KEY)

	return http.DefaultClient.Do(req)
}

func main() {
	resp, err := fetchExchangRate("BTC", "EUR")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var event coinapi.ExchangeRateEvent
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		log.Fatal(err)
	}

	price := decimal.NewFromInt(100)
	converted := price.Div(event.Rate).Round(8)
	fmt.Printf("100 EUR = %s BTC\n", converted)
}
