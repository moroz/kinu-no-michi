package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/lib/encrypt"
)

type OrderService struct {
	db queries.DBTX
	rs coinapi.ExchangeRateService
}

func NewOrderService(db queries.DBTX, rs coinapi.ExchangeRateService) *OrderService {
	return &OrderService{db, rs}
}

type CreateOrderParams struct {
	CartID uuid.UUID
	Email  string
}

func (s *OrderService) CreateOrder(ctx context.Context, params *CreateOrderParams) (*queries.Order, error) {
	db, ok := s.db.(*pgxpool.Pool)
	if !ok {
		return nil, errors.New("failed to check out database connection")
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rate, err := s.rs.GetLatestRate("BTC", "EUR")
	if err != nil {
		return nil, err
	}

	cart, err := queries.New(tx).GetCartByID(ctx, params.CartID)
	if err != nil {
		return nil, err
	}

	order, err := queries.New(tx).InsertOrder(ctx, queries.InsertOrderParams{
		ID:             uuid.Must(uuid.NewV7()),
		EmailEncrypted: encrypt.NewEncryptedString(params.Email),
		ExchangeRate:   rate.Rate,
		GrandTotalEur:  cart.GrandTotal,
		GrandTotalBtc:  cart.GrandTotal.Div(rate.Rate),
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}
