package encrypt

type EncryptionProvider interface {
	Encrypt(value []byte) ([]byte, error)
	Decrypt(value []byte) ([]byte, error)
}

var globalProvider EncryptionProvider

func SetProvider(provider EncryptionProvider) {
	globalProvider = provider
}
