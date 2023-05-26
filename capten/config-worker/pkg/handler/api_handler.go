package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/capten/config-worker/api"
)

type APIHandler struct {
	worker *workerframework.Worker
}

const (
	appJSONContentType = "application/json"
	contentType        = "Content-Type"
)

func NewAPIHandler(worker *workerframework.Worker) (*APIHandler, error) {
	return &APIHandler{
		worker: worker,
	}, nil
}

func (ah *APIHandler) BindRequest(mux *chi.Mux) {
	mux.Route("/", func(r chi.Router) {
		api.HandlerFromMux(ah, r)
	})
}

func (ah *APIHandler) GetApiDocs(w http.ResponseWriter, r *http.Request) {
	swagger, err := api.GetSwagger()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set(contentType, appJSONContentType)
	_ = json.NewEncoder(w).Encode(swagger)
}

func (ah *APIHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, appJSONContentType)
	w.WriteHeader(http.StatusOK)
}
