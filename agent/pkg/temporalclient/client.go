package temporalclient

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/agent/pkg/logging"

	"go.temporal.io/sdk/client"
)

type Configuration struct {
	TemporalServiceAddress string `envconfig:"TEMPORAL_SERVICE_URL" default:"localhost:7233"`
}

type Client struct {
	conf           *Configuration
	temporalClient client.Client
	log            logging.Logger
}

func NewClient(log logging.Logger) (*Client, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}

	clnt := &Client{
		conf: cfg,
		log:  log,
	}

	err = clnt.dial()
	if err != nil {
		return nil, err
	}

	return clnt, nil
}

func (c *Client) dial() (err error) {
	c.temporalClient, err = client.Dial(client.Options{
		HostPort: c.conf.TemporalServiceAddress,
		Logger:   c.log,
	})
	if err != nil {
		c.log.Errorf("failed to dail temporal", err)
		return err
	}
	return nil
}

func (c *Client) Close() {
	c.temporalClient.Close()
}

func (c *Client) ExecuteWorkflow(
	ctx context.Context,
	options client.StartWorkflowOptions,
	workflowName string,
	args ...interface{},
) (client.WorkflowRun, error) {
	return c.temporalClient.ExecuteWorkflow(ctx, options, workflowName, args)
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
