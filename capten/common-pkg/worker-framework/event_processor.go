package workerframework

import (
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/temporalclient"

	"go.temporal.io/sdk/worker"
)

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
	logger         logging.Logger
}

func NewWorker(taskQueueName string, wf, activity interface{}, logger logging.Logger) (*Worker, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}

	worker := &Worker{
		conf:   cfg,
		logger: logger,
	}

	err = worker.RegisterToTemporal(taskQueueName, wf, activity)
	if err != nil {
		return nil, err
	}

	return worker, nil
}

func NewWorkerV2(taskQueueName string, logger logging.Logger) (*Worker, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}

	worker := &Worker{
		conf:   cfg,
		logger: logger,
	}

	err = worker.createWorker(taskQueueName)
	if err != nil {
		return nil, err
	}

	return worker, nil
}

func (w *Worker) RegisterToTemporal(taskQueueName string, wf, activity interface{}) error {
	// The client and worker are heavyweight objects that should be created once per process.
	err := w.createWorker(taskQueueName)
	if err != nil {
		return err
	}
	w.temporalWorker.RegisterWorkflow(wf)
	w.temporalWorker.RegisterActivity(activity)

	return nil
}

func (w *Worker) RegisterWorkflows(wfList ...interface{}) {
	for _, wf := range wfList {
		w.temporalWorker.RegisterWorkflow(wf)
	}
}

func (w *Worker) RegisterActivities(activityList ...interface{}) {
	for _, a := range activityList {
		w.temporalWorker.RegisterActivity(a)
	}
}

func (w *Worker) createWorker(taskQueueName string) (err error) {
	w.temporalClient, err = temporalclient.NewClient(w.logger)
	if err != nil {
		return fmt.Errorf("unable to create client, %v", err)
	}

	err = w.temporalClient.RegisterNamespace()
	if err != nil {
		return fmt.Errorf("default namespace create verification failed, %v", err)
	}

	w.temporalWorker = worker.New(w.temporalClient.TemporalClient, taskQueueName, worker.Options{
		MaxConcurrentWorkflowTaskExecutionSize: 2,
		MaxConcurrentActivityExecutionSize:     4,
	})
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
