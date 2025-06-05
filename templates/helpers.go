package templates

import (
	"context"

	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/shopspring/decimal"
)

type menuItem struct {
	label string
	href  string
}

var menuItems = []menuItem{
	{"Home", "/"},
}

func fiatToBTC(fiat, rate decimal.Decimal) string {
	return fiat.Div(rate).Round(8).String()
}

func getCartFromContext(ctx context.Context) *queries.GetCartByIDRow {
	return ctx.Value("cart").(*queries.GetCartByIDRow)
}
