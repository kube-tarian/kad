package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	fmt.Println("Server testing")
	gr, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("GRPC failed", err)
		return
	}
	sc := serverpb.NewServerClient(gr)
	ctx := context.TODO()
	//ctx = metadata.AppendToOutgoingContext(ctx, "organizationID", "996162a1-1df7-44b7-8347-1cb1acc70666")
	ctx = metadata.AppendToOutgoingContext(ctx, "organizationID", "996162a1-1df7-44b7-8347-1cb1acc70688")
	//res, err := sc.GetClusters(ctx, &serverpb.GetClustersRequest{})
	//os.ReadFile("ca.crt")
	//ca, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/new-certs/ca.crt")
	ca, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/capten/cert/ca.crt")
	if err != nil {
		fmt.Println("ca failed", err)
		return
	}

	//clinet, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/new-certs/client.crt")
	clinet, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/capten/cert/client.crt")
	if err != nil {
		fmt.Println("client failed", err)
		return
	}
	//ckeynew, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/new-certs/client.key")
	ckeynew, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/capten/cert/client.key")
	if err != nil {
		fmt.Println("key failed", err)
		return
	}

	// fmt.Println("ClientKeyData", string(clinet))
	// fmt.Println("ClientCertData", string(ckeynew))
	// fmt.Println("ClientCAChainData", string(ca))
	base64.StdEncoding.EncodeToString(ca)
	resp, err := sc.NewClusterRegistration(ctx, &serverpb.NewClusterRegistrationRequest{
		AgentEndpoint:     "https://captenagent.dev.optimizor.app",
		ClusterName:       "NewCluster",
		ClientKeyData:     base64.StdEncoding.EncodeToString(ckeynew),
		ClientCertData:    base64.StdEncoding.EncodeToString(clinet),
		ClientCAChainData: base64.StdEncoding.EncodeToString(ca),
	})
	//fmt.Println(res, err)
	fmt.Println("Server testing done", resp, err)
	/*
		uresp, err := sc.UpdateClusterRegistration(ctx, &serverpb.UpdateClusterRegistrationRequest{
			ClusterID:         resp.ClusterID,
			AgentEndpoint:     "https://captenagent.dev.optimizor.app",
			ClusterName:       "testCluster",
			ClientKeyData:     string(ckeynew),
			ClientCertData:    string(clinet),
			ClientCAChainData: string(ca),
		})
		//fmt.Println(res, err)
		fmt.Println("Server update testing done", uresp, err)*/
}
