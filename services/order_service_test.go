package services_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/lib/encrypt"
	"github.com/moroz/kinu-no-michi/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	db, err := pgxpool.New(context.Background(), config.MustGetenv("TEST_DATABASE_URL"))
	require.NoError(t, err)

	_, err = db.Exec(context.Background(), "truncate orders cascade;")
	require.NoError(t, err)

	var key, _ = base64.StdEncoding.DecodeString("mAcJgVR/M5/OejAQNtbTeIeb+O7AoLuZ2purSckAKuM=")
	provider, err := encrypt.NewXChacha20Provider(key)
	require.NoError(t, err)
	encrypt.SetProvider(provider)

	srv := services.NewOrderService(db)

	order, err := srv.CreateOrder(context.Background(), "user@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "user@example.com", order.EmailEncrypted.String())

	var actualEmail []byte
	err = db.QueryRow(context.Background(), "select email_encrypted from orders where id = $1", order.ID).Scan(&actualEmail)
	assert.NoError(t, err)
	assert.NotEqual(t, string("user@example.com"), actualEmail)
}
