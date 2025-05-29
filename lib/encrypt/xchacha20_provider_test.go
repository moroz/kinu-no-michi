package encrypt_test

import (
	"crypto/rand"
	"testing"

	"github.com/moroz/kinu-no-michi/lib/encrypt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20"
)

func TestInitProvider(t *testing.T) {
	t.Run("is valid with valid key length", func(t *testing.T) {
		var key = make([]byte, chacha20.KeySize)
		rand.Read(key)
		actual, err := encrypt.NewXChacha20Provider(key)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("is invalid with too short key", func(t *testing.T) {
		var key = make([]byte, chacha20.KeySize-1)
		rand.Read(key)
		actual, err := encrypt.NewXChacha20Provider(key)
		assert.ErrorIs(t, err, encrypt.ErrKeyLength)
		assert.Nil(t, actual)
	})

	t.Run("is invalid with too long key", func(t *testing.T) {
		var key = make([]byte, chacha20.KeySize+1)
		rand.Read(key)
		actual, err := encrypt.NewXChacha20Provider(key)
		assert.ErrorIs(t, err, encrypt.ErrKeyLength)
		assert.Nil(t, actual)
	})
}

func TestEncryptDecrypt(t *testing.T) {
	var key = make([]byte, chacha20.KeySize)
	rand.Read(key)
	provider, err := encrypt.NewXChacha20Provider(key)
	require.NoError(t, err)
	require.NotNil(t, provider)

	msg := []byte("Hello, deer YouTube fans!")
	ciphertext, err := provider.Encrypt(msg)
	assert.NoError(t, err)
	assert.Len(t, ciphertext, len(msg)+24+16)

	decrypted, err := provider.Decrypt(ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, msg, decrypted)
}
