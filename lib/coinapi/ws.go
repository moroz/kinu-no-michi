package coinapi

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/net/websocket"
)

var apiEndpoint = &url.URL{
	Scheme: "wss",
	Host:   "api-realtime.exrates.coinapi.io",
	Path:   "/",
}

var origin = &url.URL{
	Scheme: "http",
	Host:   "localhost",
}

type subscribeMessage struct {
	Type                         string   `json:"type"`
	Heartbeat                    bool     `json:"heartbeat"`
	SubscribeFilterAssetID       []string `json:"subscribe_filter_asset_id"`
	SubscribeUpdateLimitMSExrate int      `json:"subscribe_update_limit_ms_exrate"`
}

type coinapiWSClient struct {
	token     string // API token
	mu        sync.RWMutex
	lastEvent *ExchangeRateEvent
	done      chan struct{}
	conn      *websocket.Conn
}

type ExchangeRateEvent struct {
	Time         time.Time
	Rate         decimal.Decimal
	AssetIDBase  string `json:"asset_id_base"`
	AssetIDQuote string `json:"asset_id_quote"`
}

func NewCoinAPIWSClient(token string) *coinapiWSClient {
	return &coinapiWSClient{
		token: token,
	}
}

func (c *coinapiWSClient) updateLastEvent(e *ExchangeRateEvent) {
	if e == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lastEvent == nil || c.lastEvent.Time.Before(e.Time) {
		c.lastEvent = e
	}
}

func (c *coinapiWSClient) GetLatestRate() *decimal.Decimal {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastEvent == nil {
		return nil
	}

	copied := c.lastEvent.Rate
	return &copied
}

func (c *coinapiWSClient) loop() {
	var buf = make([]byte, 1024)
	for {
		select {
		case <-c.done:
			return

		default:
			n, err := c.conn.Read(buf)
			if err != nil {
				log.Printf("Failed to read message: %s", err)
				return
			}

			var payload ExchangeRateEvent
			err = json.Unmarshal(buf[:n], &payload)
			if err != nil {
				log.Printf("Failed to parse incoming message as JSON: %s", err)
				log.Println(string(buf[:n]))
				return
			}

			log.Printf("Received event: %v", payload)

			c.updateLastEvent(&payload)
		}
	}
}

func (c *coinapiWSClient) Start() error {
	log.Printf("Connecting to %s", apiEndpoint.String())

	config := websocket.Config{
		Location: apiEndpoint,
		Origin:   &url.URL{Scheme: "http", Host: "localhost", Path: "/"},
		Header: http.Header{
			"Authorization": {c.token},
		},
		Version: websocket.ProtocolVersionHybi13,
	}

	conn, err := websocket.DialConfig(&config)
	if err != nil {
		return err
	}

	msg, _ := json.Marshal(subscribeMessage{
		Type:                         "hello",
		SubscribeFilterAssetID:       []string{"BTC/EUR"},
		SubscribeUpdateLimitMSExrate: 5000,
	})

	if _, err := conn.Write(msg); err != nil {
		return err
	}

	c.conn = conn
	c.done = make(chan struct{})

	go c.loop()
	return nil
}

func (c *coinapiWSClient) Stop() {
	c.done <- struct{}{}
}
