package config

import (
	"github.com/kelseyhightower/envconfig"
)

type SericeConfig struct {
	Host string `envconfig:"HOST" default:"0.0.0.0"`
	Port int    `envconfig:"PORT" default:"9091"`
}

func GetServiceConfig() (*SericeConfig, error) {
	cfg := &SericeConfig{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
