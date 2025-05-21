package handlers

import (
	"net/http"

	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/templates"
)

type pageController struct {
	rs coinapi.ExchangeRateService
}

func PageController(rs coinapi.ExchangeRateService) pageController {
	return pageController{rs}
}

func (c *pageController) Index(w http.ResponseWriter, r *http.Request) {
	rate, err := c.rs.GetLatestRate("BTC", "EUR")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	templates.PagesIndex(rate).Render(r.Context(), w)
}
