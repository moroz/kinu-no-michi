package coinapi

import (
	"database/sql"

	"github.com/shopspring/decimal"
	_ "modernc.org/sqlite"
)

type coinapiRestClient struct {
	token string
	db    *sql.DB
}

const initSQL = `
create table exchange_rates (base text, quote text, rate text, updated_at_ms bigint);
`

func NewCoinAPIRESTClient(token string) (*coinapiRestClient, error) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(initSQL)
	if err != nil {
		return nil, err
	}

	return &coinapiRestClient{
		token: token,
		db:    conn,
	}, nil
}

func (c *coinapiRestClient) GetLatestRate() *decimal.Decimal {
	value := decimal.NewFromInt(42)
	return &value
}
