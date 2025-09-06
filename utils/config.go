package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServerAddr 		 string 	   `mapstructure:"HTTP_SERVER_ADDR"`
	GRPCServerAddr 		 string 	   `mapstructure:"GRPC_SERVER_ADDR"`
	Dbsource 	   		 string 	   `mapstructure:"DB_SOURCE"`
	MigrationURL 		 string 	   `mapstructure:"MIGRATION_URL"`
	Dbdriver 	   	     string 	   `mapstructure:"DB_DRIVER"`
	PasetoSymmetricKey   string 	   `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddr            string 	   `mapstructure:"REDIS_ADDR"`
	SenderName           string 	   `mapstructure:"SENDER_NAME"`
	SenderEmail          string 	   `mapstructure:"SENDER_EMAIL"`
	SenderPassword       string 	   `mapstructure:"SENDER_PASSWORD"`
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