package config

import (
	"github.com/kelseyhightower/envconfig"
)

type SericeConfig struct {
	Host        string `envconfig:"HOST" default:"0.0.0.0"`
	Port        int    `envconfig:"PORT" default:"9091"`
	Mode        string `envconfig:"MODE" default:"production"`
	AuthEnabled bool   `envconfig:"AUTH_ENABLED" default:"false"`
	CreatePR    bool   `envconfig:"CREATE_PR" default:"false"`
}

func GetServiceConfig() (*SericeConfig, error) {
	cfg := &SericeConfig{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
