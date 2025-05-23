package handlers

import (
	"net/http"

	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/templates"
)

type pageController struct {
	rs coinapi.ExchangeRateService
	db queries.DBTX
}

func PageController(db queries.DBTX, rs coinapi.ExchangeRateService) pageController {
	return pageController{rs, db}
}

func (c *pageController) Index(w http.ResponseWriter, r *http.Request) {
	products, err := queries.New(c.db).ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	rate, err := c.rs.GetLatestRate("BTC", "EUR")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	templates.PagesIndex(&templates.ProductsIndexProps{
		Products:     products,
		ExchangeRate: rate,
	}).Render(r.Context(), w)
}
