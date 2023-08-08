package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/climon/pkg/activities"
	"github.com/kube-tarian/kad/capten/climon/pkg/handler"
	"github.com/kube-tarian/kad/capten/climon/pkg/workflows"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
)

const (
	WorkflowTaskQueueName = "CLIMON_HELM_TASK_QUEUE"
	HelmPluginName        = "helm"
)

type Configuration struct {
	Port int `envconfig:"PORT" default:"9080"`
}

type Application struct {
	conf       *Configuration
	apiServer  *handler.APIHandler
	httpServer *http.Server
	worker     *workerframework.Worker
	logger     logging.Logger
}

func New(logger logging.Logger) *Application {
	cfg := &Configuration{}
	if err := envconfig.Process("", cfg); err != nil {
		logger.Fatalf("Could not parse env Config: %v\n", err)
	}

	worker, err := workerframework.NewWorker(WorkflowTaskQueueName, workflows.Workflow, &activities.Activities{}, logger)
	if err != nil {
		logger.Fatalf("Worker initialization failed, Reason: %v\n", err)
	}

	apiServer, err := handler.NewAPIHandler(worker)
	if err != nil {
		logger.Fatalf("API Handler initialisation failed: %v\n", err)
	}

	mux := chi.NewMux()
	apiServer.BindRequest(mux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: mux,
	}

	return &Application{
		conf:       cfg,
		apiServer:  apiServer,
		httpServer: httpServer,
		worker:     worker,
		logger:     logger,
	}
}

func (app *Application) Start() {
	app.logger.Infof("Starting worker\n")
	go func() {
		err := app.worker.Run()
		if err != nil {
			app.logger.Errorf("Worker stopped listening on temporal, exiting. Readon: %v\n", err)
			app.Close()
		}
	}()

	app.logger.Infof("Starting server at %v", app.httpServer.Addr)
	var err error
	if err = app.httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		app.logger.Fatalf("Unexpected server close: %v", err)
	}
	app.logger.Fatalf("Server closed")
}

func (app *Application) Close() {
	app.logger.Infof("Closing the service gracefully")
	app.worker.Close()

	if err := app.httpServer.Shutdown(context.Background()); err != nil {
		app.logger.Errorf("Could not close the service gracefully: %v", err)
	}
}
