package main

import (
	"context"
	"fmt"
	"os"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/server/pkg/agent"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
)

func main() {
	log := logging.NewLogger()
	cfg, err := getAgentConfig("https://captenagent.dev.test.app")
	if err != nil {
		log.Fatalf("failed to load agent config: ", err)
	}
	oryclient, err := oryclient.NewOryClient(log)
	if err != nil {
		log.Fatal("OryClient initialization failed", err)
	}
	ac, err := agent.NewAgent(log, cfg, oryclient)
	if err != nil {
		log.Fatalf("failed to connect to agent: ", err)
		return
	}
	storeServiceCred(ac.GetClient())
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

func getAgentConfig(address string) (*agent.Config, error) {
	cadata, err := os.ReadFile("/var/capten/cert/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("ca failed, %v", err)
	}
	cdata, err := os.ReadFile("/var/capten/cert/client.crt")
	if err != nil {
		return nil, fmt.Errorf("client failed, %v", err)
	}
	ckey, err := os.ReadFile("/var/capten/cert/client.key")
	if err != nil {
		return nil, fmt.Errorf("key failed, %v", err)
	}
	return &agent.Config{Address: address, CaCert: string(cadata), Cert: string(cdata), Key: string(ckey)}, nil
}
