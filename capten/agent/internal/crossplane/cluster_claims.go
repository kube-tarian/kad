package crossplane

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"
	"github.com/kube-tarian/kad/capten/agent/internal/temporalclient"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"

	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var (
	readyStatusType          = "Ready"
	clusterNotReadyStatus    = "NotReady"
	clusterReadyStatus       = "Ready"
	readyStatusValue         = "True"
	NorReadyStatusValue      = "False"
	clusterSecretName        = "%s-cluster"
	kubeConfig               = "kubeconfig"
	k8sEndpoint              = "endpoint"
	k8sClusterCA             = "clusterCA"
	managedClusterEntityName = "managedcluster"
)

var (
	cgvk = schema.GroupVersionResource{Group: "prodready.cluster", Version: "v1alpha1", Resource: "clusterclaims"}
)

type ClusterClaimSyncHandler struct {
	log     logging.Logger
	tc      *temporalclient.Client
	dbStore *captenstore.Store
}

func NewClusterClaimSyncHandler(log logging.Logger, dbStore *captenstore.Store) (*ClusterClaimSyncHandler, error) {
	tc, err := temporalclient.NewClient(log)
	if err != nil {
		return nil, err
	}

	return &ClusterClaimSyncHandler{log: log, dbStore: dbStore, tc: tc}, nil
}

func RegisterK8SClusterClaimWatcher(log logging.Logger, dbStore *captenstore.Store, dynamicClient dynamic.Interface) error {
	obj, err := NewClusterClaimSyncHandler(log, dbStore)
	if err != nil {
		return err
	}
	return k8s.RegisterDynamicInformers(obj, dynamicClient, cgvk)
}

func getClusterClaimObj(obj any) (*model.ClusterClaim, error) {
	clusterClaimByte, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var clObj model.ClusterClaim
	err = json.Unmarshal(clusterClaimByte, &clObj)
	if err != nil {
		return nil, err
	}

	return &clObj, nil
}

func (h *ClusterClaimSyncHandler) OnAdd(obj interface{}) {
	h.log.Info("Crossplane ClusterCliam Add Callback")

	newCcObj, err := getClusterClaimObj(obj)
	if newCcObj == nil {
		h.log.Errorf("failed to read ClusterCliam object, %v", err)
		return
	}

	if err = h.updateManagedClusters([]model.ClusterClaim{*newCcObj}); err != nil {
		h.log.Errorf("failed to update ClusterCliam object, %v", err)
		return
	}

	h.log.Info("cluster-claims resources synched")
}

func (h *ClusterClaimSyncHandler) OnUpdate(oldObj, newObj interface{}) {
	h.log.Info("Crossplane ClusterCliam Update Callback")

	prevObj, err := getClusterClaimObj(oldObj)
	if prevObj == nil {
		h.log.Errorf("failed to read ClusterCliam old object %v", err)
		return
	}

	newCcObj, err := getClusterClaimObj(oldObj)
	if newCcObj == nil {
		h.log.Errorf("failed to read ClusterCliam new object %v", err)
		return
	}

	if err = h.updateManagedClusters([]model.ClusterClaim{*newCcObj}); err != nil {
		h.log.Errorf("failed to update ClusterCliam object, %v", err)
		return
	}

	h.log.Info("cluster-claims resources synched")
}

func (h *ClusterClaimSyncHandler) OnDelete(obj interface{}) {
	h.log.Info("Crossplane ClusterCliam Delete Callback")
}

func (h *ClusterClaimSyncHandler) Sync() error {
	h.log.Debug("started to sync ClusterCliam resources")

	k8sclient, err := k8s.NewK8SClient(h.log)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	objList, err := k8sclient.DynamicClient.ListAllNamespaceResource(context.TODO(), cgvk)
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

	if err = h.updateManagedClusters(clObj.Items); err != nil {
		return fmt.Errorf("failed to update clusters in DB, %v", err)
	}
	h.log.Info("cluster-claims resources synched")
	return nil
}

func (h *ClusterClaimSyncHandler) updateManagedClusters(clusterCliams []model.ClusterClaim) error {
	_, err := k8s.NewK8SClient(h.log)
	if err != nil {
		return fmt.Errorf("failed to get k8s client, %v", err)
	}

	clusters, err := h.getManagedClusters()
	if err != nil {
		return fmt.Errorf("failed to get managed clusters from DB, %v", err)
	}

	for _, clusterCliam := range clusterCliams {
		h.log.Infof("processing cluster claim %s", clusterCliam.Metadata.Name)
		for _, status := range clusterCliam.Status.Conditions {
			if status.Type != readyStatusType {
				continue
			}

			managedCluster := &captenpluginspb.ManagedCluster{}
			managedCluster.ClusterName = clusterCliam.Metadata.Name

			clusterObj, ok := clusters[managedCluster.ClusterName]
			if !ok {
				managedCluster.Id = uuid.New().String()
			} else {
				h.log.Infof("found existing managed clusterId %s, updating", clusterObj.Id)
				managedCluster.Id = clusterObj.Id
				managedCluster.ClusterDeployStatus = clusterObj.ClusterDeployStatus
			}

			if status.Status != readyStatusValue {
				// secretName := fmt.Sprintf(clusterSecretName, clusterCliam.Spec.Id)
				// resp, err := k8sclient.GetSecretData(clusterCliam.Metadata.Namespace, secretName)
				// if err != nil {
				// 	h.log.Errorf("failed to get secret %s/%s, %v", clusterCliam.Metadata.Namespace, secretName, err)
				// 	continue
				// }

				resp := map[string]string{}
				clusterEndpoint := resp[k8sEndpoint]
				managedCluster.ClusterEndpoint = clusterEndpoint
				cred := map[string]string{}
				cred[kubeConfig] = resp[kubeConfig]
				cred[k8sClusterCA] = resp[k8sClusterCA]
				cred[k8sEndpoint] = clusterEndpoint

				err = credential.PutGenericCredential(context.TODO(), managedClusterEntityName, managedCluster.Id, cred)
				if err != nil {
					h.log.Errorf("failed to store credential for %s, %v", managedCluster.Id, err)
					continue
				}

				managedCluster.ClusterDeployStatus = clusterReadyStatus
				// call config-worker.
				err = h.UpdateClusterEndpoint(managedCluster.Id, clusterCliam.Metadata.Namespace, clusterCliam.Spec.Id, clusterEndpoint, resp[kubeConfig])
				if err != nil {
					h.log.Info("failed to update cluster endpoint information %v", err)
					continue
				}
			} else {
				managedCluster.ClusterDeployStatus = clusterNotReadyStatus
			}

			err = h.dbStore.UpsertManagedCluster(managedCluster)
			if err != nil {
				h.log.Info("failed to update information to db, %v", err)
				continue
			}
			h.log.Infof("updated the cluster claim %s with status %s", managedCluster.ClusterName, managedCluster.ClusterDeployStatus)
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
		clusterEndpointMap[cluster.ClusterName] = cluster
	}
	return clusterEndpointMap, nil
}

func (h *ClusterClaimSyncHandler) UpdateClusterEndpoint(id, ns, name, endpoint, k8sConfig string) error {
	proj, err := h.dbStore.GetCrossplaneProject()
	if err != nil {
		return err
	}
	ci := model.CrossplaneClusterEndpoint{RepoURL: proj.GitProjectUrl, Id: proj.Id, Name: name, Endpoint: endpoint, Namespace: ns, Kubeconfig: k8sConfig}

	wd := workers.NewConfig(h.tc, h.log)

	_, err = wd.SendEvent(context.TODO(), &model.ConfigureParameters{Resource: "crossplane", Action: "configClusterEndpoint"}, ci)
	if err != nil {
		return err
	}

	return nil
}
