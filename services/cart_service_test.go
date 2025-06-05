package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/services"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddProductToCart(t *testing.T) {
	db, err := pgxpool.New(context.Background(), config.MustGetenv("TEST_DATABASE_URL"))
	require.NoError(t, err)

	_, err = db.Exec(context.Background(), "truncate carts cascade;")
	require.NoError(t, err)

	srv := services.NewCartService(db)

	productID, err := uuid.Parse("019709a2-5c37-73e2-a05b-9ee9f8a470b5")
	assert.NoError(t, err)

	countCarts := func() int {
		var count int
		err := db.QueryRow(context.Background(), "select count(*) from carts;").Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		return count
	}

	countCartItems := func() int {
		var count int
		err := db.QueryRow(context.Background(), "select count(*) from cart_items;").Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		return count
	}

	t.Run("creates new cart with nil CartID", func(t *testing.T) {
		cartsBefore := countCarts()
		itemsBefore := countCartItems()

		_, err := srv.AddProductToCart(context.Background(), &services.AddProductToCartParams{
			CartID:    uuid.Nil,
			ProductID: productID,
			Quantity:  decimal.NewFromInt(1),
		})

		assert.NoError(t, err)

		assert.Equal(t, cartsBefore+1, countCarts())
		assert.Equal(t, itemsBefore+1, countCartItems())
	})

	t.Run("reuses existing cart if CartID is present", func(t *testing.T) {
		cartID, err := uuid.NewV7()
		assert.NoError(t, err)
		_, err = db.Exec(context.Background(), "insert into carts (id) values ($1)", cartID)
		assert.NoError(t, err)

		cartsBefore := countCarts()
		itemsBefore := countCartItems()

		item, err := srv.AddProductToCart(context.Background(), &services.AddProductToCartParams{
			CartID:    cartID,
			ProductID: productID,
			Quantity:  decimal.NewFromInt(1),
		})
		assert.NoError(t, err)
		assert.Equal(t, cartID, item.CartID)

		assert.Equal(t, cartsBefore, countCarts())
		assert.Equal(t, itemsBefore+1, countCartItems())
	})

	t.Run("reuses cart items with the same product_id", func(t *testing.T) {
		cartsBefore := countCarts()
		itemsBefore := countCartItems()

		item, err := srv.AddProductToCart(context.Background(), &services.AddProductToCartParams{
			CartID:    uuid.Nil,
			ProductID: productID,
			Quantity:  decimal.NewFromInt(2),
		})
		assert.NoError(t, err)

		assert.Equal(t, cartsBefore+1, countCarts())
		assert.Equal(t, itemsBefore+1, countCartItems())

		item, err = srv.AddProductToCart(context.Background(), &services.AddProductToCartParams{
			CartID:    item.CartID,
			ProductID: productID,
			Quantity:  decimal.NewFromInt(3),
		})
		assert.NoError(t, err)

		assert.Equal(t, decimal.NewFromInt(5), item.Quantity)
	})
}
