package gokinesis

import (
	"github.com/spf13/viper"
)

type Config struct {
	StreamName      string `mapstructure:"STREAM_NAME"`
	AWSRegion       string `mapstructure:"AWS_REGION"`
	Endpoint        string `mapstructure:"ENDPOINT"`
	AccessKeyID     string `mapstructure:"ACCESS_KEY_ID"`
	SecretAccessKey string `mapstructure:"SECRET_ACCESS_KEY"`
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
