package workerframework

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/temporalclient"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Plugin interface {
	DeployActivities(payload interface{}) (json.RawMessage, error)
	// Create(payload interface{}) (json.RawMessage, error)
	// Delete(payload interface{}) (json.RawMessage, error)
	// List(payload interface{}) (json.RawMessage, error)

	// ConfigurationActivities(payload interface{}) (json.RawMessage, error)
	// ConfgiureTarget(payload interface{}) (json.RawMessage, error)
	// SetTarget(payload interface{}) (json.RawMessage, error)
	// SetDefaultTarget(payload interface{}) (json.RawMessage, error)
}

type DeploymentWorker interface {
	DeployActivities(payload interface{}) (json.RawMessage, error)
	// Create(payload interface{}) (json.RawMessage, error)
	// Delete(payload interface{}) (json.RawMessage, error)
	// List(payload interface{}) (json.RawMessage, error)
}

type ConfigurationWorker interface {
	// ConfigurationActivities(payload interface{}) (json.RawMessage, error)
	// ConfgiureTarget(payload interface{}) (json.RawMessage, error)
	// SetTarget(payload interface{}) (json.RawMessage, error)
	// SetDefaultTarget(payload interface{}) (json.RawMessage, error)
}

type Action interface {
	GetStatus()
}

type Configuration struct {
	TemporalServiceAddress string `envconfig:"TEMPORAL_SERVICE_URL" default:"localhost:7233"`
}

type Worker struct {
	conf           *Configuration
	temporalClient *temporalclient.Client
	temporalWorker worker.Worker
	plugins        map[string]Plugin
	logger         logging.Logger
}

func NewWorker(taskQueueName string, wf, activity interface{}, logger logging.Logger) (*Worker, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}

	worker := &Worker{
		conf:    cfg,
		plugins: make(map[string]Plugin),
		logger:  logger,
	}

	err = worker.RegisterToTemporal(taskQueueName, wf, activity)
	if err != nil {
		return nil, err
	}

	return worker, nil
}

func (w *Worker) RegisterToTemporal(taskQueueName string, wf, activity interface{}) (err error) {
	// The client and worker are heavyweight objects that should be created once per process.
	w.temporalClient, err = temporalclient.NewClient(w.logger)
	if err != nil {
		return fmt.Errorf("unable to create client, %v", err)
	}

	// Below code is to register namespace
	// Reference: https://docs.temporal.io/application-development/features?lang=go#namespaces
	// Equivalent cli command: tctl --ns default namespace register -rd 3
	// Reference: https://docs.temporal.io/tctl-v1/namespace#register
	client, err := client.NewNamespaceClient(client.Options{HostPort: "127.0.0.1:7233"})
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
	w.logger.Infof("default namespace registered. Verifying whether reflected or not")
	for i := 0; i < 3; i++ {
		_, err = client.Describe(context.Background(), "default")
		if err != nil {
			w.logger.Errorf("retrying, namesapce not found, %v", err)
			continue
		}
		break
	}
	w.logger.Infof("default namespace registered and verified")

	w.temporalWorker = worker.New(w.temporalClient.TemporalClient, taskQueueName, worker.Options{})
	w.temporalWorker.RegisterWorkflow(wf)
	w.temporalWorker.RegisterActivity(activity)

	return nil
}

func (w *Worker) Run() error {
	err := w.temporalWorker.Run(worker.InterruptCh())
	if err != nil {
		return fmt.Errorf("unable to start worker, %v", err)
	}
	return nil
}

func (w *Worker) Close() {
	w.temporalClient.Close()
	w.logger.Infof("Stopping temporal worker client\n")
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
