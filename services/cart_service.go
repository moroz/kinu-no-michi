package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/shopspring/decimal"
)

type CartService struct {
	db queries.DBTX
}

func NewCartService(db queries.DBTX) *CartService {
	return &CartService{db}
}

type AddProductToCartParams struct {
	CartID    uuid.UUID
	ProductID uuid.UUID
	Quantity  decimal.Decimal
}

func (s *CartService) AddProductToCart(ctx context.Context, params *AddProductToCartParams) (*queries.CartItem, error) {
	tx, err := s.db.(*pgxpool.Pool).BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if params.CartID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		params.CartID = id

		err = queries.New(s.db).InsertCart(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	item, err := queries.New(s.db).InsertCartItem(ctx, queries.InsertCartItemParams{
		ID:        id,
		CartID:    params.CartID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})
	if err != nil {
		return nil, err
	}

	tx.Commit(ctx)

	return item, nil
}
