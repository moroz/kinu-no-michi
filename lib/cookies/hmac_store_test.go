package cookies_test

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/moroz/kinu-no-michi/lib/cookies"
	"github.com/stretchr/testify/assert"
)

var SIGNER, _ = base64.StdEncoding.DecodeString("A5wRvFTyPZupkaKPnU7zISfhYgwpOmQUFhUHAlOThB8=")

func TestHMACStoreEncDec(t *testing.T) {
	store := cookies.HMACStore(sha256.New, SIGNER)
	value := []byte("Hello, world!")

	t.Run("happy path", func(t *testing.T) {
		cookie := store.Encode(value)
		assert.NotEqual(t, "", cookie)

		decoded, err := store.Decode(cookie)
		assert.NoError(t, err)
		assert.Equal(t, decoded, value)
	})

	t.Run("empty value", func(t *testing.T) {
		cookie := store.Encode(nil)
		assert.NotEqual(t, "", cookie)

		decoded, err := store.Decode(cookie)
		assert.NoError(t, err)
		assert.Equal(t, decoded, []byte{})
	})

	t.Run("cookie too short", func(t *testing.T) {
		cookie := base64.RawURLEncoding.EncodeToString([]byte("Hello!"))
		actual, err := store.Decode(cookie)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, cookies.ErrMalformedCookie)
	})

	t.Run("invalid signature", func(t *testing.T) {
		otherStore := cookies.HMACStore(sha256.New, []byte("Hello!"))
		cookie := otherStore.Encode([]byte("I am a fake cookie"))
		assert.NotEqual(t, "", cookie)

		actual, err := store.Decode(cookie)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, cookies.ErrInvalidSignature)
	})
}
