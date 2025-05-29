package encrypt_test

import (
	"crypto/rand"
	"testing"

	"github.com/moroz/kinu-no-michi/lib/encrypt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20"
)

func TestSerializeDeserialize(t *testing.T) {
	plaintext := "Hello world!"
	fromString := encrypt.NewEncryptedString(plaintext)
	fromBytes := encrypt.NewEncryptedBytes([]byte(plaintext))
	assert.Equal(t, fromBytes, fromString)

	ciphertext, err := fromString.Value()
	assert.ErrorIs(t, err, encrypt.ErrProviderNotSet)
	assert.NotEqual(t, []byte(plaintext), ciphertext)

	var key = make([]byte, chacha20.KeySize)
	_, err = rand.Read(key)
	require.NoError(t, err)
	provider, err := encrypt.NewXChacha20Provider(key)
	require.NoError(t, err)
	encrypt.SetProvider(provider)

	ciphertext, err = fromString.Value()
	assert.NoError(t, err)
	assert.NotEqual(t, []byte(plaintext), ciphertext)

	dst := encrypt.EncryptedBytes{}
	err = dst.Scan(ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, []byte(plaintext), []byte(dst))
}
