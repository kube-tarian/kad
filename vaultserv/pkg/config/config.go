package config

import "github.com/kelseyhightower/envconfig"

type Configuration struct {
	Host string `envconfig:"HOST" default:"0.0.0.0"`
	Port int    `envconfig:"PORT" default:"9098"`
}

func FetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
