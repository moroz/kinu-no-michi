package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/lib/cookies"
)

func Router(db queries.DBTX, rs coinapi.ExchangeRateService, cs cookies.SessionStore) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Use(FetchSession(cs))
	r.Use(LoadCartFromSession(db))

	pages := PageController(db, rs)
	r.Get("/", pages.Index)

	carts := CartController(db, rs)
	r.Get("/cart", carts.Show)

	cartItems := CartItemController(db, cs)
	r.Post("/cart_items", cartItems.Create)

	fs := http.Dir("./assets/dist")
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(fs)))

	return r
}
