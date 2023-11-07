package sync

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"

	pb "github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"

	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
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

type Fetch struct {
	log         logging.Logger
	client      *k8s.K8SClient
	db          *captenstore.Store
	creds       credentials.CredentialAdmin
	avlClusters map[string]*pb.ManagedCluster
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

	credAdmin, err := credentials.NewCredentialAdmin(context.TODO())
	if err != nil {
		log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client")
		return nil, err
	}

	avlClusters, err := getManagedClusterEndpointMap(db)
	if err != nil {
		return nil, fmt.Errorf("failed to execute  getManagedClusterEndpointMap, err: %v", err)
	}

	return &Fetch{log: log, client: k8sclient, db: db, creds: credAdmin, avlClusters: avlClusters}, nil
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

	fetch.UpdateClusterDetails(clObj.Items)

	fetch.log.Info("succesfully sync-ed cluster-claims resources")
}

func (fetch *Fetch) UpdateClusterDetails(clObj []model.ClusterClaim) {
	for _, obj := range clObj {
		for _, status := range obj.Status.Conditions {
			if status.Type != statusType {
				continue
			}

			// if status.Status != statusValue {
			// 	fetch.log.Info("%s in namespace %s, status is %s, so skiping update to db.",
			// 		obj.Spec.Id, obj.Metadata.Namespace, status.Status)
			// 	continue
			// }

			// get the cluster endpoint and kubeconfig file from the secrets
			req := &k8s.SecretDetailsRequest{Namespace: "crossplane-system",
				SecretName: fmt.Sprintf(clusterSecretName, obj.Spec.Id)}
			resp, err := fetch.client.FetchSecretDetails(req)
			if err != nil {
				fetch.log.Info("%s in namespace %s, failed to get secret: %v",
					req.SecretName, req.Namespace, err)
				continue
			}

			clusterEndpoint, err := getBase64DecodedString(resp.Data[k8sEndpoint])
			if err != nil {
				fetch.log.Info("failed to decode base64 value: %v", err)
				continue
			}

			managedCluster := &pb.ManagedCluster{}
			clusterObj, ok := fetch.avlClusters[clusterEndpoint]
			if !ok {
				managedCluster.Id = uuid.New().String()
			} else {
				fetch.log.Info("found existing Id: %s, updating the latest information ", clusterObj.Id)
				managedCluster.Id = clusterObj.Id
			}

			managedCluster.ClusterName = fmt.Sprintf(clusterName, obj.Spec.Id, obj.Metadata.Namespace)
			managedCluster.ClusterDeployStatus = status.Status
			managedCluster.ClusterEndpoint = clusterEndpoint

			clusterDetails := map[string]string{}
			clusterDetails[kubeConfig] = resp.Data[kubeConfig]
			clusterDetails[k8sClusterCA] = resp.Data[k8sClusterCA]

			err = fetch.creds.PutCredential(context.TODO(), credentials.GenericCredentialType, managedClusterEntityName, managedCluster.Id, clusterDetails)

			if err != nil {
				fetch.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", managedCluster.Id)
				fetch.log.Errorf("failed to store credential for %s, %v", managedCluster.Id, err)
				continue
			}

			managedCluster.LastUpdateTime = time.Now().Format(time.RFC3339)
			err = fetch.db.UpsertManagedCluster(managedCluster)
			if err != nil {
				fetch.log.Info("failed to update information to db: %v", err)
				continue
			}

			fetch.avlClusters[clusterEndpoint] = managedCluster
		}

	}
}

func getManagedClusterEndpointMap(db *captenstore.Store) (map[string]*pb.ManagedCluster, error) {
	clusters, err := db.GetManagedClusters()
	if err != nil {
		return nil, fmt.Errorf("failed to get the managed cluster information from db: %v", err)
	}

	clusterEndpointMap := map[string]*pb.ManagedCluster{}
	for _, cluster := range clusters {
		clusterEndpointMap[cluster.ClusterEndpoint] = cluster
	}

	return clusterEndpointMap, nil
}

func getBase64DecodedString(encodedString string) (string, error) {
	decodedByte, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "", err
	}
	return string(decodedByte), nil
}
