package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestConfig() {
	configVars := map[string]string{
		"APP_PORT":                  "4545",
		"LOG_LEVEL":                 "debug",
		"CHROME_DP_TIMEOUT_SECONDS": "10",
	}

	for k, v := range configVars {
		_ = os.Setenv(k, v)
	}

	Load()
	s.Equal(4545, Port())
	s.Equal("debug", LogLevel())
	s.Equal(time.Duration(10)*time.Second, ChromeDPTimeout())

	for k := range configVars {
		_ = os.Unsetenv(k)
	}
}
