package handlers

import (
	"net/http"

	"github.com/moroz/kinu-no-michi/templates"
)

type pageController struct {
	rs ExchangeRateService
}

func PageController(rs ExchangeRateService) pageController {
	return pageController{rs}
}

func (c *pageController) Index(w http.ResponseWriter, r *http.Request) {
	rate := c.rs.GetLatestRate()
	templates.PagesIndex(rate).Render(r.Context(), w)
}
