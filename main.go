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
	restClient, err := coinapi.NewCoinAPIRESTClient(COINAPI_API_KEY, 60000)
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.Router(restClient)

	log.Printf("Listening on %s", LISTEN_ON)

	log.Fatal(http.ListenAndServe(LISTEN_ON, r))
}
