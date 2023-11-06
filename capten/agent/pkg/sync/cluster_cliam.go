package sync

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"

	pb "github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"

	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Fetch struct {
	log    logging.Logger
	client *k8s.K8SClient
	db     *captenstore.Store
}

func NewFetch() (*Fetch, error) {
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

	return &Fetch{log: log, client: k8sclient, db: db}, nil
}

// Run ...
func (fetch *Fetch) Run() {
	fetch.log.Info("started to sync cluster-claims resources")

	objList, err := fetch.client.DynamicClient.ListAllNamespaceResource(context.TODO(), schema.GroupVersionResource{Group: "prodready.cluster", Version: "v1alpha1", Resource: "clusterclaims"})
	if err != nil {
		fetch.log.Error("Failed to fetch all the resource, err:", err)

		return
	}

	clusterClaimByte, err := json.Marshal(objList)
	if err != nil {
		fetch.log.Error("Failed to marshall the data, err:", err)

		return
	}

	var clObj model.ClusterClaimList
	err = json.Unmarshal(clusterClaimByte, &clObj)
	if err != nil {
		fetch.log.Error("Failed to un-marshall the data, err:", err)

		return
	}

	err = fetch.UpdateClusterDetails(clObj.Items)
	if err != nil {
		fetch.log.Error("Failed to UpdateClusterDetails, err:", err)

		return
	}

	fetch.log.Info("succesfully sync-ed cluster-claims resources")
}

func (fetch *Fetch) UpdateClusterDetails(clObj []model.ClusterClaim) error {
	for _, obj := range clObj {
		managedCluster := &pb.ManagedCluster{}
		managedCluster.Id = obj.Spec.Id
		// add code to fetch the kubeconfig file
		//updated the cluster name and other details
		// err := fetch.db.UpsertManagedCluster(managedCluster)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}
