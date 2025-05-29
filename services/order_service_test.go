package services_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/lib/encrypt"
	"github.com/moroz/kinu-no-michi/services"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	db, err := pgxpool.New(context.Background(), config.MustGetenv("TEST_DATABASE_URL"))
	require.NoError(t, err)

	_, err = db.Exec(context.Background(), "delete from orders")
	require.NoError(t, err)

	var key, _ = base64.StdEncoding.DecodeString("mAcJgVR/M5/OejAQNtbTeIeb+O7AoLuZ2purSckAKuM=")
	provider, err := encrypt.NewXChacha20Provider(key)
	require.NoError(t, err)
	encrypt.SetProvider(provider)

	// pre-insert some cart items
	cs := services.NewCartService(db)
	item, err := cs.AddProductToCart(context.Background(), &services.AddProductToCartParams{
		CartID:    nil,
		ProductID: uuid.MustParse("019709a2-5c37-73e2-a05b-9ee9f8a470b5"),
		Quantity:  decimal.NewFromInt(2),
	})
	require.NoError(t, err)
	_, err = cs.AddProductToCart(context.Background(), &services.AddProductToCartParams{
		CartID:    &item.CartID,
		ProductID: uuid.MustParse("019709a2-5c3a-7c63-b811-090eaedf0835"),
		Quantity:  decimal.NewFromInt(3),
	})
	require.NoError(t, err)

	rs := coinapi.NewMockClient(100000)

	srv := services.NewOrderService(db, rs)

	order, err := srv.CreateOrder(context.Background(), &services.CreateOrderParams{
		CartID: item.CartID,
		Email:  "user@example.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, order)

	assert.Equal(t, decimal.NewFromInt(100000), order.ExchangeRate)
	assert.Equal(t, decimal.NewFromInt(280), order.GrandTotalEur)
	expected, _ := decimal.NewFromString("0.0028")
	assert.Equal(t, expected, order.GrandTotalBtc)
}
