package workerframework

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Plugin interface {
	Exec(payload interface{})
}

type Action interface {
	GetStatus()
}

type Configuration struct {
	TemporalServiceAddress string `envconfig:"TEMPORAL_SERVICE_URL" default:"http://localhost:7233"`
}

type Worker struct {
	conf           *Configuration
	temporalClient client.Client
	temporalWorker worker.Worker
	plugins        map[string]Plugin
}

func NewWorker(workflowName string, wf, activity interface{}) (*Worker, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}

	worker := &Worker{
		conf:    cfg,
		plugins: make(map[string]Plugin),
	}

	err = worker.RegisterToTemporal(workflowName, wf, activity)
	if err != nil {
		return nil, err
	}

	return worker, nil
}

func (w *Worker) RegisterPlugin(pluginName string, plugin Plugin) {
	w.plugins[pluginName] = plugin
}

func (w *Worker) DeregisterPlugin(pluginName string) {
	delete(w.plugins, pluginName)
}

func (w *Worker) RegisterToTemporal(workflowName string, wf, activity interface{}) (err error) {
	// The client and worker are heavyweight objects that should be created once per process.
	w.temporalClient, err = client.Dial(client.Options{})
	if err != nil {
		return fmt.Errorf("unable to create client, %v", err)
	}

	w.temporalWorker = worker.New(w.temporalClient, workflowName, worker.Options{})
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
	log.Printf("Stopping temporal worker client\n")
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
