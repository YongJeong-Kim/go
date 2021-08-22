package config

import (
	"github.com/spf13/viper"
	"time"
)

type App struct {
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RefreshSecret        string        `mapstructure:"REFRESH_SECRET"`
	TokenSymmetricKey    string        `mapstructur:"TOKEN_SYMMETRIC_KEY"`
}

func LoadConfig(path string) (config App, err error) {
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
