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
)

func main() {
	seed := deriveSeed(config.SECRET_KEY_BASE, []byte("seed"))
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

	client, err := NewClient("username", "password", "127.0.0.1:18443")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	_, _ = client.CreateWallet("watchonly", rpcclient.WithCreateWalletDisablePrivateKeys())

	descInfo, err := client.GetDescriptorInfo(fmt.Sprintf("addr(%s)", addr.EncodeAddress()))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(descInfo.Descriptor)

	resp, err := client.ImportDescriptor(&ImportDescriptorsItem{
		Desc: descInfo.Descriptor,
		Timestamp: ImportDescriptorsTimestamp{
			Timestamp: time.Now().Add(-2 * time.Hour),
		},
		Label: addr.EncodeAddress(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp))

	unspent, err := client.ListUnspent()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", unspent)

	// client.SendCmd(ImportDescriptorsCmd{
	// 	ImportDescriptorsItem{
	// 		Desc:      fmt.Sprintf("addr(%s)", addr.EncodeAddress()),
	// 		Timestamp: time.Now().Add(-2 * time.Hour).Unix(),
	// 	},
	// })
	// blockCount, err := client.GetBlockCount()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Println(blockCount)
}
