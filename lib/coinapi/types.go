package coinapi

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRateService interface {
	GetLatestRate(base, quote string) (*ExchangeRate, error)
}

type ExchangeRate struct {
	Base      string
	Quote     string
	Rate      decimal.Decimal
	UpdatedAt time.Time
}

type ExchangeRateEvent struct {
	Time         time.Time
	Rate         decimal.Decimal
	AssetIDBase  string `json:"asset_id_base"`
	AssetIDQuote string `json:"asset_id_quote"`
}
