package handlers

import (
	"net/http"

	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/templates"
)

type pageController struct {
	db queries.DBTX
	rs coinapi.ExchangeRateService
}

func PageController(db queries.DBTX, rs coinapi.ExchangeRateService) pageController {
	return pageController{db, rs}
}

func (c *pageController) Index(w http.ResponseWriter, r *http.Request) {
	products, err := queries.New(c.db).ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rate, err := c.rs.GetLatestRate("BTC", "EUR")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	templates.PagesIndex(&templates.PagesIndexProps{
		Products: products,
		Rate:     rate,
	}).Render(r.Context(), w)
}
