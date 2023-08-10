package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	fmt.Println("Registration testing")
	gr, err := grpc.Dial("captenserver.test.app:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("GRPC failed", err)
		return
	}
	sc := serverpb.NewServerClient(gr)
	ctx := context.TODO()
	ctx = metadata.AppendToOutgoingContext(ctx, "organizationID", "eab9af58-1cbd-4f25-b3f3-d07ee99a5baa")

	resp1, err := sc.GetClusters(ctx, &serverpb.GetClustersRequest{})
	if err != nil {
		fmt.Println("error with get registrations ", err)
	} else {
		fmt.Println("registrations fetch done", resp1.Status, resp1.StatusMessage, len(resp1.Data))
		for _, cluster := range resp1.Data {
			fmt.Printf("cluster: %+v\n", *cluster)
		}
	}

	cadata, err := os.ReadFile("/var/capten/cert/ca.crt")
	if err != nil {
		fmt.Println("ca failed", err)
		return
	}

	cdata, err := os.ReadFile("/var/capten/cert/client.crt")
	if err != nil {
		fmt.Println("client failed", err)
		return
	}

	ckey, err := os.ReadFile("/var/dev/capten/cert/client.key")
	if err != nil {
		fmt.Println("key failed", err)
		return
	}

	nresp, err := sc.NewClusterRegistration(ctx, &serverpb.NewClusterRegistrationRequest{
		AgentEndpoint:     "https://captenagent.dev.test.app",
		ClusterName:       "awscluster-3",
		ClientKeyData:     string(ckey),
		ClientCertData:    string(cdata),
		ClientCAChainData: string(cadata),
	})
	if err != nil {
		fmt.Println("error with registration ", err)
	} else {
		fmt.Println("Server testing done", nresp.Status, nresp.StatusMessage, nresp.ClusterID)
	}

	nresp1, err := sc.NewClusterRegistration(ctx, &serverpb.NewClusterRegistrationRequest{
		AgentEndpoint:     "https://captenagent.dev.test.app",
		ClusterName:       "awscluster-4",
		ClientKeyData:     "fsdfgsdfsdf",
		ClientCertData:    "bcvbvcbhhhfh",
		ClientCAChainData: "tyrytertyeyeyey",
	})
	if err != nil {
		fmt.Println("error with registration ", err)
	} else {
		fmt.Println("Server testing done", nresp1.Status, nresp1.StatusMessage, nresp1.ClusterID)
	}

	gresp2, err := sc.GetClusters(ctx, &serverpb.GetClustersRequest{})
	if err != nil {
		fmt.Println("error with get registrations ", err)
	} else {
		fmt.Println("registrations fetch done", gresp2.Status, gresp2.StatusMessage, len(gresp2.Data))
		for _, cluster := range gresp2.Data {
			fmt.Printf("cluster: %+v\n", *cluster)
		}
	}

	uresp, err := sc.UpdateClusterRegistration(ctx, &serverpb.UpdateClusterRegistrationRequest{
		ClusterID:         nresp.ClusterID,
		AgentEndpoint:     "https://captenagent.dev.test.app1",
		ClusterName:       "awscluster-3-update",
		ClientKeyData:     string(ckey),
		ClientCertData:    string(cdata),
		ClientCAChainData: string(cadata),
	})
	if err != nil {
		fmt.Println("error with get registrations ", err)
	} else {
		fmt.Println("registrations fetch done", uresp.Status, uresp.StatusMessage)
	}

	gresp3, err := sc.GetClusters(ctx, &serverpb.GetClustersRequest{})
	if err != nil {
		fmt.Println("error with get registrations ", err)
	} else {
		fmt.Println("registrations fetch done", gresp3.Status, gresp3.StatusMessage, len(gresp3.Data))
		for _, cluster := range gresp3.Data {
			fmt.Printf("cluster: %+v\n", *cluster)
		}
	}

	for _, cluster := range gresp3.Data {
		dresp, err := sc.DeleteClusterRegistration(ctx, &serverpb.DeleteClusterRegistrationRequest{ClusterID: cluster.ClusterID})
		if err != nil {
			fmt.Println("error with get registrations ", err, cluster.ClusterID)
		} else {
			fmt.Println("registrations delete done", dresp.Status, dresp.StatusMessage, cluster.ClusterID)
		}
	}
}
