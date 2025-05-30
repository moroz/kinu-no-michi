package main

import (
	"fmt"
	"log"
	"time"

	"github.com/btcsuite/btcd/rpcclient"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/moroz/kinu-no-michi/config"
	"golang.org/x/crypto/argon2"
)

type ImportDescriptorsPayload struct {
}

func deriveKey(parent *hdkeychain.ExtendedKey, path ...uint32) (key *hdkeychain.ExtendedKey, err error) {
	key = parent
	for _, i := range path {
		key, err = key.Derive(i)
		if err != nil {
			return
		}
	}
	return
}

func initRpcClient() *rpcclient.Client {
	cfg := &rpcclient.ConnConfig{
		User:         "username",
		Pass:         "password",
		HTTPPostMode: true,
		Host:         "localhost:18443",
		DisableTLS:   true,
	}

	client, err := rpcclient.New(cfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	seed := argon2.IDKey(config.SECRET_KEY_BASE, []byte("seed"), 2, 46*1024, 1, 32)
	master, err := hdkeychain.NewMaster(seed, &chaincfg.RegressionNetParams)
	if err != nil {
		log.Fatal(err)
	}

	// m/84'/0'/0'/0/0
	key, err := deriveKey(master, 84+hdkeychain.HardenedKeyStart, hdkeychain.HardenedKeyStart, hdkeychain.HardenedKeyStart, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	pubKey, err := key.ECPubKey()
	if err != nil {
		log.Fatal(err)
	}

	witnessProg := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.RegressionNetParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(addr.EncodeAddress())

	client := initRpcClient()
	defer client.Shutdown()

	_, _ = client.CreateWallet("watchonly", rpcclient.WithCreateWalletDisablePrivateKeys())

	if err != nil {
		log.Fatal(err)
	}

	client.SendCmd(ImportDescriptorsCmd{
		ImportDescriptorsItem{
			Desc:      fmt.Sprintf("addr(%s)", addr.EncodeAddress()),
			Timestamp: time.Now().Add(-2 * time.Hour).Unix(),
		},
	})
	// blockCount, err := client.GetBlockCount()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Println(blockCount)
}
