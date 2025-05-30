package services_test

import (
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
	db, err := pgxpool.New(t.Context(), config.MustGetenv("TEST_DATABASE_URL"))
	require.NoError(t, err)

	_, err = db.Exec(t.Context(), "truncate orders cascade;")
	require.NoError(t, err)

	var key, _ = base64.StdEncoding.DecodeString("mAcJgVR/M5/OejAQNtbTeIeb+O7AoLuZ2purSckAKuM=")
	provider, err := encrypt.NewXChacha20Provider(key)
	require.NoError(t, err)
	encrypt.SetProvider(provider)

	cs := services.NewCartService(db)
	item, err := cs.AddProductToCart(t.Context(), &services.AddProductToCartParams{
		CartID:    nil,
		ProductID: uuid.MustParse("019709a2-5c37-73e2-a05b-9ee9f8a470b5"),
		Quantity:  decimal.NewFromInt(2),
	})
	assert.NoError(t, err)

	_, err = cs.AddProductToCart(t.Context(), &services.AddProductToCartParams{
		CartID:    &item.CartID,
		ProductID: uuid.MustParse("01971efc-d170-7664-9ae8-82f386ff59fe"),
		Quantity:  decimal.NewFromInt(3),
	})
	assert.NoError(t, err)

	rs := coinapi.NewMockClient(100_000)
	srv := services.NewOrderService(db, rs)

	email := "user@example.com"
	order, err := srv.CreateOrder(t.Context(), &services.CreateOrderParams{
		CartID: item.CartID,
		Email:  email,
	})
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "user@example.com", order.EmailEncrypted.String())
	assert.Equal(t, "280", order.GrandTotalEur.String())
	assert.Equal(t, "0.0028", order.GrandTotalBtc.String())
	assert.Equal(t, "100000", order.ExchangeRate.String())

	var actualEmail []byte
	err = db.QueryRow(t.Context(), "select email_encrypted from orders where id = $1", order.ID).Scan(&actualEmail)
	assert.NoError(t, err)
	assert.NotEqual(t, string("user@example.com"), actualEmail)

	var itemCount int
	err = db.QueryRow(t.Context(), `select count(id) from order_line_items where order_id = $1`, order.ID).Scan(&itemCount)
	assert.Equal(t, 2, itemCount)

	var exists bool
	err = db.QueryRow(t.Context(), `select exists (select 1 from carts where id = $1)`, item.CartID).Scan(&exists)
	assert.NoError(t, err)
	assert.False(t, exists)
}
