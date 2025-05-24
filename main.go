package main

import (
	"context"
	"crypto/sha256"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/handlers"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
	"github.com/moroz/kinu-no-michi/lib/cookies"
)

const LISTEN_ON = ":3000"

func main() {
	cookieStore := cookies.HMACStore(sha256.New, config.COOKIE_HMAC_KEY)

	db, err := pgxpool.New(context.Background(), config.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	restClient, err := coinapi.NewCoinAPIRESTClient(config.COINAPI_API_KEY, 60000)
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.Router(db, restClient, cookieStore)

	log.Printf("Listening on %s", LISTEN_ON)

	log.Fatal(http.ListenAndServe(LISTEN_ON, r))
}
