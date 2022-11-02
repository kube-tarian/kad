package k8s

import (
	"fmt"
	"os"
)

type info struct {
	configPath     string
	tokenPath      string
	kubeApiServer  string
	kubeCaFilePath string
}

func GetInfo() *info {
	//TODO: use the env variable
	return &info{
		configPath:     "/Users/srikrishnabh/workspace/intelops/k8s/temporalk8scfg",
		tokenPath:      "/var/run/secrets/kubernetes.io/serviceaccount/token",
		kubeApiServer:  "https://kubernetes.default.svc.cluster.local",
		kubeCaFilePath: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
	}
}

func (i *info) GetConfigPath() string {
	return ""
}

func (i *info) GetToken() string {
	token, err := os.ReadFile(i.tokenPath)
	if err != nil {
		fmt.Println("failed to read the tokenfile")
	}

	return string(token)
}

func (i *info) GetK8sEndpoint() string {
	return i.kubeApiServer
}

func (i *info) GetK8sCaFilePath() string {
	return i.kubeCaFilePath
}
