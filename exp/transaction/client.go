package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/shopspring/decimal"
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

func (c *Client) SendRawCmd(method string, params ...any) ([]byte, error) {
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
	return c.SendRawCmd("importdescriptors", items)
}

type GetRawTransactionWithBlockResult struct {
	Result btcjson.TxRawResult
	Err    error
}

func (c *Client) GetRawTransactionWithBlock(txId, blockHash string) (*btcjson.TxRawResult, error) {
	bytes, err := c.SendRawCmd("getrawtransaction", txId, true, blockHash)
	if err != nil {
		return nil, err
	}

	var result GetRawTransactionWithBlockResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return &result.Result, nil
}

func (c *Client) ImportSegWitAddress(address string) ([]byte, error) {
	descInfo, err := c.GetDescriptorInfo(fmt.Sprintf("addr(%s)", address))
	if err != nil {
		return nil, err
	}

	return c.ImportDescriptor(&ImportDescriptorsItem{
		Desc: descInfo.Descriptor,
		Timestamp: ImportDescriptorsTimestamp{
			Timestamp: time.Now().Add(-2 * time.Hour),
		},
		Label: address,
	})
}

func (c *Client) ValidateTransaction(tx *wire.MsgTx) ([]byte, error) {
	var buf bytes.Buffer
	tx.Serialize(&buf)
	hex := hex.EncodeToString(buf.Bytes())
	return c.SendRawCmd("testmempoolaccept", []string{hex})
}

func (c *Client) SendRawTransaction(tx *wire.MsgTx, maxFeeRate int64) ([]byte, error) {
	var buf bytes.Buffer
	tx.Serialize(&buf)
	hex := hex.EncodeToString(buf.Bytes())
	rate := decimal.NewFromInt(maxFeeRate).Div(decimal.NewFromInt(1e5)).String()
	return c.SendRawCmd("sendrawtransaction", hex, rate)
}
