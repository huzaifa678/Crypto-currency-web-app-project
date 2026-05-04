//nolint:revive
package utils

import (
	"os"
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
	GoogleClientID       string        `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret   string        `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL    string        `mapstructure:"REDIRECT_URL"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	Origin               string        `mapstructure:"ORIGIN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetConfigName("app")
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if os.Getenv("APP_ENV") == "test" {
		viper.SetConfigName("app.test")
		_ = viper.MergeInConfig() 
	}


	err = viper.Unmarshal(&config)

	return
}