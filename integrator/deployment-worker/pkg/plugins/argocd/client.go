package argocd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
	"github.com/kube-tarian/kad/integrator/pkg/logging"
)

type ArgoCDCLient struct {
	conf       *Configuration
	httpClient *http.Client
	logger     logging.Logger
}

func NewClient(logger logging.Logger) (*ArgoCDCLient, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if cfg.IsSSLEnabled {
		// TODO: Configure SSL certificates
		logger.Errorf("SSL not yet supported, continuing with insecure verify true")
	}
	return &ArgoCDCLient{
		conf:       cfg,
		httpClient: &http.Client{Transport: tr},
		logger:     logger,
	}, nil
}

func (a *ArgoCDCLient) Exec(payload model.RequestPayload) (json.RawMessage, error) {
	var err error

	switch payload.Action {
	case "install":
		err = a.Create(payload)
	case "delete":
		err = a.Delete(payload)
	default:
		err = fmt.Errorf("unsupported action for argocd plugin: %v", payload.Action)
	}
	if err != nil {
		a.logger.Errorf("argocd %v of application failed, %v", payload.Action, err)
		return nil, err
	}

	return nil, nil
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
