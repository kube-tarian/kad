package argocd

import (
	"fmt"
	appsv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"testing"

	cmdutil "github.com/argoproj/argo-cd/v2/cmd/util"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/model"
)

func Test_ClusterAdd(t *testing.T) {
	client, err := newTestClient(logging.NewLogger())
	if err != nil {
		fmt.Printf("\n failed to get test client %v", err)
	}
	data := map[string]interface{}{}
	data["clusterOpts"] = &cmdutil.ClusterOptions{
		InCluster:               false,
		Upsert:                  false,
		ServiceAccount:          "",
		AwsRoleArn:              "",
		AwsClusterName:          "",
		SystemNamespace:         "",
		Namespaces:              nil,
		ClusterResources:        false,
		Name:                    "dev-cluster1",
		Project:                 "dev1",
		Shard:                   0,
		ExecProviderCommand:     "",
		ExecProviderArgs:        nil,
		ExecProviderEnv:         nil,
		ExecProviderAPIVersion:  "",
		ExecProviderInstallHint: "",
	}
	//data["pathOpts"] = &clientcmd.PathOptions{
	//	GlobalFile:        "",
	//	EnvVar:            "",
	//	ExplicitFileFlag:  "",
	//	GlobalFileSubpath: "",
	//	LoadingRules:      nil,
	//}
	data["contextName"] = "admin@talosconfig-userdata"
	req := &model.ConfigPayload{
		Action:     "add",
		Data:       data,
		PluginName: "argocd",
		Resource:   "cluster",
	}
	_, err = client.ClusterAdd(req)
	if err != nil {
		t.Errorf("\n failed to add cluster: %v", err)
		return
	}
	fmt.Println("success")
}

func Test_ClusterDelete(t *testing.T) {
	client, err := newTestClient(logging.NewLogger())
	if err != nil {
		fmt.Printf("\n failed to get test client %v", err)
	}
	data := map[string]interface{}{}
	data["clusterSelector"] = "http://my-server"
	req := &model.ConfigPayload{
		Action:     "delete",
		Data:       data,
		PluginName: "argocd",
		Resource:   "cluster",
	}
	_, err = client.ClusterDelete(req)
	if err != nil {
		t.Errorf("\n failed to add cluster: %v", err)
		return
	}
	fmt.Println("success")
}

func Test_RepoAdd(t *testing.T) {
	client, err := newTestClient(logging.NewLogger())
	if err != nil {
		fmt.Printf("\n failed to get test client %v", err)
	}
	data := map[string]interface{}{}
	inputRepo := appsv1.Repository{
		Name:     "TestRepo",
		Repo:     "git@github.com:argoproj/argo-cd.git",
		Username: "someUsername",
		Password: "somePassword",
	}
	data["RepoOptions"] = &cmdutil.RepoOptions{
		Repo:                           inputRepo,
		Upsert:                         false,
		SshPrivateKeyPath:              "",
		InsecureIgnoreHostKey:          false,
		InsecureSkipServerVerification: false,
		TlsClientCertPath:              "",
		TlsClientCertKeyPath:           "",
		EnableLfs:                      false,
		EnableOci:                      false,
		GithubAppId:                    0,
		GithubAppInstallationId:        0,
		GithubAppPrivateKeyPath:        "",
		GitHubAppEnterpriseBaseURL:     "",
		Proxy:                          "",
	}
	//data["pathOpts"] = &clientcmd.PathOptions{
	//	GlobalFile:        "",
	//	EnvVar:            "",
	//	ExplicitFileFlag:  "",
	//	GlobalFileSubpath: "",
	//	LoadingRules:      nil,
	//}
	data["contextName"] = "admin@talosconfig-userdata"
	req := &model.ConfigPayload{
		Action:     "add",
		Data:       data,
		PluginName: "argocd",
		Resource:   "repo",
	}
	_, err = client.RepoAdd(req)
	if err != nil {
		fmt.Printf("\n failed to add cluster: %v", err)
		return
	}
	fmt.Println("success")
}

func Test_RepoDelete(t *testing.T) {
	client, err := newTestClient(logging.NewLogger())
	if err != nil {
		fmt.Printf("\n failed to get test client %v", err)
	}
	data := map[string]interface{}{}
	inputRepo := appsv1.Repository{
		Name:     "TestRepo",
		Repo:     "git@github.com:argoproj/argo-cd.git",
		Username: "someUsername",
		Password: "somePassword",
	}
	data["RepoOptions"] = &cmdutil.RepoOptions{
		Repo:                           inputRepo,
		Upsert:                         false,
		SshPrivateKeyPath:              "",
		InsecureIgnoreHostKey:          false,
		InsecureSkipServerVerification: false,
		TlsClientCertPath:              "",
		TlsClientCertKeyPath:           "",
		EnableLfs:                      false,
		EnableOci:                      false,
		GithubAppId:                    0,
		GithubAppInstallationId:        0,
		GithubAppPrivateKeyPath:        "",
		GitHubAppEnterpriseBaseURL:     "",
		Proxy:                          "",
	}
	data["contextName"] = "admin@talosconfig-userdata"
	req := &model.ConfigPayload{
		Action:     "delete",
		Data:       data,
		PluginName: "argocd",
		Resource:   "repo",
	}
	_, err = client.RepoDelete(req)
	if err != nil {
		fmt.Printf("\n failed to add cluster: %v", err)
		return
	}
	fmt.Println("success")
}

func newTestClient(logger logging.Logger) (*ArgoCDCLient, error) {

	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: "localhost:8080",
		Insecure:   true,
		AuthToken:  "testtoken",
	})
	if err != nil {
		fmt.Printf("\n failed to get apiclient %v", err)
	}
	return &ArgoCDCLient{
		conf:   &Configuration{},
		logger: logger,
		client: client,
	}, nil
}
