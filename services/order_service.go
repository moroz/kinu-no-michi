package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

	cart, err := queries.New(tx).GetCartByID(ctx, params.CartID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("order does not exist: %w", err)
	}
	if cart.ItemCount == 0 {
		return nil, fmt.Errorf("no items in cart")
	}

	rate, err := s.rs.GetLatestRate("BTC", "EUR")
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	order, err := queries.New(tx).InsertOrder(ctx, queries.InsertOrderParams{
		ID:             id,
		EmailEncrypted: encrypt.NewEncryptedString(params.Email),
		GrandTotalEur:  cart.GrandTotal,
		GrandTotalBtc:  cart.GrandTotal.DivRound(rate.Rate, 8),
		ExchangeRate:   rate.Rate,
	})
	if err != nil {
		return nil, err
	}

	tx.Commit(ctx)
	return order, nil
}
