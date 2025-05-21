package coinapi

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRateService interface {
	GetLatestRate(base, quote string) (*ExchangeRate, error)
}

// ExchangeRate is for internal use in the application
type ExchangeRate struct {
	Base      string          `json:"base"`
	Quote     string          `json:"quote"`
	Rate      decimal.Decimal `json:"rate"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// exchangeRateEvent is the type of events sent by CoinAPI
type exchangeRateEvent struct {
	Time         time.Time
	Rate         decimal.Decimal
	AssetIDBase  string `json:"asset_id_base"`
	AssetIDQuote string `json:"asset_id_quote"`
}
