package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"

	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
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

func registerK8SProviderWatcher(log logging.Logger, dbStore *captenstore.Store, dynamicClient dynamic.Interface) error {
	return k8s.RegisterDynamicInformers(NewProvidersSyncHandler(log, dbStore), dynamicClient, pgvk)
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
	h.log.Info("Crossplane Provider Update Callback")
	prevObj, err := getProviderObj(oldObj)
	if prevObj == nil {
		h.log.Errorf("failed to read Provider old object %v", err)
		return
	}

	newCcObj, err := getProviderObj(oldObj)
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
}

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
		h.log.Infof("processing Crossplane Provider %s", k8sProvider.Name)
		for _, providerStatus := range k8sProvider.Status.Conditions {
			if providerStatus.Type != model.TypeHealthy {
				continue
			}

			dbProvider, ok := dbProviderMap[k8sProvider.Name]
			if !ok {
				h.log.Infof("Provider name %s is not found in the db, skipping the update", k8sProvider.Name)
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

			v, _ := json.Marshal(provider)
			fmt.Println("Provider ===>" + string(v))

			if err := h.dbStore.UpdateCrossplaneProvider(&provider); err != nil {
				h.log.Errorf("failed to update provider %s details in db, %v", k8sProvider.Name, err)
				continue
			}
			h.log.Infof("updated the crossplane provider %s", k8sProvider.Name)
		}
	}
	return nil
}
