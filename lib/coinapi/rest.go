package coinapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type coinapiRestClient struct {
	token  string
	db     *sql.DB
	maxAge int
}

const initSQL = `
create table if not exists exchange_rates (
	base text, quote text, rate text, updated_at_ms bigint,
	primary key (base, quote)
);
`

const DEFAULT_MAX_AGE = 5 * 60 * 1000

func NewCoinAPIRESTClient(token string, maxAge int) (*coinapiRestClient, error) {
	conn, err := sql.Open("sqlite", "cache.db")
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(initSQL)
	if err != nil {
		return nil, err
	}

	client := coinapiRestClient{
		token:  token,
		db:     conn,
		maxAge: maxAge,
	}

	if maxAge <= 0 {
		client.maxAge = DEFAULT_MAX_AGE
	}

	return &client, err
}

func (c *coinapiRestClient) GetLatestRate(base, quote string) (*ExchangeRate, error) {
	cached, err := c.getCachedRate(base, quote)

	threshold := time.Now().Add(-1 * time.Duration(c.maxAge) * time.Millisecond)

	// check if the quote is present in the database and if it's newer than maxAge ms
	if err == nil && cached.UpdatedAt.After(threshold) {
		return cached, nil
	}

	return c.refreshRate(base, quote)
}

func (c *coinapiRestClient) getCachedRate(base, quote string) (*ExchangeRate, error) {
	var rate ExchangeRate
	var ts int64
	err := c.db.QueryRow(
		"select base, quote, rate, updated_at_ms from exchange_rates where base = ? and quote = ?", base, quote,
	).Scan(&rate.Base, &rate.Quote, &rate.Rate, &ts)
	if err != nil {
		return nil, err
	}
	rate.UpdatedAt = time.Unix(0, ts*1000000)
	return &rate, nil
}

const baseURL = "https://api-realtime.exrates.coinapi.io/v1/exchangerate"

func (c *coinapiRestClient) refreshRate(base, quote string) (*ExchangeRate, error) {
	log.Printf("Refreshing rates for %s/%s", base, quote)
	url := fmt.Sprintf(`%s/%s/%s`, baseURL, base, quote)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var event exchangeRateEvent
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		return nil, err
	}

	_, err = c.db.Exec(`insert into exchange_rates (base, quote, rate, updated_at_ms) values (?, ?, ?, ?) on conflict (base, quote) do update set rate = excluded.rate, updated_at_ms = excluded.updated_at_ms`, event.AssetIDBase, event.AssetIDQuote, event.Rate.String(), event.Time.UnixMilli())

	if err != nil {
		return nil, err
	}

	return &ExchangeRate{
		Base:      event.AssetIDBase,
		Quote:     event.AssetIDQuote,
		Rate:      event.Rate,
		UpdatedAt: event.Time,
	}, nil
}
