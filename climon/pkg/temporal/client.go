package temporal

import (
	"log"

	"github.com/pkg/errors"

	tc "go.temporal.io/sdk/client"
	tw "go.temporal.io/sdk/worker"
)

type Client struct {
	options        tc.Options
	temporalClient tc.Client
	worker         tw.Worker
}

func NewClient(address string) (*Client, error) {
	options := tc.Options{
		HostPort: address,
	}

	temporalClient, err := tc.Dial(options)
	if err != nil {
		log.Println("failed to dail temporal", err)
		return nil, err
	}

	return &Client{
		options:        options,
		temporalClient: temporalClient,
	}, nil
}

func (c *Client) Close() {
	c.worker.Stop()
	c.temporalClient.Close()
}

func (c *Client) CreateWorker(taskQue string) {
	c.worker = tw.New(c.temporalClient, taskQue, tw.Options{})
}

func (c *Client) RegisterWorkflow(workflow interface{}) {
	c.worker.RegisterWorkflow(workflow)
}

func (c *Client) RegisterActivity(activity interface{}) {
	c.worker.RegisterActivity(activity)
}

func (c *Client) StartWorker() error {
	log.Println("starting the worker")
	err := c.worker.Run(tw.InterruptCh())
	if err != nil {
		return errors.Wrapf(err, "unable to start the worker")
	}
	return nil
}
