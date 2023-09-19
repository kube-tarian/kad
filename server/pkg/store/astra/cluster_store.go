package astra

import (
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

const (
	insertClusterQuery     = "INSERT INTO %s.capten_clusters (cluster_id, org_id, cluster_name, endpoint) VALUES (%s, %s, '%s', '%s');"
	updateClusterQuery     = "UPDATE %s.capten_clusters SET cluster_name='%s', endpoint='%s' WHERE org_id=%s AND cluster_id=%s;"
	deleteClusterQuery     = "DELETE FROM %s.capten_clusters WHERE org_id=%s AND cluster_id=%s;"
	getClusterDetailsQuery = "SELECT cluster_name, endpoint FROM %s.capten_clusters WHERE org_id=%s AND cluster_id=%s;"
	getClustersForOrgQuery = "SELECT cluster_id, cluster_name, endpoint FROM %s.capten_clusters WHERE org_id=%s;"

	getClusterAppLaunchesQuery = `SELECT release_name, description, category, icon, launch_url, launch_ui_description 
	 FROM %s.app_launches WHERE org_id=%s AND cluster_id=%s;`
	insertClusterAppLaunchesQuery = `INSERT INTO %s.app_launches 
	(cluster_id, org_id, release_name, description, category, icon, launch_url, launch_ui_description)
	VALUES (%s, %s, '%s', '%s', '%s', textAsBlob('%s'),'%s', '%s');`
	updateClusterAppLaunchesQuery = `UPDATE %s.app_launches SET description='%s', category='%s', 
	icon=textAsBlob('%s'), launch_url='%s', launch_ui_description='%s' WHERE org_id=%s AND cluster_id=%s AND release_name='%s';`
	deleteClusterAppLaunchesQuery = `DELETE FROM %s.app_launches 
	WHERE org_id=%s AND cluster_id=%s AND release_name='%s';`
	deleteFullClusterAppLaunchesQuery = `DELETE FROM %s.app_launches WHERE org_id=%s AND cluster_id=%s;`
)

func (a *AstraServerStore) GetClusterAppLaunches(orgID, clusterID string) (*agentpb.GetClusterAppLaunchesResponse, error) {
	q := &pb.Query{
		Cql: fmt.Sprintf(getClusterAppLaunchesQuery, a.keyspace, orgID, clusterID),
	}

	response, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return nil, fmt.Errorf("failed get Cluster cluster app launches: %w", err)
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

func (a *AstraServerStore) DeleteFullClusterAppLaunches(orgID, clusterID string) error {
	q := &pb.Query{
		Cql: fmt.Sprintf(deleteFullClusterAppLaunchesQuery, a.keyspace, orgID, clusterID),
	}

	_, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return fmt.Errorf("failed get Cluster cluster app launches: %w", err)
	}
	return nil
}

func (a *AstraServerStore) deleteClusterAppLaunches(orgID, clusterID, appLaunches []*agentpb.AppLaunchConfig) error {
	if len(appLaunches) == 0 {
		return nil
	}

	batchQuery := make([]*pb.BatchQuery, len(appLaunches))
	for index, app := range appLaunches {
		batchQuery[index] = &pb.BatchQuery{Cql: fmt.Sprintf(deleteClusterAppLaunchesQuery, a.keyspace, orgID, clusterID, app.ReleaseName)}
	}
	_, err := a.c.Session().ExecuteBatch(&pb.Batch{Queries: batchQuery})

	return err
}

func (a *AstraServerStore) InsertClusterAppLaunches(orgID, clusterID string, appLaunches []*agentpb.AppLaunchConfig) error {
	if len(appLaunches) == 0 {
		return nil
	}

	batchQuery := make([]*pb.BatchQuery, len(appLaunches))
	for index, app := range appLaunches {
		batchQuery[index] = &pb.BatchQuery{Cql: fmt.Sprintf(insertClusterAppLaunchesQuery, a.keyspace,
			clusterID, orgID, app.ReleaseName, app.Description, app.Category, app.Icon, app.LaunchURL, app.LaunchUIDescription)}
	}
	_, err := a.c.Session().ExecuteBatch(&pb.Batch{Queries: batchQuery})

	return err
}

func (a *AstraServerStore) UpdateClusterAppLaunches(orgID, clusterID string, appLaunches []*agentpb.AppLaunchConfig) error {
	appResponse, err := a.GetClusterAppLaunches(orgID, clusterID)
	if err != nil {
		return fmt.Errorf("failed to update the applaunches Cluster, err %w", err)
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
				deleteClusterAppLaunchesQuery, a.keyspace, orgID, clusterID, dbApp.ReleaseName)})
		}
	}

	// compare the give data with DB, if there is a mismatch do CUD
	for _, sentApp := range appLaunches {
		dbAppLaunches, found := dbAvailableReleaseName[sentApp.ReleaseName]
		if !found {
			insertBatchQuery = append(insertBatchQuery, &pb.BatchQuery{Cql: fmt.Sprintf(insertClusterAppLaunchesQuery,
				a.keyspace, clusterID, orgID, sentApp.ReleaseName, sentApp.Description, sentApp.Category, sentApp.Icon, sentApp.LaunchURL, sentApp.LaunchUIDescription)})

			continue
		}

		// update the any mismatch data
		if sentApp.ReleaseName == dbAppLaunches.ReleaseName && (sentApp.Category != sentApp.Category ||
			sentApp.Description != dbAppLaunches.Description || string(sentApp.Icon) != string(dbAppLaunches.Icon) ||
			sentApp.LaunchUIDescription != dbAppLaunches.LaunchUIDescription || sentApp.LaunchURL != dbAppLaunches.LaunchURL) {
			updateBatchQuery = append(updateBatchQuery, &pb.BatchQuery{Cql: fmt.Sprintf(updateClusterAppLaunchesQuery, a.keyspace,
				sentApp.Description, sentApp.Category, sentApp.Icon, sentApp.LaunchURL, sentApp.LaunchUIDescription, orgID, clusterID, sentApp.ReleaseName)})

		}

	}

	if len(insertBatchQuery) == 0 && len(updateBatchQuery) == 0 && len(deleteBatchQuery) == 0 {
		return nil
	}

	finalBatchQuery := append(insertBatchQuery, updateBatchQuery...)
	finalBatchQuery = append(finalBatchQuery, deleteBatchQuery...)
	_, err = a.c.Session().ExecuteBatch(&pb.Batch{Queries: finalBatchQuery})

	return err
}

func (a *AstraServerStore) AddCluster(orgID, clusterID, clusterName, endpoint string) error {
	q := &pb.Query{
		Cql: fmt.Sprintf(insertClusterQuery, a.keyspace, clusterID, orgID, clusterName, endpoint),
	}

	_, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return fmt.Errorf("failed store cluster details %w", err)
	}
	return nil
}

func (a *AstraServerStore) UpdateCluster(orgID, clusterID, clusterName, endpoint string) error {
	q := &pb.Query{
		Cql: fmt.Sprintf(updateClusterQuery, a.keyspace, clusterName, endpoint, orgID, clusterID),
	}

	_, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return fmt.Errorf("failed to update cluster info: %w", err)
	}
	return nil
}

func (a *AstraServerStore) DeleteCluster(orgID, clusterID string) error {
	q := &pb.Query{
		Cql: fmt.Sprintf(deleteClusterQuery, a.keyspace, orgID, clusterID),
	}

	_, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return fmt.Errorf("failed to update cluster info: %w", err)
	}
	return nil
}

func (a *AstraServerStore) GetClusterDetails(orgID, clusterID string) (*types.ClusterDetails, error) {
	q := &pb.Query{
		Cql: fmt.Sprintf(getClusterDetailsQuery, a.keyspace, orgID, clusterID),
	}

	response, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return nil, err
	}
	result := response.GetResultSet()

	if len(result.Rows) != 1 {
		return nil, fmt.Errorf("cluster: %s not found", clusterID)
	}

	clusterName, err := client.ToString(result.Rows[0].Values[0])
	if err != nil {
		return nil, fmt.Errorf("cluster: %s unable to convert clusterName to string", clusterID)
	}

	clusterEndpoint, err := client.ToString(result.Rows[0].Values[1])
	if err != nil {
		return nil, fmt.Errorf("cluster: %s unable to convert endpoint to string", clusterID)
	}

	return &types.ClusterDetails{
		ClusterID: clusterID, OrgID: orgID,
		ClusterName: clusterName, Endpoint: clusterEndpoint}, nil
}

func (a *AstraServerStore) GetClusters(orgID string) ([]types.ClusterDetails, error) {
	q := &pb.Query{
		Cql: fmt.Sprintf(getClustersForOrgQuery, a.keyspace, orgID),
	}

	response, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info from db: %w", err)
	}

	result := response.GetResultSet()
	var clusterDetails []types.ClusterDetails
	for _, row := range result.Rows {
		cqlClusterID, err := client.ToUUID(row.Values[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterID: %w", err)
		}

		cqlClusterName, err := client.ToString(row.Values[1])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterName: %w", err)
		}

		cqlEndpoint, err := client.ToString(row.Values[2])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterEndpoint: %w", err)
		}

		clusterDetails = append(clusterDetails,
			types.ClusterDetails{
				OrgID:       orgID,
				ClusterID:   cqlClusterID.String(),
				ClusterName: cqlClusterName, Endpoint: cqlEndpoint,
			})
	}

	return clusterDetails, nil
}
