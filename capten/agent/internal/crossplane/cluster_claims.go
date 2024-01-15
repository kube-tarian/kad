package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"
	"github.com/kube-tarian/kad/capten/agent/internal/temporalclient"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"

	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	managedcluster "github.com/kube-tarian/kad/capten/common-pkg/managed-cluster"
	"github.com/kube-tarian/kad/capten/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

const (
	readyStatusType = "ready"

	clusterNotReadyStatus       = "NotReady"
	clusterReadyStatus          = "Ready"
	clusterDeletingStatus       = "Deleting"
	clusterFailedToDeleteStatus = "FailedToDelete"
	readyStatusValue            = "True"
	NorReadyStatusValue         = "False"
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

func registerK8SClusterClaimWatcher(log logging.Logger, dbStore *captenstore.Store, dynamicClient dynamic.Interface) error {
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

	newCcObj, err := getClusterClaimObj(obj)
	if newCcObj == nil {
		h.log.Errorf("failed to read ClusterCliam object, %v", err)
		return
	}

	if err = h.deleteManagedClusters(*newCcObj); err != nil {
		h.log.Errorf("failed to delete ClusterCliam object, %v", err)
		return
	}

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

	if err := h.syncClusterClaimsWithDB(clObj.Items); err != nil {
		return fmt.Errorf("failed to sync clusters in DB, %v", err)
	}

	return nil
}

func (h *ClusterClaimSyncHandler) updateManagedClusters(clusterCliams []model.ClusterClaim) error {
	clusters, err := h.getManagedClusters()
	if err != nil {
		return fmt.Errorf("failed to get managed clusters from DB, %v", err)
	}

	for _, clusterCliam := range clusterCliams {
		h.log.Infof("processing cluster claim %s", clusterCliam.Metadata.Name)
		readyStatus := h.getClusterClaimStatus(clusterCliam.Status.Conditions)
		h.log.Infof("cluster claim %s status: %s-%s-%s", clusterCliam.Metadata.Name,
			clusterCliam.Status.NodePoolStatus, clusterCliam.Status.ControlPlaneStatus, readyStatus)

		if !(strings.EqualFold(clusterCliam.Status.NodePoolStatus, "active") &&
			strings.EqualFold(clusterCliam.Status.ControlPlaneStatus, "active")) {
			h.log.Infof("cluster %s is not created", clusterCliam.Metadata.Name)
			return nil
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

		if strings.EqualFold(readyStatus, readyStatusValue) {
			managedCluster.ClusterDeployStatus = clusterReadyStatus
		} else {
			managedCluster.ClusterDeployStatus = clusterNotReadyStatus
		}

		err = managedcluster.StoreClusterAccessData(context.Background(), clusterCliam.Metadata.Namespace, managedCluster.Id)
		if err != nil {
			h.log.Errorf("failed to store cluster access data for %s, %v", managedCluster.Id, err)
			continue
		}

		err = h.dbStore.UpsertManagedCluster(managedCluster)
		if err != nil {
			h.log.Info("failed to update information to db, %v", err)
			continue
		}
		h.log.Infof("updated the cluster claim %s with status %s", managedCluster.ClusterName, managedCluster.ClusterDeployStatus)

		if managedCluster.ClusterDeployStatus == clusterReadyStatus {
			err = h.triggerClusterUpdates(clusterCliam.Spec.Id, managedCluster.Id)
			if err != nil {
				h.log.Info("failed to trigger cluster update workflow, %v", err)
				continue
			}
			h.log.Infof("triggered cluster update workflow for cluster %s", managedCluster.ClusterName)
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

func (h *ClusterClaimSyncHandler) triggerClusterUpdates(clusterName, managedClusterID string) error {
	proj, err := h.dbStore.GetCrossplaneProject()
	if err != nil {
		return err
	}

	ci := model.CrossplaneClusterUpdate{RepoURL: proj.GitProjectUrl, GitProjectId: proj.GitProjectId,
		ManagedClusterName: clusterName, ManagedClusterId: managedClusterID}
	wd := workers.NewConfig(h.tc, h.log)
	_, err = wd.SendEvent(context.TODO(), &model.ConfigureParameters{Resource: model.CrossPlaneResource, Action: model.CrossPlaneClusterUpdate}, ci)
	return err
}

func (h *ClusterClaimSyncHandler) getClusterClaimStatus(conditions []model.ClusterClaimCondition) (readyStatus string) {
	for _, condition := range conditions {
		switch strings.ToLower(condition.Type) {
		case readyStatusType:
			readyStatus = condition.Status
		}
	}
	return
}

func (h *ClusterClaimSyncHandler) deleteManagedClusters(clusterCliam model.ClusterClaim) error {

	clusters, err := h.getManagedClusters()
	if err != nil {
		return fmt.Errorf("failed to get managed clusters from DB, %v", err)
	}

	var clusterFound bool
	var managedCluster *captenpluginspb.ManagedCluster
	for _, v := range clusters {
		if v.ClusterName == clusterCliam.Metadata.Name {
			clusterFound = true
			managedCluster = v
			break
		}
	}

	if !clusterFound {
		h.log.Info("failed to delete managed cluster from DB, %s Cluster is not stored in ManagedClusters table", clusterCliam.Metadata.Name)
		return nil
	}

	managedCluster.ClusterDeployStatus = clusterDeletingStatus
	if err := h.dbStore.UpsertManagedCluster(managedCluster); err != nil {
		return fmt.Errorf("failed to update managed cluster from DB, %v", err)
	}

	err = h.triggerClusterDelete(clusterCliam.Spec.Id, managedCluster)
	if err != nil {
		return fmt.Errorf("failed to trigger cluster delete workflow, %v", err)
	}

	h.log.Infof("triggered cluster delete workflow for cluster %s", managedCluster.ClusterName)

	return nil
}

func (h *ClusterClaimSyncHandler) triggerClusterDelete(clusterName string, managedCluster *captenpluginspb.ManagedCluster) error {
	wd := workers.NewConfig(h.tc, h.log)

	proj, err := h.dbStore.GetCrossplaneProject()
	if err != nil {
		return err
	}
	ci := model.CrossplaneClusterUpdate{RepoURL: proj.GitProjectUrl, GitProjectId: proj.GitProjectId,
		ManagedClusterName: clusterName, ManagedClusterId: managedCluster.Id}

	wkfId, err := wd.SendAsyncEvent(context.TODO(), &model.ConfigureParameters{Resource: model.CrossPlaneResource, Action: model.CrossPlaneProjectDelete}, ci)
	if err != nil {
		managedCluster.ClusterDeployStatus = clusterFailedToDeleteStatus
		if err := h.dbStore.UpsertManagedCluster(managedCluster); err != nil {
			return fmt.Errorf("failed to update managed cluster from DB, %v", err)
		}
		return fmt.Errorf("failed to send event to workflow to configure %s, %v", managedCluster.ClusterEndpoint, err)
	}

	go h.monitorCrossplaneWorkflow(managedCluster, wkfId)

	h.log.Infof("Crossplane project delete %s config workflow %s created", managedCluster.ClusterEndpoint, wkfId)

	return nil
}

func (h *ClusterClaimSyncHandler) monitorCrossplaneWorkflow(managedCluster *captenpluginspb.ManagedCluster, wkfId string) {
	// during system reboot start monitoring, add it in map or somewhere.
	wd := workers.NewConfig(h.tc, h.log)
	_, err := wd.GetWorkflowInformation(context.TODO(), wkfId)
	if err != nil {
		managedCluster.ClusterDeployStatus = clusterFailedToDeleteStatus
		if err := h.dbStore.UpsertManagedCluster(managedCluster); err != nil {
			h.log.Errorf("failed to update managed cluster from DB, %v", err)
			return
		}
		h.log.Errorf("failed to send event to workflow to configure %s, %v", managedCluster.ClusterEndpoint, err)
		return
	}
	h.log.Infof("Successfuly removed the %s app config from clusters", managedCluster.ClusterName)

	if err := h.dbStore.DeleteManagedClusterById(managedCluster.Id); err != nil {
		h.log.Errorf("failed to delete managed cluster from DB, %v", err)
		return
	}
	h.log.Infof("Successfuly deleted managed cluster record for %s. cluster Id - %s", managedCluster.ClusterName, managedCluster.Id)

	if err = managedcluster.DeleteClusterAccessData(context.Background(), managedCluster.Id); err != nil {
		h.log.Errorf("failed to delete credential for %s, %v", managedCluster.Id, err)
		return
	}
	h.log.Infof("Successfuly deleted managed cluster credential for %s. cluster Id - %s", managedCluster.ClusterName, managedCluster.Id)

	h.log.Infof("Crossplane project delete %s config workflow %s completed", managedCluster.ClusterEndpoint, wkfId)
}

func (h *ClusterClaimSyncHandler) syncClusterClaimsWithDB(clusterClaims []model.ClusterClaim) error {

	clusters, err := h.getManagedClusters()
	if err != nil {
		return fmt.Errorf("failed to get managed clusters from DB, %v", err)
	}

	for _, cm := range clusterClaims {
		var isDeleteManagedCluster = true
		for _, c := range clusters {
			if c.ClusterName == cm.Metadata.Name {
				isDeleteManagedCluster = false
				break
			}
		}

		if isDeleteManagedCluster {
			if err = h.deleteManagedClusters(cm); err != nil {
				return fmt.Errorf("failed to delete ClusterCliam object, %v", err)
			}
		}
	}
	return nil
}
