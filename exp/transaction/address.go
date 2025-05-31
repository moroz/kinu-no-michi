package main

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"golang.org/x/crypto/argon2"
)

func deriveSeed(base, salt []byte) []byte {
	return argon2.IDKey(base, salt, 2, 46*1024, 1, 32)
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

func deriveAddress(parent *hdkeychain.ExtendedKey, netParams *chaincfg.Params, path ...uint32) (string, error) {
	key, err := deriveKey(parent, path...)
	if err != nil {
		return "", err
	}

	pubKey, err := key.ECPubKey()
	if err != nil {
		return "", err
	}

	witnessProg := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, netParams)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}
