package templates

import (
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
)

type ProductsIndexProps struct {
	Products     []*queries.Product
	ExchangeRate *coinapi.ExchangeRate
}
