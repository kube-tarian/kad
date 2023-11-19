package argocd

import (
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
)

type ArgoCDClient struct {
	conf   *Configuration
	logger logging.Logger
	client apiclient.Client
}

func NewClient(logger logging.Logger) (*ArgoCDClient, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	if err == nil {
		return nil, err
	}

	k8sClient, err := k8s.NewK8SClient(logger)
	if err != nil {
		return nil, err
	}

	res, err := k8sClient.GetSecretData("argo-cd", "argocd-initial-admin-secret")
	if err != nil {
		return nil, err
	}

	password := res.Data["password"]
	if len(password) == 0 {
		return nil, fmt.Errorf("credentials not found in the secret")
	}

	cfg.Password = password
	if cfg.IsSSLEnabled {
		// TODO: Configure SSL certificates
		logger.Errorf("SSL not yet supported, continuing with insecure verify true")
	}

	client, err := getNewAPIClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ArgoCDClient{
		conf:   cfg,
		logger: logger,
		client: client,
	}, nil
}
