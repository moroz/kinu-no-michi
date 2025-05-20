package config

import (
	"log"
	"os"
)

func MustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("FATAL: Environment variable %s is not set!", key)
	}
	return val
}
