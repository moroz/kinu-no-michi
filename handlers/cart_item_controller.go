package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/cookies"
	"github.com/moroz/kinu-no-michi/services"
	"github.com/shopspring/decimal"
)

type cartItemController struct {
	cs  cookies.SessionStore
	srv *services.CartService
}

func CartItemController(db queries.DBTX, cs cookies.SessionStore) *cartItemController {
	return &cartItemController{
		cs,
		services.NewCartService(db),
	}
}

func parseParams(r *http.Request) (*services.AddProductToCartParams, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	productId, err := uuid.Parse(r.FormValue("product_id"))
	if err != nil {
		return nil, err
	}

	quantity, err := decimal.NewFromString(r.FormValue("quantity"))
	if err != nil {
		return nil, err
	}

	return &services.AddProductToCartParams{
		ProductID: productId,
		Quantity:  quantity,
	}, nil
}

func (c *cartItemController) Create(w http.ResponseWriter, r *http.Request) {
	params, err := parseParams(r)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	item, err := c.srv.AddProductToCart(r.Context(), params)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	session := encodeSession(&appSession{
		CartID: item.CartID,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     config.SESSION_COOKIE_NAME,
		Value:    c.cs.Encode(session),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(204)
}
