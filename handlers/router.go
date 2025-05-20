package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shopspring/decimal"
)

type ExchangeRateService interface {
	GetLatestRate() *decimal.Decimal
}

func Router(rs ExchangeRateService) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	pages := PageController(rs)
	r.Get("/", pages.Index)

	fs := http.Dir("./assets/dist")
	r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(fs)))

	return r
}
