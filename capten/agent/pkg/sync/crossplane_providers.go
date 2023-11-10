package sync

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "github.com/crossplane/crossplane/apis/pkg/v1"
	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"

	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type FetchCrossPlaneProviders struct {
	log    logging.Logger
	client *k8s.K8SClient
	db     *captenstore.Store
	creds  credentials.CredentialAdmin
}

func NewFetchCrossPlaneProviders() (*FetchCrossPlaneProviders, error) {
	log := logging.NewLogger()
	db, err := captenstore.NewStore(log)
	if err != nil {
		// ignoring store failure until DB user creation working
		// return nil, err
		log.Errorf("failed to initialize store, %v", err)
	}

	k8sclient, err := k8s.NewK8SClient(log)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	credAdmin, err := credentials.NewCredentialAdmin(context.TODO())
	if err != nil {
		log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client")
		return nil, err
	}

	// avlClusters, err := getManagedClusterEndpointMap(db)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to execute  getManagedClusterEndpointMap, err: %v", err)
	// }

	return &FetchCrossPlaneProviders{log: log, client: k8sclient, db: db, creds: credAdmin}, nil
}

func (fetch *FetchCrossPlaneProviders) Run() {
	fetch.log.Info("started to sync CrossplaneProvider resources")

	objList, err := fetch.client.DynamicClient.ListAllNamespaceResource(context.TODO(), schema.GroupVersionResource{Group: v1.Group, Version: v1.Version, Resource: "providers"})
	if err != nil {
		fetch.log.Error("Failed to fetch all the resource, err:", err)

		return
	}

	providers, err := json.Marshal(objList)
	if err != nil {
		fetch.log.Error("Failed to marshall the data, err:", err)

		return
	}

	var providerObj v1.ProviderList
	err = json.Unmarshal(providers, &providerObj)
	if err != nil {
		fetch.log.Error("Failed to un-marshall the data, err:", err)

		return
	}

	fetch.UpdateCrossplaneProvider(providerObj.Items)

	fetch.log.Info("succesfully sync-ed CrossplaneProvider resources")
}

func (fetch *FetchCrossPlaneProviders) UpdateCrossplaneProvider(clObj []v1.Provider) {
	for _, obj := range clObj {
		for _, status := range obj.Status.Conditions {
			if status.Type != v1.TypeHealthy {
				continue
			}

			provider := model.CrossplaneProvider{}
			provider.ProviderName = obj.Name
			provider.Status = string(v1.TypeHealthy)
			//fetch.db.UpdateCrossplaneProvider(&provider)
		}
	}
}
