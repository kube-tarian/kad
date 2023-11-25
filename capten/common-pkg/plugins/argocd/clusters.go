package argocd

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/io"
	"gopkg.in/yaml.v2"
	k8sapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	CredEntityName = "k8s"
	CredIdentifier = "kubeconfig"
)

func parseKubeconfigString(kubeconfigString string) (*k8sapi.Config, error) {
	config := &k8sapi.Config{}
	err := yaml.Unmarshal([]byte(kubeconfigString), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (a *ArgoCDClient) CreateOrUpdateCluster(ctx context.Context, clusterName, kubeconfigData string) error {
	kubeConfig, err := parseKubeconfigString(kubeconfigData)
	if err != nil {
		return err
	}

	clusterData, ok := kubeConfig.Clusters[clusterName]
	if !ok {
		return fmt.Errorf("cluster %s not found in kubeconfig", clusterName)
	}
	caData, err := base64.StdEncoding.DecodeString(string(clusterData.CertificateAuthorityData))
	if err != nil {
		return err
	}

	var clusterCAuthInfo *k8sapi.AuthInfo
	for _, authInfo := range kubeConfig.AuthInfos {
		clusterCAuthInfo = authInfo
		break
	}

	if clusterCAuthInfo == nil {
		return fmt.Errorf("auth info not found for cluster")
	}

	var clientCertData []byte
	if len(clusterCAuthInfo.ClientCertificateData) != 0 {
		clientCertData, err = base64.StdEncoding.DecodeString(string(clusterCAuthInfo.ClientCertificateData))
		if err != nil {
			return err
		}
	}

	a.logger.Infof("Cluster create or update request for cluster %s", clusterName)
	var clientKeyData []byte
	if len(clusterCAuthInfo.ClientKeyData) != 0 {
		clientKeyData, err = base64.StdEncoding.DecodeString(string(clusterCAuthInfo.ClientKeyData))
		if err != nil {
			return err
		}
	}

	conn, appClient, err := a.client.NewClusterClient()
	if err != nil {
		return fmt.Errorf("failed to create argocd cluster client, %v", err)
	}
	defer io.Close(conn)

	var update bool
	_, err = appClient.Create(ctx, &cluster.ClusterCreateRequest{
		Cluster: &v1alpha1.Cluster{
			Server: clusterData.Server,
			Name:   clusterName,
			Config: v1alpha1.ClusterConfig{
				BearerToken: clusterCAuthInfo.Token,
				Username:    clusterCAuthInfo.Username,
				Password:    clusterCAuthInfo.Password,
				TLSClientConfig: v1alpha1.TLSClientConfig{
					ServerName: clusterData.Server,
					CAData:     caData,
					CertData:   clientCertData,
					KeyData:    clientKeyData,
				},
			},
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			update = true
		}
		return fmt.Errorf("failed to create cluster %s, %v", clusterName, err)
	}

	if update {
		_, err := appClient.Update(ctx, &cluster.ClusterUpdateRequest{
			Cluster: &v1alpha1.Cluster{
				Server: clusterData.Server,
				Name:   clusterName,
				Config: v1alpha1.ClusterConfig{
					BearerToken: clusterCAuthInfo.Token,
					Username:    clusterCAuthInfo.Username,
					Password:    clusterCAuthInfo.Password,
					TLSClientConfig: v1alpha1.TLSClientConfig{
						ServerName: clusterData.Server,
						CAData:     caData,
						CertData:   clientCertData,
						KeyData:    clientKeyData,
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to update cluster %s, %v", clusterName, err)
		}
		a.logger.Infof("Cluster %s created", clusterName)
	} else {
		a.logger.Infof("Cluster %s updated", clusterName)
	}
	return nil
}

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
