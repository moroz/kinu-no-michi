package encrypt

import (
	"database/sql"
	"database/sql/driver"
	"errors"
)

type EncryptedBytes []byte

var (
	_ driver.Valuer = &EncryptedBytes{}
	_ sql.Scanner   = &EncryptedBytes{}
)

var ErrProviderNotSet = errors.New("encryption provider is not set")

func NewEncryptedString(s string) EncryptedBytes {
	return NewEncryptedBytes([]byte(s))
}

func NewEncryptedBytes(b []byte) EncryptedBytes {
	if len(b) == 0 {
		return nil
	}

	return EncryptedBytes(b)
}

func (e EncryptedBytes) String() string {
	return string(e)
}

func (e EncryptedBytes) Bytes() []byte {
	return e[:]
}

// Value implements driver.Valuer interface
func (e EncryptedBytes) Value() (driver.Value, error) {
	if len(e) == 0 {
		return nil, nil
	}

	if globalProvider == nil {
		return nil, ErrProviderNotSet
	}

	return globalProvider.Encrypt(e)
}

func (e *EncryptedBytes) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("failed to read byte slice")
	}

	if globalProvider == nil {
		return ErrProviderNotSet
	}

	decrypted, err := globalProvider.Decrypt(b)
	if err != nil {
		return err
	}

	*e = decrypted
	return nil
}
