package argocd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/io"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"gopkg.in/yaml.v2"
)

const (
	CredEntityName = "k8s"
	CredIdentifier = "kubeconfig"
)

func (a *ArgoCDClient) CreateCluster(ctx context.Context, serverName string, config map[string]string) (*v1alpha1.Cluster, error) {
	cred, err := credential.GetGenericCredential(ctx, CredEntityName, CredIdentifier)
	if err != nil {
		fmt.Println("Error occured while fetching kubeconfig")
		fmt.Printf(err.Error())
	}

	byteConfig, err := json.Marshal(cred)
	if err != nil {
		return nil, err
	}

	kubeconfig := KubeConfig{}
	if err := yaml.Unmarshal(byteConfig, &kubeconfig); err != nil {
		return nil, err
	}

	fmt.Println("byteConfig =>", string(byteConfig))
	fmt.Println("Config =>", kubeconfig)

	fmt.Println("clusterReq.Server" + kubeconfig.Clusters[0].Cluster.Server)

	conn, appClient, err := a.client.NewClusterClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	caData, err := base64.StdEncoding.DecodeString(strings.TrimSpace(kubeconfig.Clusters[0].Cluster.CertificateAuthorityData))
	if err != nil {
		return nil, err
	}

	resp, err := appClient.Create(ctx, &cluster.ClusterCreateRequest{
		Cluster: &v1alpha1.Cluster{
			Server: kubeconfig.Clusters[0].Cluster.Server,
			Name:   kubeconfig.Contexts[0].Context.Cluster,
			Config: v1alpha1.ClusterConfig{
				BearerToken: kubeconfig.Users[0].User.ClientCertificateData,
				TLSClientConfig: v1alpha1.TLSClientConfig{
					ServerName: serverName,
					CAData:     caData,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// func parseKubeConfig(kubeconfigMap map[string]string) error {
// 	clusters, ok := kubeconfigMap["clusters"].([]interface{})
// 	if !ok {
// 		fmt.Println("Error: clusters not found in kubeconfig")
// 		return
// 	}
// }

func (a *ArgoCDClient) DeleteCluster(ctx context.Context, clusterURL string) (*cluster.ClusterResponse, error) {
	conn, appClient, err := a.client.NewClusterClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	resp, err := appClient.Delete(ctx, &cluster.ClusterQuery{
		Id: &cluster.ClusterID{
			Value: clusterURL,
		},
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *ArgoCDClient) GetCluster(ctx context.Context, clusterURL string) (*v1alpha1.Cluster, error) {
	conn, appClient, err := a.client.NewClusterClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	cluster, err := appClient.Get(ctx, &cluster.ClusterQuery{
		Id: &cluster.ClusterID{
			Value: clusterURL,
		},
	})
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (a *ArgoCDClient) ListClusters(ctx context.Context) (*v1alpha1.ClusterList, error) {
	conn, appClient, err := a.client.NewClusterClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	list, err := appClient.List(ctx, &cluster.ClusterQuery{})
	if err != nil {
		return nil, err
	}

	return list, nil
}
