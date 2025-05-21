package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/moroz/kinu-no-michi/config"
	"github.com/shopspring/decimal"
)

var COINAPI_API_KEY = config.MustGetenv("COINAPI_API_KEY")

const baseURL = "https://api-realtime.exrates.coinapi.io/v1/exchangerate"

type exchangeRateEvent struct {
	Time         time.Time
	Rate         decimal.Decimal
	AssetIDBase  string `json:"asset_id_base"`
	AssetIDQuote string `json:"asset_id_quote"`
}

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

	var event exchangeRateEvent
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		log.Fatal(err)
	}

	price := decimal.NewFromInt(100)
	converted := price.Div(event.Rate).Round(8)
	fmt.Printf("100 EUR = %s BTC\n", converted)
}
