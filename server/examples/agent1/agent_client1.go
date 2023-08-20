package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
)

func main() {
	log := logging.NewLogger()
	cfg, err := getAgentConfig("https://captenagent.dev.optimizor.app")
	if err != nil {
		log.Fatalf("failed to load agent config: ", err)
	}

	ac, err := agent.NewAgent(log, cfg, nil)
	if err != nil {
		log.Fatalf("failed to connect to agent: ", err)
		return
	}
	//storeServiceCred(ac.GetClient())
	//storeGenericCred(ac.GetClient())
	getLaunchUIApps(ac.GetClient())
}

func storeServiceCred(ac agentpb.AgentClient) {
	serviceCred := credentials.ServiceCredential{
		UserName: "testuser",
		Password: "password2",
	}
	serviceCredMap := credentials.PrepareServiceCredentialMap(serviceCred)
	_, err := ac.StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: credentials.ServiceUserCredentialType,
		CredEntityName: "testentity",
		CredIdentifier: "testentityuser",
		Credential:     serviceCredMap,
	})
	if err != nil {
		fmt.Println("store error: ", err)
		return
	}
	fmt.Println("successful")
}

func storeGenericCred(ac agentpb.AgentClient) {
	cred := map[string]string{
		"UserName": "gentestuser",
		"Password": "password3",
	}

	_, err := ac.StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: credentials.GenericCredentialType,
		CredEntityName: "gentestentity",
		CredIdentifier: "gentestentityuser",
		Credential:     cred,
	})
	if err != nil {
		fmt.Println("store error: ", err)
		return
	}
	fmt.Println("successful")
}

func getLaunchUIApps(ac agentpb.AgentClient) {
	resp, err := ac.GetClusterAppLaunches(context.Background(), &agentpb.GetClusterAppLaunchesRequest{})
	if err != nil {
		fmt.Println("get app launch error: ", err)
		return
	}
	if resp.Status != agentpb.StatusCode_OK {
		fmt.Println("get app launch error: ", resp.StatusMessage)
		return
	}

	for _, app := range resp.LaunchConfigList {
		fmt.Printf("app launch %s : %s, '%s', %s\n", app.ReleaseName, app.Category, app.Description, app.LaunchURL)
		fmt.Printf("app launch icon %s : %s\n", app.ReleaseName, string(app.Icon))

		dat, _ := json.Marshal(app)
		fmt.Printf("%s", string(dat))
	}

	fmt.Println("successful")
}

func getAgentConfig(address string) (*agent.Config, error) {
	cadata, err := os.ReadFile("/home/venkatk/dev/capten/cert/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("ca failed, %v", err)
	}
	cdata, err := os.ReadFile("/home/venkatk/dev/capten/cert/client.crt")
	if err != nil {
		return nil, fmt.Errorf("client failed, %v", err)
	}
	ckey, err := os.ReadFile("/home/venkatk/dev/capten/cert/client.key")
	if err != nil {
		return nil, fmt.Errorf("key failed, %v", err)
	}
	return &agent.Config{Address: address, CaCert: string(cadata), Cert: string(cdata), Key: string(ckey)}, nil
}
