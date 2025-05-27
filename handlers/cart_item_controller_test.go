package handlers_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/handlers"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/lib/cookies"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddProductToCart(t *testing.T) {
	db, err := pgxpool.New(context.Background(), config.MustGetenv("TEST_DATABASE_URL"))
	require.NoError(t, err)

	_, err = db.Exec(context.Background(), "truncate carts cascade;")
	require.NoError(t, err)

	rs := coinapi.NewMockClient(95000)

	cookieKey, err := base64.StdEncoding.DecodeString("A5wRvFTyPZupkaKPnU7zISfhYgwpOmQUFhUHAlOThB8=")
	require.NoError(t, err)
	cs := cookies.HMACStore(sha256.New, cookieKey)

	router := handlers.Router(db, rs, cs)

	count := func(table string) int {
		var res int
		query := fmt.Sprintf("select count(*) from %s", table)
		if err := db.QueryRow(context.Background(), query).Scan(&res); err != nil {
			t.Fatal(err)
		}
		return res
	}

	const CONTENT_TYPE = "application/x-www-form-urlencoded"
	const ENDPOINT = "/cart_items"
	const PRODUCT_ID = "019709a2-5c37-73e2-a05b-9ee9f8a470b5"

	cartsBefore := count("carts")
	itemsBefore := count("cart_items")

	payload := url.Values{
		"product_id": {PRODUCT_ID},
		"quantity":   {"3"},
	}.Encode()

	req, err := http.NewRequest("POST", ENDPOINT, bytes.NewBufferString(payload))
	require.NoError(t, err)

	req.Header.Add("Content-Type", CONTENT_TYPE)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.GreaterOrEqual(t, rr.Code, 200)
	assert.Less(t, rr.Code, 400)

	assert.Equal(t, cartsBefore+1, count("carts"))
	assert.Equal(t, itemsBefore+1, count("cart_items"))

	actual := rr.Result().Cookies()
	assert.Len(t, actual, 1)

	cookie := actual[0]
	assert.Equal(t, config.SESSION_COOKIE_NAME, cookie.Name)

	payload = url.Values{
		"product_id": {PRODUCT_ID},
		"quantity":   {"2"},
	}.Encode()
	req, err = http.NewRequest("POST", ENDPOINT, bytes.NewBufferString(payload))
	req.Header.Add("Content-Type", CONTENT_TYPE)
	req.AddCookie(&http.Cookie{
		Name:     config.SESSION_COOKIE_NAME,
		Value:    cookie.Value,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, cartsBefore+1, count("carts"))
	assert.Equal(t, itemsBefore+1, count("cart_items"))

	var cartID uuid.UUID
	err = db.QueryRow(context.Background(), "select id from carts order by id desc limit 1").Scan(&cartID)
	require.NoError(t, err)

	cart, err := queries.New(db).GetCartByID(context.Background(), cartID)
	assert.Equal(t, decimal.NewFromInt(250), cart.GrandTotal)
	assert.Equal(t, int64(1), cart.ItemCount)
}
