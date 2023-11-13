package crossplane

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	providerNamePrefix = "provider-"
)

type ProvidersSyncHandler struct {
	log     logging.Logger
	client  *k8s.K8SClient
	dbStore *captenstore.Store
	creds   credentials.CredentialAdmin
}

func NewProvidersSyncHandler(log logging.Logger, dbStore *captenstore.Store) (*ProvidersSyncHandler, error) {
	return &ProvidersSyncHandler{log: log, dbStore: dbStore}, nil
}

func (h *ProvidersSyncHandler) Sync() error {
	h.log.Debug("started to sync CrossplaneProvider resources")

	k8sclient, err := k8s.NewK8SClient(h.log)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	objList, err := k8sclient.DynamicClient.ListAllNamespaceResource(context.TODO(),
		schema.GroupVersionResource{Group: "pkg.crossplane.io", Version: "v1", Resource: "providers"})
	if err != nil {
		return fmt.Errorf("failed to fetch providers resources, %v", err)
	}

	providers, err := json.Marshal(objList)
	if err != nil {
		return fmt.Errorf("failed to marshall the data, err:", err)
	}

	var providerObj model.ProviderList
	err = json.Unmarshal(providers, &providerObj)
	if err != nil {
		return fmt.Errorf("failed to un-marshall the data, err:", err)
	}

	if err = h.updateCrossplaneProvider(providerObj.Items); err != nil {
		return fmt.Errorf("failed to update providers in DB, %v", err)
	}
	h.log.Debug("Crossplane Provider resources synched")
	return nil
}

func (h *ProvidersSyncHandler) updateCrossplaneProvider(clObj []model.Provider) error {
	prvList, err := h.dbStore.GetCrossplaneProviders()
	if err != nil {
		return fmt.Errorf("failed to get Crossplane Providers, %v", err)
	}

	prvMap := make(map[string]*captenpluginspb.CrossplaneProvider)
	for _, prov := range prvList {
		prvMap[providerNamePrefix+prov.ProviderName] = prov
	}

	for _, obj := range clObj {
		for _, status := range obj.Status.Conditions {
			if status.Type != model.TypeHealthy {
				continue
			}

			prvObj, ok := prvMap[obj.Name]
			if !ok {
				h.log.Infof("Provider name %s is not found in the db, skipping the update", obj.Name)
				continue
			}

			provider := model.CrossplaneProvider{
				Id:              prvObj.Id,
				Status:          string(status.Type),
				CloudType:       prvObj.CloudType,
				CloudProviderId: prvObj.CloudProviderId,
				ProviderName:    prvObj.ProviderName,
			}

			if err := h.dbStore.UpdateCrossplaneProvider(&provider); err != nil {
				h.log.Errorf("failed to update provider %s details in db, %v", prvObj.ProviderName, err)
				continue
			}
			h.log.Infof("successfully updated the details for %s", prvObj.ProviderName)
		}
	}
	return nil
}
