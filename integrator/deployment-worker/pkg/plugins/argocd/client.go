package argocd

import (
	"github.com/kelseyhightower/envconfig"
)

type ArgoCDCLient struct {
	conf *Configuration
}

func NewClient() (*ArgoCDCLient, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}

	return &ArgoCDCLient{
		conf: cfg,
	}, nil
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
