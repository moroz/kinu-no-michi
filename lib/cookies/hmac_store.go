package cookies

import (
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"hash"
)

type hmacStore struct {
	hash func() hash.Hash
	key  []byte
}

var ErrMalformedCookie = errors.New("the provided cookie is incorrectly encoded or too short")
var ErrInvalidSignature = errors.New("the provided signature does not match the expected value")

func HMACStore(hash func() hash.Hash, key []byte) SessionStore {
	return &hmacStore{
		hash: hash,
		key:  key,
	}
}

func (s *hmacStore) digest(v []byte) []byte {
	mac := hmac.New(s.hash, s.key)
	mac.Write(v)
	return mac.Sum(nil)
}

func (s *hmacStore) Encode(v []byte) string {
	sum := append(v, s.digest(v)...)
	return base64.RawURLEncoding.EncodeToString(sum)
}

func constTimeEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var acc byte

	for i := range a {
		acc |= a[i] ^ b[i]
	}

	return acc == 0
}

func (s *hmacStore) Decode(cookie string) ([]byte, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(cookie)
	if err != nil {
		return nil, err
	}

	size := s.hash().Size()
	if len(bytes) < size {
		return nil, ErrMalformedCookie
	}

	msg, sum := bytes[:len(bytes)-size], bytes[len(bytes)-size:]
	expected := s.digest(msg)

	if constTimeEq(sum, expected) {
		return msg, nil
	}

	return nil, ErrInvalidSignature
}
