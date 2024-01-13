package tekton

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"

	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var (
	pgvk = schema.GroupVersionResource{Group: "triggers.tekton.dev", Version: "v1beta1", Resource: "eventlisteners"}
)

type TektonPipelineSyncHandler struct {
	log     logging.Logger
	dbStore *captenstore.Store
}

func NewTektonPipelineSyncHandler(log logging.Logger, dbStore *captenstore.Store) *TektonPipelineSyncHandler {
	return &TektonPipelineSyncHandler{log: log, dbStore: dbStore}
}

func registerK8STektonPipelineSync(log logging.Logger, dbStore *captenstore.Store, dynamicClient dynamic.Interface) error {
	return k8s.RegisterDynamicInformers(NewTektonPipelineSyncHandler(log, dbStore), dynamicClient, pgvk)
}

func getEventListenerObj(obj any) (*model.EventListener, error) {
	clusterClaimByte, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var clObj model.EventListener
	err = json.Unmarshal(clusterClaimByte, &clObj)
	if err != nil {
		return nil, err
	}

	return &clObj, nil
}

func (h *TektonPipelineSyncHandler) OnAdd(obj interface{}) {
	h.log.Info("TektonPipeline Add Callback")
	newCcObj, err := getEventListenerObj(obj)
	if newCcObj == nil {
		h.log.Errorf("failed to read TektonPipeline object, %v", err)
		return
	}

	if err := h.updateTektonPipelines([]model.EventListener{*newCcObj}); err != nil {
		h.log.Errorf("failed to update TektonPipeline object, %v", err)
		return
	}
}

func (h *TektonPipelineSyncHandler) OnUpdate(oldObj, newObj interface{}) {
	h.log.Info("TektonPipeline Update Callback")
	prevObj, err := getEventListenerObj(oldObj)
	if prevObj == nil {
		h.log.Errorf("failed to read TektonPipeline old object %v", err)
		return
	}

	newCcObj, err := getEventListenerObj(oldObj)
	if newCcObj == nil {
		h.log.Errorf("failed to read TektonPipeline new object %v", err)
		return
	}

	if err := h.updateTektonPipelines([]model.EventListener{*newCcObj}); err != nil {
		h.log.Errorf("failed to update TektonPipeline object, %v", err)
		return
	}
}

func (h *TektonPipelineSyncHandler) OnDelete(obj interface{}) {
	h.log.Info("TektonPipeline Delete Callback")
}

func (h *TektonPipelineSyncHandler) Sync() error {
	h.log.Debug("started to sync TektonPipeline resources")

	k8sclient, err := k8s.NewK8SClient(h.log)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	objList, err := k8sclient.DynamicClient.ListAllNamespaceResource(context.TODO(), pgvk)
	if err != nil {
		return fmt.Errorf("failed to fetch pipelines resources, %v", err)
	}

	pipelines, err := json.Marshal(objList)
	if err != nil {
		return fmt.Errorf("failed to marshall the data, %v", err)
	}

	var pipelineObj model.EventListeners
	err = json.Unmarshal(pipelines, &pipelineObj)
	if err != nil {
		return fmt.Errorf("failed to un-marshall the data, %s", err)
	}

	if err = h.updateTektonPipelines(pipelineObj.Items); err != nil {
		return fmt.Errorf("failed to update TektonPipeline in DB, %v", err)
	}
	h.log.Debug("TektonPipeline resources synched")
	return nil
}

func (h *TektonPipelineSyncHandler) updateTektonPipelines(k8spipelines []model.EventListener) error {
	dbpipelines, err := h.dbStore.GetTektonPipeliness()
	if err != nil {
		return fmt.Errorf("failed to get TektonPipeline pipelines, %v", err)
	}

	dbpipelineMap := make(map[string]*model.TektonPipeline)
	for _, dbpipeline := range dbpipelines {
		dbpipelineMap[dbpipeline.PipelineName] = dbpipeline
	}

	for _, k8spipeline := range k8spipelines {
		h.log.Infof("processing TektonPipeline %s", k8spipeline.Name)
		for _, pipelineStatus := range k8spipeline.Status.Conditions {
			if pipelineStatus.Type != "Ready" {
				continue
			}

			dbpipeline, ok := dbpipelineMap[k8spipeline.Name]
			if !ok {
				h.log.Infof("TektonPipeline name %s is not found in the db, skipping the update", k8spipeline.Name)
				continue
			}

			status := model.TektonPipelineNotReady
			if strings.EqualFold(string(pipelineStatus.Status), "true") {
				status = model.TektonPipelineReady
			}

			dbpipeline.Status = string(status)

			v, _ := json.Marshal(dbpipeline)
			fmt.Println("TektonPipeline ===>" + string(v))

			if err := h.dbStore.UpsertTektonPipelines(dbpipeline); err != nil {
				h.log.Errorf("failed to update TektonPipeline %s details in db, %v", k8spipeline.Name, err)
				continue
			}
			h.log.Infof("updated the TektonPipeline eventlistener %s", k8spipeline.Name)
		}
	}
	return nil
}
