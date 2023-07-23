package config

import (
	"github.com/kelseyhightower/envconfig"
)

type ServiceConfig struct {
	ServerHost     string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	ServerPort     int    `envconfig:"SERVER_PORT" default:"8080"`
	ServerHTTPHost string `envconfig:"SERVER_HTTP_HOST" default:"0.0.0.0"`
	ServerHTTPPort int    `envconfig:"SERVER_HTTP_PORT" default:"8081"`
	Database       string `envconfig:"DATABASE" default:"astra"`
}

func GetServiceConfig() (ServiceConfig, error) {
	cfg := ServiceConfig{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}
