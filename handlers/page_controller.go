package handlers

import (
	"net/http"

	"github.com/moroz/kinu-no-michi/templates"
)

type pageController struct{}

func PageController() pageController {
	return pageController{}
}

func (c *pageController) Index(w http.ResponseWriter, r *http.Request) {
	templates.PagesIndex().Render(r.Context(), w)
}
