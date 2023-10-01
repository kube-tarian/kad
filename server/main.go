package main

import (
	"fmt"

	"crypto/tls"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// func main() {
// 	fmt.Println("Server testing")
// 	gr, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		fmt.Println("GRPC failed", err)
// 		return
// 	}
// 	sc := serverpb.NewServerClient(gr)
// 	ctx := context.TODO()
// 	//ctx = metadata.AppendToOutgoingContext(ctx, "organizationID", "996162a1-1df7-44b7-8347-1cb1acc70666")
// 	ctx = metadata.AppendToOutgoingContext(ctx, "organizationID", "996162a1-1df7-44b9-8347-1cb1acc78888")
// 	//res, err := sc.GetClusters(ctx, &serverpb.GetClustersRequest{})
// 	//os.ReadFile("ca.crt")
// 	//ca, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/new-certs/ca.crt")
// 	ca, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/capten-client-demo-auth-certs/ca.crt")
// 	if err != nil {
// 		fmt.Println("ca failed", err)
// 		return
// 	}

// 	//clinet, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/new-certs/client.crt")
// 	clinet, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/capten-client-demo-auth-certs/client.crt")
// 	if err != nil {
// 		fmt.Println("client failed", err)
// 		return
// 	}
// 	//ckeynew, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/new-certs/client.key")
// 	ckeynew, err := os.ReadFile("/mnt/c/Users/hkin/Desktop/free/capten-client-demo-auth-certs/client.key")
// 	if err != nil {
// 		fmt.Println("key failed", err)
// 		return
// 	}

// 	fmt.Println("ClientKeyData", string(clinet))
// 	fmt.Println("ClientCertData", string(ckeynew))
// 	fmt.Println("ClientCAChainData", string(ca))
// 	base64.StdEncoding.EncodeToString(ca)
// 	resp, err := sc.NewClusterRegistration(ctx, &serverpb.NewClusterRegistrationRequest{
// 		AgentEndpoint:     "https://captenagent.demo.optimizor.app",
// 		ClusterName:       "DemoCluster-27",
// 		ClientKeyData:     base64.StdEncoding.EncodeToString(ckeynew),
// 		ClientCertData:    base64.StdEncoding.EncodeToString(clinet),
// 		ClientCAChainData: base64.StdEncoding.EncodeToString(ca),
// 	})
// 	fmt.Println(resp, err)
// 	//fmt.Println("Server testing done", resp, err)

// 	// uresp, err := sc.UpdateClusterRegistration(ctx, &serverpb.UpdateClusterRegistrationRequest{
// 	// 	ClusterID:         "6d25b32f-4da3-11ee-8943-361c8f970b48",
// 	// 	AgentEndpoint:     "https://captenagent.demo.optimizor.app",
// 	// 	ClusterName:       "testCluster",
// 	// 	ClientKeyData:     base64.StdEncoding.EncodeToString(ckeynew),
// 	// 	ClientCertData:    base64.StdEncoding.EncodeToString(clinet),
// 	// 	ClientCAChainData: base64.StdEncoding.EncodeToString(ca),
// 	// })
// 	// //fmt.Println(res, err)
// 	//fmt.Println("Server update testing done", uresp, err)
// }

func NewClient() (*client.StargateClient, error) {
	// serviceCredential, err := credential.GetGenericCredential(context.Background(),
	// 	conf.EntityName, conf.CredentailIdentifier)
	// if err != nil {
	// 	return nil, err
	// }

	serviceCredential := map[string]string{"TOKEN": "AstraCS:MNexxhzQYbrZtExtcewSWflp:823ee0fcd257ac4b738f7c66eaa383dccb8d8dd5aea7bbaf40e40c3de62b7b5d",
		"ASTRA_DB_ID": "fd4fc682-686a-479d-be68-881a6cd3c2df", "ASTRA_DB_REGION": "us-west-2"}

	token := serviceCredential["TOKEN"]
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial("fd4fc682-686a-479d-be68-881a6cd3c2df-us-west-2.apps.astra.datastax.com:443", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(token),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to astra db, %w", err)
	}

	session, err := client.NewStargateClientWithConn(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating stargate client, %w", err)
	}

	return session, nil
}

func UpdateCacheAppLaunches(a *client.StargateClient, orgID, clusterID string, appLaunches []*agentpb.AppLaunchConfig) error {
	appResponse, err := GetCacheAppLaunches(a, orgID, clusterID)
	if err != nil {
		return fmt.Errorf("failed to update the applaunches cache, err %w", err)
	}

	if len(appLaunches) == 0 && len(appResponse.LaunchConfigList) == 0 {
		return nil
	}

	// First get the data and then insert
	insertBatchQuery := make([]*pb.BatchQuery, 0)
	updateBatchQuery := make([]*pb.BatchQuery, 0)
	deleteBatchQuery := make([]*pb.BatchQuery, 0)
	dbAvailableReleaseName := make(map[string]*agentpb.AppLaunchConfig)
	for _, dbAppLaunches := range appResponse.LaunchConfigList {
		dbAvailableReleaseName[dbAppLaunches.ReleaseName] = dbAppLaunches
	}

	sentAvailableReleaseName := make(map[string]bool)
	for _, sentApp := range appLaunches {
		sentAvailableReleaseName[sentApp.ReleaseName] = true
	}

	for _, dbApp := range dbAvailableReleaseName {
		if _, found := sentAvailableReleaseName[dbApp.ReleaseName]; !found {
			deleteBatchQuery = append(deleteBatchQuery, &pb.BatchQuery{Cql: fmt.Sprintf(
				deleteCacheAppLaunchesQuery, "capten", orgID, clusterID, dbApp.ReleaseName)})
		}
	}

	// compare the give data with DB, if there is a mismatch do CUD
	for _, sentApp := range appLaunches {
		dbAppLaunches, found := dbAvailableReleaseName[sentApp.ReleaseName]
		if !found {
			insertBatchQuery = append(insertBatchQuery, &pb.BatchQuery{Cql: fmt.Sprintf(insertCacheAppLaunchesQuery,
				"capten", clusterID, orgID, sentApp.ReleaseName, sentApp.Description, sentApp.Category, sentApp.Icon, sentApp.LaunchURL, sentApp.LaunchUIDescription)})

			continue
		}

		// update the any mismatch data
		if sentApp.ReleaseName == dbAppLaunches.ReleaseName && (sentApp.Category != sentApp.Category ||
			sentApp.Description != dbAppLaunches.Description || string(sentApp.Icon) != string(dbAppLaunches.Icon) ||
			sentApp.LaunchUIDescription != dbAppLaunches.LaunchUIDescription || sentApp.LaunchURL != dbAppLaunches.LaunchURL) {
			updateBatchQuery = append(updateBatchQuery, &pb.BatchQuery{Cql: fmt.Sprintf(updateCacheAppLaunchesQuery, "capten",
				sentApp.Description, sentApp.Category, sentApp.Icon, sentApp.LaunchURL, sentApp.LaunchUIDescription, orgID, clusterID, sentApp.ReleaseName)})

		}

	}

	finalBatchQuery := append(insertBatchQuery, updateBatchQuery...)
	finalBatchQuery = append(finalBatchQuery, deleteBatchQuery...)
	fmt.Println("QUERY:", finalBatchQuery)
	_, err = a.ExecuteBatch(&pb.Batch{Queries: finalBatchQuery})

	return err
}

const (
	insertClusterQuery     = "INSERT INTO %s.capten_clusters (cluster_id, org_id, cluster_name, endpoint) VALUES (%s, %s, '%s', '%s');"
	updateClusterQuery     = "UPDATE %s.capten_clusters SET cluster_name='%s', endpoint='%s' WHERE org_id=%s AND cluster_id=%s;"
	deleteClusterQuery     = "DELETE FROM %s.capten_clusters WHERE org_id=%s AND cluster_id=%s;"
	getClusterDetailsQuery = "SELECT cluster_name, endpoint FROM %s.capten_clusters WHERE org_id=%s AND cluster_id=%s;"
	getClustersForOrgQuery = "SELECT cluster_id, cluster_name, endpoint FROM %s.capten_clusters WHERE org_id=%s;"

	getCacheAppLaunchesQuery = `SELECT release_name, description, category, icon, launch_url, launch_ui_description 
	 FROM %s.app_launches WHERE org_id=%s AND cluster_id=%s;`
	insertCacheAppLaunchesQuery = `INSERT INTO %s.app_launches 
	(cluster_id, org_id, release_name, description, category, icon, launch_url, launch_ui_description)
	VALUES (%s, %s, '%s', '%s', '%s', textAsBlob('%s'),'%s', '%s');`
	updateCacheAppLaunchesQuery = `UPDATE %s.app_launches SET description='%s', category='%s', 
	icon=textAsBlob('%s'), launch_url='%s', launch_ui_description='%s' WHERE org_id=%s AND cluster_id=%s AND release_name='%s';`
	deleteCacheAppLaunchesQuery = `DELETE FROM %s.app_launches 
	WHERE org_id=%s AND cluster_id=%s AND release_name='%s';`
	deleteFullCacheAppLaunchesQuery = `DELETE FROM %s.app_launches WHERE org_id=%s AND cluster_id=%s;`
)

func GetCacheAppLaunches(a *client.StargateClient, orgID, clusterID string) (*agentpb.GetClusterAppLaunchesResponse, error) {
	q := &pb.Query{
		Cql: fmt.Sprintf(getCacheAppLaunchesQuery, "capten", orgID, clusterID),
	}

	response, err := a.ExecuteQuery(q)
	if err != nil {
		return nil, fmt.Errorf("failed get cache cluster app launches: %w", err)
	}

	result := response.GetResultSet()

	clusterAppLaunches := make([]*agentpb.AppLaunchConfig, len(result.Rows))
	for index, row := range result.Rows {
		releaseName, err := client.ToString(row.Values[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterID: %w", err)
		}

		description, err := client.ToString(row.Values[1])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterName: %w", err)
		}

		category, err := client.ToString(row.Values[2])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterEndpoint: %w", err)
		}

		icon, err := client.ToBlob(row.Values[3])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterEndpoint: %w", err)
		}

		launchUrl, err := client.ToString(row.Values[4])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterEndpoint: %w", err)
		}

		launchUrlDesc, err := client.ToString(row.Values[5])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterEndpoint: %w", err)
		}

		clusterAppLaunches[index] = &agentpb.AppLaunchConfig{ReleaseName: releaseName, Category: category, Icon: icon,
			Description: description, LaunchURL: launchUrl,
			LaunchUIDescription: launchUrlDesc}
	}

	return &agentpb.GetClusterAppLaunchesResponse{Status: agentpb.StatusCode_OK, StatusMessage: "Successfully fetched the cluster launches",
		LaunchConfigList: clusterAppLaunches}, nil
}

func main() {
	fmt.Println("Start Query")
	sess, err := NewClient()
	if err != nil {
		fmt.Println("test, ", err)
		return
	}

	fmt.Println(UpdateCacheAppLaunches(sess, "996162a1-1df7-44b9-8347-1cb1acc78888",
		"f4d316da-5174-11ee-bf43-e6f881eee158", []*agentpb.AppLaunchConfig{{ReleaseName: "test"}, {ReleaseName: "prometheus"}}))

	// descQuery := "select * from capten.app_launches;"
	// resp, err := sess.ExecuteQuery(&pb.Query{Cql: descQuery})
	// if err != nil {
	// 	fmt.Println("failed to execute query: ", err)
	// 	return
	// }

	// fmt.Println("select response:", resp.GetResultSet().Rows)
	// //c, err := cl.GetClusters("996162a1-1df7-44b9-8347-1cb1acc78888")
	// fmt.Println("End Query")

}
