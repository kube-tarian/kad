package iamclient

import (
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"github.com/pkg/errors"
)

type Config struct {
	IAMURL          string `envconfig:"IAM_URL" required:"true"`
	ServiceRegister bool   `envconfig:"SERVICE_REGISTER" default:"false"`
	ServiceName     string `envconfig:"SERVICE_NAME" default:"capten-server"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func RegisterService(log logging.Logger) error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	if !cfg.ServiceRegister {
		log.Infof("service registration disabled")
		return nil
	}

	oryclient, err := oryclient.NewOryClient(log)
	if err != nil {
		return errors.WithMessage(err, "OryClient initialization failed")
	}

	IC, err := NewClient(log, oryclient, cfg)
	if err != nil {
		return errors.WithMessage(err, "Error occured while created IAM client")
	}

	err = IC.RegisterWithIam()
	if err != nil {
		return errors.WithMessage(err, "Registering capten server as oauth client failed")
	}

	err = IC.RegisterRolesActions()
	if err != nil {
		return errors.WithMessage(err, "Registering Roles and Actions in IAM failed")
	}
	return nil
}
