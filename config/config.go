package config

import (
	"encoding/base64"
	"log"
	"os"

	"golang.org/x/crypto/argon2"
)

func MustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("FATAL: Environment variable %s is not set!", key)
	}
	return val
}

func MustGetenvBase64(key string) []byte {
	val := MustGetenv(key)
	decoded, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		log.Fatalf("FATAL: Could not decode environment variable %s from Base64: %s", key, err)
	}
	return decoded
}

var COINAPI_API_KEY = MustGetenv("COINAPI_API_KEY")
var DATABASE_URL = MustGetenv("DATABASE_URL")
var SECRET_KEY_BASE = MustGetenvBase64("SECRET_KEY_BASE")

var COOKIE_HMAC_KEY = argon2.IDKey(SECRET_KEY_BASE, []byte("cookie hmac"), 1, 64*1024, 1, 32)
