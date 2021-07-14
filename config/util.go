package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func mustGetInt(key string) int {
	v, err := strconv.Atoi(mustGetString(key))
	mustParseKey(err, key)
	return v
}

func mustGetString(key string) string {
	mustHaveKey(key)
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}

func mustHaveKey(key string) {
	if !viper.IsSet(key) && os.Getenv(key) == "" {
		log.Fatalf("%s key is not set", key)
	}
}

func mustParseKey(err error, key string) {
	if err != nil {
		log.Fatalf("Could not parse key: %s, Error: %s", key, err)
	}
}
