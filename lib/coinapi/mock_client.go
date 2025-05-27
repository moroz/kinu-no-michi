package coinapi

import (
	"time"

	"github.com/shopspring/decimal"
)

type mockClient struct {
	rate decimal.Decimal
}

func NewMockClient(rate int) ExchangeRateService {
	return &mockClient{
		rate: decimal.NewFromInt(int64(rate)),
	}
}

func (c *mockClient) GetLatestRate(base, quote string) (*ExchangeRate, error) {
	return &ExchangeRate{
		Base:      base,
		Quote:     quote,
		Rate:      c.rate,
		UpdatedAt: time.Now(),
	}, nil
}
