package config

import (
	"github.com/kelseyhightower/envconfig"
)

type ServiceConfig struct {
	ServerHost               string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	ServerPort               int    `envconfig:"SERVER_PORT" default:"8080"`
	ServerGRPCHost           string `envconfig:"SERVER_GRPC_HOST" default:"0.0.0.0"`
	ServerGRPCPort           int    `envconfig:"SERVER_GRPC_PORT" default:"8081"`
	ServiceName              string `envconfig:"SERVICE_NAME" default:"capten-server"`
	Database                 string `envconfig:"DATABASE" default:"astra"`
	AuthEnabled              bool   `envconfig:"AUTH_ENABLED" default:"false"`
	RegisterLaunchAppsConifg bool   `envconfig:"REGISTER_LAUNCH_APPS_CONFIG" default:"true"`
	CaptenOAuthURL           string `envconfig:"CAPTEN_OAUTH_URL" default:"https://alpha.optimizor.app/api/.ory"`
}

func GetServiceConfig() (ServiceConfig, error) {
	cfg := ServiceConfig{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}
