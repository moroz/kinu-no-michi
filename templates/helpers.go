package templates

import "github.com/shopspring/decimal"

type menuItem struct {
	label string
	href  string
}

var menuItems = []menuItem{
	{"Home", "/"},
	{"Fruit", "/fruit"},
	{"Vegetables", "/vegetables"},
	{"Books", "/books"},
	{"Firearms", "/firearms"},
	{"Cosmetics", "/cosmetics"},
}

func fiatToBTC(fiat, rate decimal.Decimal) string {
	return fiat.Div(rate).Round(8).String()
}
