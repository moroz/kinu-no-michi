package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/encrypt"
)

type OrderService struct {
	db queries.DBTX
}

func NewOrderService(db queries.DBTX) *OrderService {
	return &OrderService{db}
}

func (s *OrderService) CreateOrder(ctx context.Context, email string) (*queries.Order, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return queries.New(s.db).InsertOrder(ctx, queries.InsertOrderParams{
		ID:             id,
		EmailEncrypted: encrypt.NewEncryptedString(email),
	})
}
