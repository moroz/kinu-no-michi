package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/moroz/kinu-no-michi/config"
)

func main() {
	netParams := &chaincfg.RegressionNetParams

	seed := deriveSeed(config.SECRET_KEY_BASE, []byte("seed"))
	master, err := hdkeychain.NewMaster(seed, netParams)
	if err != nil {
		log.Fatal(err)
	}

	mainAddress, err := deriveAddress(master, netParams, 84+hdkeychain.HardenedKeyStart, hdkeychain.HardenedKeyStart, hdkeychain.HardenedKeyStart, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	otherAddress, err := deriveAddress(master, netParams, 84+hdkeychain.HardenedKeyStart, hdkeychain.HardenedKeyStart, hdkeychain.HardenedKeyStart, 0, 1)
	if err != nil {
		log.Fatal(err)
	}

	client, err := NewClient("username", "password", "127.0.0.1:18443")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	_, _ = client.CreateWallet("watchonly", rpcclient.WithCreateWalletDisablePrivateKeys())

	for _, addr := range []string{mainAddress, otherAddress} {
		_, err := client.ImportSegWitAddress(addr)
		if err != nil {
			log.Fatal(err)
		}
	}

	unspent, err := client.ListUnspent()
	if err != nil {
		log.Fatal(err)
	}

	if len(unspent) == 0 {
		log.Fatal("No money to spend")
	}

	utxo := unspent[0]
	inHash, err := chainhash.NewHashFromStr(utxo.TxID)
	if err != nil {
		log.Fatal(err)
	}

	msgTx := wire.NewMsgTx(wire.TxVersion)
	msgTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(inHash, utxo.Vout), nil, nil))

	txAmount := int64(100_000)
	totalInputAmount := int64(utxo.Amount * 1e8)
	fee := int64(500)
	change := totalInputAmount - txAmount - fee

	recipient, err := btcutil.DecodeAddress(otherAddress, netParams)
	if err != nil {
		log.Fatal(err)
	}
	pkScript, _ := txscript.PayToAddrScript(recipient)
	msgTx.AddTxOut(wire.NewTxOut(txAmount, pkScript))

	recipient, err = btcutil.DecodeAddress(mainAddress, netParams)
	if err != nil {
		log.Fatal(err)
	}
	pkScript, _ = txscript.PayToAddrScript(recipient)
	msgTx.AddTxOut(wire.NewTxOut(change, pkScript))

	mainKey, _ := deriveKey(master, 84+hdkeychain.HardenedKeyStart, hdkeychain.HardenedKeyStart, hdkeychain.HardenedKeyStart, 0, 0)
	mainPriv, _ := mainKey.ECPrivKey()

	prevPkScript, _ := hex.DecodeString(utxo.ScriptPubKey)
	inputFetcher := txscript.NewCannedPrevOutputFetcher(prevPkScript, totalInputAmount)
	witness, err := txscript.WitnessSignature(msgTx, txscript.NewTxSigHashes(msgTx, inputFetcher), 0, totalInputAmount, prevPkScript, txscript.SigHashAll, mainPriv, true)
	if err != nil {
		log.Fatal(err)
	}
	msgTx.TxIn[0].Witness = witness

	var buf bytes.Buffer
	msgTx.Serialize(&buf)
	fmt.Printf("%X\n", buf.Bytes())
}
