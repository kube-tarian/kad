package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"

	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	statusType               = "Ready"
	statusValue              = "True"
	clusterName              = "%s-%s"
	clusterSecretName        = "%s-cluster"
	kubeConfig               = "kubeconfig"
	k8sEndpoint              = "endpoint"
	k8sClusterCA             = "clusterCA"
	managedClusterEntityName = "managedcluster"
)

type ClusterClaimSyncHandler struct {
	log      logging.Logger
	dbStore  *captenstore.Store
	clusters map[string]*captenpluginspb.ManagedCluster
}

func NewClusterClaimSyncHandler(log logging.Logger, dbStore *captenstore.Store) (*ClusterClaimSyncHandler, error) {
	return &ClusterClaimSyncHandler{log: log, dbStore: dbStore, clusters: map[string]*captenpluginspb.ManagedCluster{}}, nil
}

func (h *ClusterClaimSyncHandler) Sync() error {
	h.log.Debug("started to sync cluster-claims resources")

	k8sclient, err := k8s.NewK8SClient(h.log)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	objList, err := k8sclient.DynamicClient.ListAllNamespaceResource(context.TODO(),
		schema.GroupVersionResource{Group: "prodready.cluster", Version: "v1alpha1", Resource: "clusterclaims"})
	if err != nil {
		return fmt.Errorf("failed to list cluster claim resources, %v", err)
	}

	clusterClaimByte, err := json.Marshal(objList)
	if err != nil {
		return fmt.Errorf("failed to marshal cluster claim resources, %v", err)
	}

	var clObj model.ClusterClaimList
	err = json.Unmarshal(clusterClaimByte, &clObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cluster claim resources, %v", err)
	}

	if err = h.updateManagedClusters(k8sclient, clObj.Items); err != nil {
		return fmt.Errorf("failed to update clusters in DB, %v", err)
	}
	h.log.Info("cluster-claims resources synched")
	return nil
}

func (h *ClusterClaimSyncHandler) updateManagedClusters(k8sClient *k8s.K8SClient, clObj []model.ClusterClaim) error {
	credAdmin, err := credentials.NewCredentialAdmin(context.TODO())
	if err != nil {
		return err
	}

	clusters, err := h.getManagedClusters()
	if err != nil {
		return fmt.Errorf("failed to get managed clusters from DB, %v", err)
	}

	for _, obj := range clObj {
		for _, status := range obj.Status.Conditions {
			if status.Type != statusType {
				continue
			}

			if status.Status != statusValue {
				h.log.Info("%s in namespace %s, status is %s, so skiping update to db.",
					obj.Spec.Id, obj.Metadata.Namespace, status.Status)
				continue
			}

			secretName := fmt.Sprintf(clusterSecretName, obj.Spec.Id)
			resp, err := k8sClient.GetSecretData(obj.Metadata.Namespace, secretName)
			if err != nil {
				h.log.Info("failed to get secret %s in namespace %s, %v",
					secretName, obj.Metadata.Namespace, err)
				continue
			}

			clusterEndpoint := resp.Data[k8sEndpoint]

			managedCluster := &captenpluginspb.ManagedCluster{}
			clusterObj, ok := h.clusters[clusterEndpoint]
			if !ok {
				managedCluster.Id = uuid.New().String()
			} else {
				h.log.Info("found existing Id: %s, updating the latest information ", clusterObj.Id)
				managedCluster.Id = clusterObj.Id
			}

			managedCluster.ClusterName = fmt.Sprintf(clusterName, obj.Spec.Id, obj.Metadata.Namespace)
			managedCluster.ClusterDeployStatus = status.Status
			managedCluster.ClusterEndpoint = clusterEndpoint

			clusterDetails := map[string]string{}
			clusterDetails[kubeConfig] = resp.Data[kubeConfig]
			clusterDetails[k8sClusterCA] = resp.Data[k8sClusterCA]

			err = credAdmin.PutCredential(context.TODO(), credentials.GenericCredentialType, managedClusterEntityName, managedCluster.Id, clusterDetails)

			if err != nil {
				h.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", managedCluster.Id)
				h.log.Errorf("failed to store credential for %s, %v", managedCluster.Id, err)
				continue
			}

			managedCluster.LastUpdateTime = time.Now().Format(time.RFC3339)
			err = h.dbStore.UpsertManagedCluster(managedCluster)
			if err != nil {
				h.log.Info("failed to update information to db, %v", err)
				continue
			}
			clusters[clusterEndpoint] = managedCluster
		}
	}
	return nil
}

func (h *ClusterClaimSyncHandler) getManagedClusters() (map[string]*captenpluginspb.ManagedCluster, error) {
	clusters, err := h.dbStore.GetManagedClusters()
	if err != nil {
		return nil, err
	}

	clusterEndpointMap := map[string]*captenpluginspb.ManagedCluster{}
	for _, cluster := range clusters {
		clusterEndpointMap[cluster.ClusterEndpoint] = cluster
	}
	return clusterEndpointMap, nil
}
