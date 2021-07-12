package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	configVars := map[string]string{
		"APP_PORT":                  "4545",
		"LOG_LEVEL":                 "debug",
		"CHROME_DP_TIMEOUT_SECONDS": "10",
	}

	for k, v := range configVars {
		_ = os.Setenv(k, v)
	}

	Load()
	assert.Equal(t, 4545, Port())
	assert.Equal(t, "debug", LogLevel())
	assert.Equal(t, time.Duration(10)*time.Second, ChromeDPTimeout())

	for k := range configVars {
		_ = os.Unsetenv(k)
	}
}
