package captensdk

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kube-tarian/kad/integrator/capten-sdk/agentpb"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
)

var log = logging.NewLogger()

// func TestApplicationCreate(t *testing.T) {
// 	appClient, err := setup(t)
// 	if err != nil {
// 		return
// 	}

// 	req := &DeploymentRequestPayload{
// 		PluginName: "argocd",
// 		Action:     "deployment",
// 		Data: DeploymentRequestData{
// 			RepoName:    "argocd-example",
// 			RepoURL:     "https://gitlab.privatecloud.sk/vladoportos/argo-cd-example.git",
// 			ChartName:   "hello-world",
// 			Namespace:   "default",
// 			ReleaseName: "hello-world",
// 			Timeout:     5,
// 		},
// 	}
// 	_, err = appClient.Create(req)
// 	assert.Nilf(t, err, "application create should be success")
// }

// func TestApplicationDelete(t *testing.T) {
// 	appClient, err := setup(t)
// 	if err != nil {
// 		return
// 	}

// 	req := &DeploymentRequestPayload{
// 		PluginName: "argocd",
// 		Action:     "deployment",
// 		Data: DeploymentRequestData{
// 			RepoName:    "argocd-example",
// 			RepoURL:     "https://gitlab.privatecloud.sk/vladoportos/argo-cd-example.git",
// 			ChartName:   "hello-world",
// 			Namespace:   "default",
// 			ReleaseName: "hello-world",
// 			Timeout:     5,
// 		},
// 	}
// 	_, err = appClient.Delete(req)
// 	assert.Nilf(t, err, "application create should be success")
// }

// func setup(t *testing.T) (*ApplicationClient, error) {
// 	client, err := NewClient(log)
// 	assert.Nilf(t, err, "New client should be initialized")
// 	if err != nil {
// 		return nil, err
// 	}

// 	appClient, err := client.NewApplicationClient(&TransportSSLOptions{IsSSLEnabled: false})
// 	assert.Nilf(t, err, "New application client should be initialized")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return appClient, err
// }

func TestApplicationDelete(t *testing.T) {
	os.Setenv("DEPLOYMENT_REQUEST_FILE_PATH", "tests/hello-world-argocd-plugin.json")
	main()
}

type JobRequest struct {
	Opertion string                    `json:"operation" required:"true"`
	Payload  *DeploymentRequestPayload `json:"payload" required:"true"`
}

func DeploymentRequestFromFile(fileName string) (*DeploymentRequestPayload, error) {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		log.Errorf("File %s read failed, %v", fileName, err)
		return nil, err
	}

	req := &JobRequest{}
	err = json.Unmarshal(fileData, req)
	if err != nil {
		log.Errorf("Failed to unmarshal job request, %v", err)
		return nil, err
	}

	return req.Payload, nil
}

func main() {
	deployRequestFilePath := os.Getenv("DEPLOYMENT_REQUEST_FILE_PATH")
	if deployRequestFilePath == "" {
		log.Fatalf("DEPLOYMENT_REQUEST_FILE_PATH not set")
	}

	req, err := DeploymentRequestFromFile(deployRequestFilePath)
	if err != nil {
		log.Fatalf("Preparing Deployment request from file path %s failed, %v", deployRequestFilePath, err)
	}

	client, err := NewClient(log)
	if err != nil {
		log.Fatalf("New Capten SDK client initialization failed, %v", err)
	}

	appClient, err := client.NewApplicationClient(&TransportSSLOptions{IsSSLEnabled: false})
	if err != nil {
		return
	}

	var response *agentpb.JobResponse
	switch req.Action {
	case "install", "update":
		response, err = appClient.Create(req)
	case "delete":
		response, err = appClient.Delete(req)
	}
	if err != nil {
		log.Fatalf("%s application failed, %v", req.Action, err)
	}
	log.Infof("application %s from plugin %s is success, response: %+v", req.PluginName, req.Action, response)
}
