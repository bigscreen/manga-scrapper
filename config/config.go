package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	port            int
	logLevel        string
	chromeDPTimeout int
}

var appConfig *Config

func Load() {
	viper.SetDefault("APP_PORT", "4545")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.AutomaticEnv()

	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")

	_ = viper.ReadInConfig()

	appConfig = &Config{
		port:            getIntOrPanic("APP_PORT"),
		logLevel:        fatalGetString("LOG_LEVEL"),
		chromeDPTimeout: getIntOrPanic("CHROME_DP_TIMEOUT_SECONDS"),
	}
}

func Port() int {
	return appConfig.port
}

func LogLevel() string {
	return appConfig.logLevel
}

func ChromeDPTimeout() time.Duration {
	return time.Duration(appConfig.chromeDPTimeout) * time.Second
}
