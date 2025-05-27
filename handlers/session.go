package handlers

import (
	"bytes"
	"encoding/gob"

	"github.com/google/uuid"
)

func init() {
	gob.Register(appSession{})
}

type appSession struct {
	CartID uuid.UUID
}

func decodeSession(binary []byte) (*appSession, error) {
	var result appSession
	err := gob.NewDecoder(bytes.NewBuffer(binary)).Decode(&result)
	return &result, err
}

func encodeSession(s *appSession) []byte {
	if s == nil {
		return nil
	}

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(*s)
	return buf.Bytes()
}
