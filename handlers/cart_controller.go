package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/templates"
)

type cartController struct {
	db queries.DBTX
	rs coinapi.ExchangeRateService
}

func CartController(db queries.DBTX, rs coinapi.ExchangeRateService) *cartController {
	return &cartController{db, rs}
}

func (c *cartController) Show(w http.ResponseWriter, r *http.Request) {
	var items []*queries.GetCartItemsByCartIDRow
	if cart, ok := r.Context().Value("cart").(*queries.GetCartByIDRow); ok && cart != nil && cart.ID != uuid.Nil {
		i, err := queries.New(c.db).GetCartItemsByCartID(r.Context(), cart.ID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		items = i
	}

	rate, err := c.rs.GetLatestRate("BTC", "EUR")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	templates.CartShow(&templates.CartShowProps{
		CartItems: items,
		Rate:      rate,
	}).Render(r.Context(), w)
}
