package encrypt

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

var (
	ErrKeyLength    = errors.New("invalid key length")
	ErrMalformedMsg = errors.New("message is too short")
)

type XChacha20Provider struct {
	aead cipher.AEAD
}

func NewXChacha20Provider(key []byte) (EncryptionProvider, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("%w (want %d, got %d)", ErrKeyLength, chacha20poly1305.KeySize, len(key))
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	return &XChacha20Provider{aead}, nil
}

func (p *XChacha20Provider) Encrypt(value []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSizeX, chacha20poly1305.NonceSizeX+len(value)+chacha20poly1305.Overhead)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	result := p.aead.Seal(nonce, nonce, value, nil)
	return result, nil
}

func (p *XChacha20Provider) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < chacha20poly1305.NonceSizeX+chacha20poly1305.Overhead {
		return nil, ErrMalformedMsg
	}

	nonce, msg := ciphertext[:chacha20poly1305.NonceSizeX], ciphertext[chacha20poly1305.NonceSizeX:]
	return p.aead.Open(nil, nonce, msg, nil)
}
