package argocd

import (
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"net/http/httptest"
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	cmdutil "github.com/argoproj/argo-cd/v2/cmd/util"
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
	data["pathOpts"] = &clientcmd.PathOptions{
		GlobalFile:        "",
		EnvVar:            "",
		ExplicitFileFlag:  "",
		GlobalFileSubpath: "",
		LoadingRules:      nil,
	}
	data["contextname"] = "dev"
	req := &model.ConfigPayload{
		Action:     "Add",
		Data:       data,
		PluginName: "argocd",
		Resource:   "",
	}
	_, err = client.ClusterAdd(req)
	if err != nil {
		fmt.Printf("\n failed to add cluster: %v", err)
	}
	fmt.Println("success")
}

func newTestClient(logger logging.Logger) (*ArgoCDCLient, error) {
	//cfg, err := fetchConfiguration(logger)
	//if err != nil {
	//	return nil, err
	//}
	var mock MockServer
	s := httptest.NewServer(mock.Handler())
	defer s.Close()
	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: s.URL,
		Insecure:   false,
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
