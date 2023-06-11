package temporalclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/intelops/go-common/logging"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type Configuration struct {
	TemporalServiceAddress  string `envconfig:"TEMPORAL_SERVICE_URL" default:"localhost:7233"`
	EncryptionKey           string `envconfig:"TEMPORAL_ENCRYPTIONKEY" default:"00000000~secretGoesHere~00000000"`
	Secure                  bool   `envconfig:"TEMPORAL_SECURE" default:"false"`
	NamespaceTasksRetention int    `envconfig:"TEMPORAL_NAMESPACE_TASKS_RETENTION" default:"3"`
	RetryCount              int    `envconfig:"TEMPORAL_RETRY_COUNT" default:"3"`
	RetryTimeout            int    `envconfig:"TEMPORAL_RETRY_TIMEOUT" default:"10"`
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
	return c.TemporalClient.ExecuteWorkflow(ctx, options, workflowName, args)
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}

func (c *Client) RegisterNamespace() error {
	// Below code is to register namespace
	// Reference: https://docs.temporal.io/application-development/features?lang=go#namespaces
	// Equivalent cli command: tctl --ns default namespace register -rd 3
	// Reference: https://docs.temporal.io/tctl-v1/namespace#register
	client, err := client.NewNamespaceClient(client.Options{HostPort: c.conf.TemporalServiceAddress})
	if err != nil {
		return fmt.Errorf("unable to create namespace client, %v", err)
	}

	var retention time.Duration = 3 * 24 * time.Hour
	err = client.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        "default",
		WorkflowExecutionRetentionPeriod: &retention,
	})
	if err != nil && !strings.Contains(err.Error(), "Namespace already exists") {
		return fmt.Errorf("unable to register namesapce, %v", err)
	}
	c.log.Infof("default namespace registered. Verifying whether reflected or not")
	for i := 0; i < 3; i++ {
		_, err = client.Describe(context.Background(), "default")
		if err != nil {
			c.log.Errorf("retrying, namesapce not found, %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	if err == nil {
		c.log.Infof("default namespace registered and verified")
	}
	return err
}
