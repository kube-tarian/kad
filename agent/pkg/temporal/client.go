package temporal

import (
	"context"
	"log"

	tc "go.temporal.io/sdk/client"
)

type client struct {
	options        tc.Options
	temporalClient tc.Client
}

func NewClient(address string) *client {
	return &client{
		options: tc.Options{
			HostPort: address,
		},
	}
}

func (c *client) Dail() (tc.Client, error) {
	temporalClient, err := tc.Dial(c.options)
	if err != nil {
		log.Println("failed to dail temporal", err)
		return nil, err
	}
	c.temporalClient = temporalClient
	return c.temporalClient, nil
}

func (c *client) Close() {
	c.temporalClient.Close()
}

func (c *client) ExecuteWorkflow(ctx context.Context, options tc.StartWorkflowOptions, workflowName string,
	args ...interface{}) (tc.WorkflowRun, error) {
	return c.temporalClient.ExecuteWorkflow(context.Background(), options, workflowName, args)
}
