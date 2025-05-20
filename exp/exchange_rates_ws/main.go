package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/moroz/kinu-no-michi/config"
	"golang.org/x/net/websocket"
)

var COINAPI_API_KEY = config.MustGetenv("COINAPI_API_KEY")

const API_ENDPOINT = "wss://api-realtime.exrates.coinapi.io"

type subscribeMessage struct {
	Type                         string   `json:"type"`
	Heartbeat                    bool     `json:"heartbeat"`
	SubscribeFilterAssetID       []string `json:"subscribe_filter_asset_id"`
	SubscribeUpdateLimitMSExrate int      `json:"subscribe_update_limit_ms_exrate"`
}

func main() {
	log.Printf("Connecting to %s", API_ENDPOINT)
	endpoint, _ := url.Parse(API_ENDPOINT)

	config := websocket.Config{
		Location: endpoint,
		Origin:   &url.URL{Scheme: "http", Host: "localhost", Path: "/"},
		Header: http.Header{
			"Authorization": {COINAPI_API_KEY},
		},
		Version: websocket.ProtocolVersionHybi13,
	}

	conn, err := websocket.DialConfig(&config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	msg, _ := json.Marshal(subscribeMessage{
		Type:                         "hello",
		SubscribeFilterAssetID:       []string{"BTC/EUR"},
		SubscribeUpdateLimitMSExrate: 5000,
	})

	if _, err := conn.Write(msg); err != nil {
		log.Fatal(err)
	}

	var buf = make([]byte, 1024)
	for {
		if _, err := conn.Read(buf); err != nil {
			log.Printf("Failed to read message from WS: %s", err)
			break
		}
		fmt.Println(string(buf))
	}

	log.Println("Closing connection...")
}
