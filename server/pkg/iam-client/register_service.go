package iamclient

import (
	cm "github.com/intelops/go-common/iam"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"github.com/pkg/errors"
)

type Config struct {
	IAMURL                       string `envconfig:"IAM_URL" required:"true"`
	ServiceRegister              bool   `envconfig:"SERVICE_REGISTER" default:"true"`
	ServiceName                  string `envconfig:"SERVICE_NAME" default:"capten-server"`
	ServiceRolesConfigFilePath   string `envconfig:"SERVICE_ROLES_CONFIG_FILE_PATH" default:"/data/service-config/roles.yaml"`
	CerbosResourcePolicyFilePath string `envconfig:"SERVICE_ROLES_CONFIG_FILE_PATH" default:"/data/cerbos-policy-config/resource-policy.yaml"`
	PolicyRegister               bool   `envconfig:"CERBOS_POLICY_REGISTER" default:"true"`
}

func NewConfig() (Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
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

	iamClient, err := NewClient(log, oryclient, cfg)
	if err != nil {
		return errors.WithMessage(err, "Error occured while created IAM client")
	}

	err = iamClient.RegisterRolesActions()
	if err != nil {
		return errors.WithMessage(err, "Registering Roles and Actions in IAM failed")
	}
	log.Infof("service registration successful")
	return nil
}

func RegisterCerbosPolicy(log logging.Logger) error {
	cfg, err := NewConfig()
	if err != nil {
		return err
	}

	if !cfg.PolicyRegister {
		log.Infof("Cerbos Policy registration disabled")
		return nil
	}

	oryclient, err := oryclient.NewOryClient(log)
	if err != nil {
		return errors.WithMessage(err, "OryClient initialization failed")
	}

	iamClient, err := NewClient(log, oryclient, cfg)
	if err != nil {
		return errors.WithMessage(err, "Error occured while created IAM client")
	}

	err = iamClient.RegisterCerbosPolicy()
	if err != nil {
		return errors.WithMessage(err, "Registering Resource policy in cerbos through IAM failed")
	}
	log.Infof("cerbos policy registration successful")
	return nil
}

func Interceptor(log logging.Logger) (*cm.ClientsAndConfigs, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}
	oryclient, err := oryclient.NewOryClient(log)
	if err != nil {
		log.Errorf("OryClient initialization failed %v", err)
		return nil, errors.WithMessage(err, "OryClient initialization failed")
	}

	iamClient, err := NewClient(log, oryclient, cfg)
	if err != nil {
		log.Errorf("Error occured while created IAM client %v", err)
		return nil, errors.WithMessage(err, "Error occured while created IAM client")
	}
	interceptor, err := iamClient.Interceptor()
	if err != nil {
		log.Errorf("Error occured while creating Unary Interceptor %v", err)
		return nil, errors.WithMessage(err, "Error occured while creating Unary Interceptor")
	}
	return interceptor, nil
}
