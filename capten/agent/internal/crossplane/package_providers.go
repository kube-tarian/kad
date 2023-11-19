package crossplane

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"

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
	log     logging.Logger
	dbStore *captenstore.Store
}

func NewProvidersSyncHandler(log logging.Logger, dbStore *captenstore.Store) *ProvidersSyncHandler {
	return &ProvidersSyncHandler{log: log, dbStore: dbStore}
}

func RegisterK8SProviderWatcher(log logging.Logger, dbStore *captenstore.Store, dynamicClient dynamic.Interface) error {
	return k8s.RegisterDynamicInformers(NewProvidersSyncHandler(log, dbStore), dynamicClient, pgvk)
}

func getProviderObj(obj any) *model.Provider {
	clusterClaimByte, err := json.Marshal(obj)
	if err != nil {
		return nil
	}

	var clObj model.Provider
	err = json.Unmarshal(clusterClaimByte, &clObj)
	if err != nil {
		return nil
	}

	return &clObj
}

func (h *ProvidersSyncHandler) OnAdd(obj interface{}) {
	newCcObj := getProviderObj(obj)
	if newCcObj == nil {
		return
	}
	if err := h.updateCrossplaneProvider([]model.Provider{*newCcObj}); err != nil {
		return
	}

	h.log.Info("Crossplane Provider resources synched")
}

func (h *ProvidersSyncHandler) OnUpdate(oldObj, newObj interface{}) {
	prevObj := getProviderObj(oldObj)
	if prevObj == nil {
		return
	}

	newCcObj := getProviderObj(oldObj)
	if newCcObj == nil {
		return
	}

	// We receive the objects details on configured interval, identify actual updates made on the obj.
	if newCcObj.ResourceVersion == newCcObj.ResourceVersion {
		return
	}

	if err := h.updateCrossplaneProvider([]model.Provider{*newCcObj}); err != nil {
		return
	}

	h.log.Info("Crossplane Provider resources synched")
}

func (h *ProvidersSyncHandler) OnDelete(obj interface{}) {}

func (h *ProvidersSyncHandler) Sync() error {
	h.log.Debug("started to sync CrossplaneProvider resources")

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
		for _, providerStatus := range k8sProvider.Status.Conditions {
			if providerStatus.Type != model.TypeHealthy {
				continue
			}

			dbProvider, ok := dbProviderMap[k8sProvider.Name]
			if !ok {
				h.log.Infof("Provider name %s is not found in the db, skipping the update", k8sProvider.Name)
				continue
			}

			provider := model.CrossplaneProvider{
				Id:              dbProvider.Id,
				Status:          string(providerStatus.Type),
				CloudType:       dbProvider.CloudType,
				CloudProviderId: dbProvider.CloudProviderId,
				ProviderName:    dbProvider.ProviderName,
			}

			if err := h.dbStore.UpdateCrossplaneProvider(&provider); err != nil {
				h.log.Errorf("failed to update provider %s details in db, %v", k8sProvider.Name, err)
				continue
			}
			h.log.Infof("updated the crossplane provider %s", k8sProvider.Name)
		}
	}
	return nil
}
