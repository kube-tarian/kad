package temporalclient

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"

	"go.temporal.io/sdk/client"
)

type Configuration struct {
	TemporalServiceAddress string `envconfig:"TEMPORAL_SERVICE_URL" default:"localhost:7233"`
	EncryptionKey          string `envconfig:"ENCRYPTIONKEY" default:"00000000~secretGoesHere~00000000"`
	Secure                 bool   `envconfig:"SECURE" default:"false"`
}

type Client struct {
	conf           *Configuration
	TemporalClient client.Client
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

	err = clnt.newTemporalClient()
	if err != nil {
		return nil, err
	}

	return clnt, nil
}

func (c *Client) newTemporalClient() (err error) {
	opts := client.Options{
		Namespace: "default",
		HostPort:  c.conf.TemporalServiceAddress,
		Logger:    c.log,
	}

	if c.conf.Secure {
		encryptedDataConverter, err := NewEncryptDataConverterV1(Options{
			EncryptionKey: []byte(c.conf.EncryptionKey),
		})
		if err != nil {
			return err
		}

		opts.DataConverter = encryptedDataConverter
	}
	c.TemporalClient, err = client.Dial(opts)
	return
}

func (c *Client) Close() {
	c.TemporalClient.Close()
}

func (c *Client) ExecuteWorkflow(
	ctx context.Context,
	options client.StartWorkflowOptions,
	workflowName string,
	args ...interface{},
) (client.WorkflowRun, error) {
	return c.TemporalClient.ExecuteWorkflow(ctx, options, workflowName, args...)
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
