package config

import (
	"github.com/kelseyhightower/envconfig"
)

type ServiceConfig struct {
	ServerHost         string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	ServerPort         int    `envconfig:"SERVER_PORT" default:"8080"`
	ServerGRPCHost     string `envconfig:"SERVER_GRPC_HOST" default:"0.0.0.0"`
	ServerGRPCPort     int    `envconfig:"SERVER_GRPC_PORT" default:"8081"`
	Database           string `envconfig:"DATABASE" default:"astra"`
	AppStorConfig      string `envconfig:"APP_STORE_CONFIG" default:"./storeconfig"`
	ReadAppStoreConfig bool   `envconfig:"READ_APP_STORE_CONFIG" default:"true"`
  ServiceRegister bool   `envconfig:"SERVICE_REGISTER" default:"false"`
}

func GetServiceConfig() (ServiceConfig, error) {
	cfg := ServiceConfig{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}
