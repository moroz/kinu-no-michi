package main

import (
	"log"
	"net/http"

	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/handlers"
	"github.com/moroz/kinu-no-michi/lib/coinapi"
)

const LISTEN_ON = ":3000"

var COINAPI_API_KEY = config.MustGetenv("COINAPI_API_KEY")

func main() {
	wsClient := coinapi.NewCoinAPIWSClient(COINAPI_API_KEY)
	err := wsClient.Start()
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.Router(wsClient)

	log.Printf("Listening on %s", LISTEN_ON)

	log.Fatal(http.ListenAndServe(LISTEN_ON, r))
}
