package main

import (
	"time"

	"github.com/spf13/viper"
)

const (
	databaseReadURLEnv           = "DATABASE_READ_URL"
	databaseWriteURLEnv          = "DATABASE_WRITE_URL"
	databaseConnRetriesEnv       = "DATABASE_CONNECTION_RETRIES"
	databaseRetryWaitDurationEnv = "DATABASE_RETRY_WAIT_DURATION"
	httpPortEnv                  = "HTTP_PORT"
)

const (
	dbConnRetriesDefault       = 10
	dbRetryWaitDurationDefault = 2 * time.Second
)

type databaseConfig struct {
	ReadURL           string
	WriteURL          string
	ConnRetries       int
	RetryWaitDuration time.Duration
}

type Config struct {
	Database databaseConfig
	HTTPPort int
}

func LoadConfig() *Config {
	viper.SetDefault(databaseConnRetriesEnv, dbConnRetriesDefault)
	viper.SetDefault(databaseRetryWaitDurationEnv, dbRetryWaitDurationDefault)

	viper.AutomaticEnv()

	return &Config{
		Database: databaseConfig{
			ReadURL:           viper.GetString(databaseReadURLEnv),
			WriteURL:          viper.GetString(databaseWriteURLEnv),
			ConnRetries:       viper.GetInt(databaseConnRetriesEnv),
			RetryWaitDuration: viper.GetDuration(databaseRetryWaitDurationEnv),
		},
		HTTPPort: viper.GetInt(httpPortEnv),
	}
}
