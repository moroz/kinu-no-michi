package cookies

type SessionStore interface {
	Encode(v []byte) string
	Decode(cookie string) ([]byte, error)
}
