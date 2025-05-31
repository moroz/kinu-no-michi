package main

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
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
