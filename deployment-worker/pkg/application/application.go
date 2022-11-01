package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/deployment-worker/pkg/activities"
	"github.com/kube-tarian/kad/deployment-worker/pkg/handler"
	"github.com/kube-tarian/kad/deployment-worker/pkg/workflows"
	workerframework "github.com/kube-tarian/kad/pkg/worker-framework"
)

const (
	WorkflowName = "Deployment"
)

type Configuration struct {
	Port int `envconfig:"PORT" default:"9080"`
}

type Application struct {
	conf       *Configuration
	apiServer  *handler.APIHandler
	httpServer *http.Server
	worker     *workerframework.Worker
}

func New() *Application {
	cfg := &Configuration{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Could not parse env Config: %v\n", err)
	}

	// TODO: Create Worker instance and store in Handler
	worker, err := workerframework.NewWorker(WorkflowName, workflows.Workflow, &activities.Activities{})
	if err != nil {
		log.Fatalf("Worker initialization failed, Reason: %v\n", err)
	}

	apiServer, err := handler.NewAPIHandler(worker)
	if err != nil {
		log.Fatalf("API Handler initialisation failed: %v\n", err)
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
	}
}

func (app *Application) Start() {
	log.Printf("Starting worker\n")
	go func() {
		err := app.worker.Run()
		if err != nil {
			app.Close()
			log.Fatalf("Worker stopped listening on temporal, exiting. Readon: %v\n", err)
		}
	}()

	log.Printf("Starting server at %v\n", app.httpServer.Addr)
	var err error
	if err = app.httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Unexpected server close: %v", err)
	}
	log.Fatalf("Server closed")
}

func (app *Application) Close() {
	log.Printf("Closing the service gracefully\n")
	app.worker.Close()

	if err := app.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("Could not close the service gracefully: %v\n", err)
	}
}
