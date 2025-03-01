package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddr string `mapstructure:"SERVER_ADDR"`
	Dbsource string `mapstructure:"DB_SOURCE"`
	Dbdriver string `mapstructure:"DB_DRIVER"`
	PasetoSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}