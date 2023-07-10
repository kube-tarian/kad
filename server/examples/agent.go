package main

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("captenagent.dev.optimizor.app:80",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("connect error: ", err)
		return
	}
	storeServiceCred(conn)
}

func storeServiceCred(conn grpc.ClientConnInterface) {
	serviceCred := credentials.ServiceCredential{
		UserName: "user",
		Password: "password2",
	}
	serviceCredMap := credentials.PrepareServiceCredentialMap(serviceCred)
	agentClient := agentpb.NewAgentClient(conn)
	_, err := agentClient.StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: credentials.ServiceUserCredentialType,
		CredEntityName: "vitess",
		CredIdentifier: "user",
		Credential:     serviceCredMap,
	})
	if err != nil {
		fmt.Println("store error: ", err)
		return
	}
	fmt.Println("successful: ", err)
}
