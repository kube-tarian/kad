package captensdk

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
)

type CaptenAgentConfiguration struct {
	AgentAddress string `envconfig:"AGENT_ADDRESS" default:"localhost"`
	AgentPort    int    `envconfig:"AGENT_PORT" default:"9091"`
}

type CaptenClient interface {
	NewApplicationClient() (*ApplicationClient, error)
	NewClimonClient(opts *TransportSSLOptions) (*ClimonClient, error)
	NewClusterClient() (*ClusterClient, error)
	NewRepoClient() (*RepoClient, error)
	NewProjectClient() (*ProjectClient, error)
}

type Client struct {
	log  logging.Logger
	conf *CaptenAgentConfiguration
}

func NewClient(log logging.Logger) (*Client, error) {
	cfg := &CaptenAgentConfiguration{}
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Errorf("Capten agent configuration not provided: %v", err)
	}

	return NewClientWithConfiguratin(cfg, log)
}

func NewClientWithConfiguratin(cfg *CaptenAgentConfiguration, log logging.Logger) (*Client, error) {
	return &Client{
		log:  log,
		conf: cfg,
	}, nil
}
