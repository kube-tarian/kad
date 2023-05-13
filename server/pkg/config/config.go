package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Configuration struct {
	Host string `envconfig:"HOST" default:"0.0.0.0"`
	Port int    `envconfig:"PORT" default:"9092"`
}

func FetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	viper.AddConfigPath()
	err := envconfig.Process("", cfg)
	return cfg, err
}
