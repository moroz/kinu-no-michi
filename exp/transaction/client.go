package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
)

type Client struct {
	*rpcclient.Client
	user string
	pass string
	host string
}

type RpcCmd struct {
	Version string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type ImportDescriptorsTimestamp struct {
	Now       bool
	Timestamp time.Time
}

type ImportDescriptorsItem struct {
	Desc      string                     `json:"desc"`
	Timestamp ImportDescriptorsTimestamp `json:"timestamp"`
	Label     string                     `json:"label,omitempty"`
	Active    bool                       `json:"active,omitempty"`
	Range     []int                      `json:"range,omitempty"`
}

func (i *ImportDescriptorsTimestamp) MarshalJSON() ([]byte, error) {
	if i.Now {
		return json.Marshal("now")
	}
	return json.Marshal(i.Timestamp.Unix())
}

func NewClient(username, password, host string) (*Client, error) {
	rpc, err := rpcclient.New(&rpcclient.ConnConfig{
		User:         username,
		Pass:         password,
		Host:         host,
		HTTPPostMode: true,
		DisableTLS:   true,
	}, nil)

	if err != nil {
		return nil, err
	}

	return &Client{rpc, username, password, host}, nil
}

func (c *Client) SendRawCmd(method string, params []any) ([]byte, error) {
	var id = make([]byte, 4)
	rand.Read(id)

	payload, err := json.Marshal(RpcCmd{
		Version: "1.0",
		Id:      hex.EncodeToString(id),
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://"+c.host, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.user, c.pass)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) ImportDescriptor(params *ImportDescriptorsItem) ([]byte, error) {
	return c.ImportDescriptors([]*ImportDescriptorsItem{params})
}

func (c *Client) ImportDescriptors(items []*ImportDescriptorsItem) ([]byte, error) {
	return c.SendRawCmd("importdescriptors", []any{items})
}
