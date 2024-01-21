package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"
	"github.com/kube-tarian/kad/capten/agent/internal/temporalclient"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var (
	pgvk = schema.GroupVersionResource{Group: "pkg.crossplane.io", Version: "v1", Resource: "providers"}
)

type ProvidersSyncHandler struct {
	log             logging.Logger
	tc              *temporalclient.Client
	dbStore         *captenstore.Store
	activeProviders map[string]bool
	mutex           sync.Mutex
}

func NewProvidersSyncHandler(log logging.Logger, dbStore *captenstore.Store) (*ProvidersSyncHandler, error) {
	tc, err := temporalclient.NewClient(log)
	if err != nil {
		return nil, err
	}

	return &ProvidersSyncHandler{log: log, dbStore: dbStore, tc: tc, activeProviders: map[string]bool{}}, nil
}

func registerK8SProviderWatcher(log logging.Logger, dbStore *captenstore.Store, dynamicClient dynamic.Interface) error {
	provider, err := NewProvidersSyncHandler(log, dbStore)
	if err != nil {
		return err
	}
	return k8s.RegisterDynamicInformers(provider, dynamicClient, pgvk)
}

func getProviderObj(obj any) (*model.Provider, error) {
	clusterClaimByte, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var clObj model.Provider
	err = json.Unmarshal(clusterClaimByte, &clObj)
	if err != nil {
		return nil, err
	}

	return &clObj, nil
}

func (h *ProvidersSyncHandler) OnAdd(obj interface{}) {
	h.log.Info("Crossplane Provider Add Callback")
	h.mutex.Lock()
	defer h.mutex.Unlock()

	newCcObj, err := getProviderObj(obj)
	if newCcObj == nil {
		h.log.Errorf("failed to read Provider object, %v", err)
		return
	}

	if err := h.updateCrossplaneProvider([]model.Provider{*newCcObj}); err != nil {
		h.log.Errorf("failed to update Provider object, %v", err)
		return
	}
}

func (h *ProvidersSyncHandler) OnUpdate(oldObj, newObj interface{}) {
	h.log.Debug("Crossplane Provider Update Callback")
	h.mutex.Lock()
	defer h.mutex.Unlock()

	newCcObj, err := getProviderObj(newObj)
	if newCcObj == nil {
		h.log.Errorf("failed to read Provider new object %v", err)
		return
	}

	if err := h.updateCrossplaneProvider([]model.Provider{*newCcObj}); err != nil {
		h.log.Errorf("failed to update Provider object, %v", err)
		return
	}
}

func (h *ProvidersSyncHandler) OnDelete(obj interface{}) {
	h.log.Info("Crossplane Provider Delete Callback")
	newCcObj, err := getProviderObj(obj)
	if newCcObj == nil {
		h.log.Errorf("failed to read Provider object, %v", err)
		return
	}
	delete(h.activeProviders, newCcObj.Name)
}

func (h *ProvidersSyncHandler) Sync() error {
	h.log.Debug("started to sync CrossplaneProvider resources")
	h.mutex.Lock()
	defer h.mutex.Unlock()

	k8sclient, err := k8s.NewK8SClient(h.log)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	objList, err := k8sclient.DynamicClient.ListAllNamespaceResource(context.TODO(), pgvk)
	if err != nil {
		return fmt.Errorf("failed to fetch providers resources, %v", err)
	}

	providers, err := json.Marshal(objList)
	if err != nil {
		return fmt.Errorf("failed to marshall the data, %v", err)
	}

	var providerObj model.ProviderList
	err = json.Unmarshal(providers, &providerObj)
	if err != nil {
		return fmt.Errorf("failed to un-marshall the data, %s", err)
	}

	if err = h.updateCrossplaneProvider(providerObj.Items); err != nil {
		return fmt.Errorf("failed to update providers in DB, %v", err)
	}
	h.log.Debug("Crossplane Provider resources synched")
	return nil
}

func (h *ProvidersSyncHandler) updateCrossplaneProvider(k8sProviders []model.Provider) error {
	dbProviders, err := h.dbStore.GetCrossplaneProviders()
	if err != nil {
		return fmt.Errorf("failed to get Crossplane Providers, %v", err)
	}

	dbProviderMap := make(map[string]*captenpluginspb.CrossplaneProvider)
	for _, dbProvider := range dbProviders {
		dbProviderMap[model.PrepareCrossplaneProviderName(dbProvider.CloudType)] = dbProvider
	}

	for _, k8sProvider := range k8sProviders {
		h.log.Debugf("processing Crossplane Provider %s", k8sProvider.Name)
		for _, providerStatus := range k8sProvider.Status.Conditions {
			if providerStatus.Type != model.TypeHealthy {
				continue
			}

			dbProvider, ok := dbProviderMap[k8sProvider.Name]
			if !ok {
				h.log.Infof("Provider name %s is not found in the db, skipping the update", k8sProvider.Name)
				delete(h.activeProviders, k8sProvider.Name)
				continue
			}

			_, ok = h.activeProviders[k8sProvider.Name]
			if !ok {
				h.log.Debugf("Provider name %s is already configured, skipping the update", k8sProvider.Name)
				continue
			}

			status := model.CrossPlaneProviderNotReady
			if strings.EqualFold(string(providerStatus.Status), "true") {
				status = model.CrossPlaneProviderReady
			}
			provider := model.CrossplaneProvider{
				Id:              dbProvider.Id,
				Status:          string(status),
				CloudType:       dbProvider.CloudType,
				CloudProviderId: dbProvider.CloudProviderId,
				ProviderName:    dbProvider.ProviderName,
			}

			if err := h.dbStore.UpdateCrossplaneProvider(&provider); err != nil {
				h.log.Errorf("failed to update provider %s details in db, %v", k8sProvider.Name, err)
				continue
			}
			h.log.Infof("updated the crossplane provider %s, status: %s", k8sProvider.Name, status)

			if status == model.CrossPlaneProviderReady {
				err = h.triggerProviderUpdate(provider.ProviderName, provider)
				if err != nil {
					return fmt.Errorf("failed to trigger crossplane provider update workflow, %v", err)
				}
				h.log.Infof("triggered crossplane provider update workflow for provider %s", provider.ProviderName)
			}
		}
	}
	return nil
}

func (h *ProvidersSyncHandler) triggerProviderUpdate(clusterName string, provider model.CrossplaneProvider) error {
	wd := workers.NewConfig(h.tc, h.log)

	proj, err := h.dbStore.GetCrossplaneProject()
	if err != nil {
		return err
	}
	ci := model.CrossplaneClusterUpdate{RepoURL: proj.GitProjectUrl, GitProjectId: proj.GitProjectId,
		ManagedClusterName: clusterName, ManagedClusterId: provider.Id}

	wkfId, err := wd.SendAsyncEvent(context.TODO(), &model.ConfigureParameters{Resource: model.CrossPlaneResource, Action: model.CrossPlaneProviderUpdate}, ci)
	if err != nil {
		return fmt.Errorf("failed to send event to crossplane provider update workflow to configure %s, %v", provider.ProviderName, err)
	}

	h.log.Infof("Crossplane provider update %s config workflow %s created", provider.ProviderName, wkfId)

	go h.monitorProviderUpdateWorkflow(&provider, wkfId)

	return nil
}

func (h *ProvidersSyncHandler) monitorProviderUpdateWorkflow(provider *model.CrossplaneProvider, wkfId string) {
	// during system reboot start monitoring, add it in map or somewhere.
	wd := workers.NewConfig(h.tc, h.log)
	_, err := wd.GetWorkflowInformation(context.TODO(), wkfId)
	if err != nil {
		h.log.Errorf("failed to send crossplane provider update event to workflow to configure %s, %v", provider.ProviderName, err)
		return
	}
	h.activeProviders[model.PrepareCrossplaneProviderName(provider.CloudType)] = true
	h.log.Infof("Crossplane provider update %s config workflow %s completed", provider.ProviderName, wkfId)
}
