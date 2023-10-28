package astra

import (
	"fmt"

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
)

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

func (a *AstraServerStore) GetClusterForOrg(orgID string) (*types.ClusterDetails, error) {
	q := &pb.Query{
		Cql: fmt.Sprintf(getClustersForOrgQuery, a.keyspace, orgID),
	}

	response, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info from db: %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("no cluster found")
	}

	cqlClusterID, err := client.ToUUID(result.Rows[0].Values[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get clusterID: %w", err)
	}

	cqlClusterName, err := client.ToString(result.Rows[0].Values[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get clusterName: %w", err)
	}

	cqlEndpoint, err := client.ToString(result.Rows[0].Values[2])
	if err != nil {
		return nil, fmt.Errorf("failed to get clusterEndpoint: %w", err)
	}

	clusterDetails := types.ClusterDetails{
		OrgID:       orgID,
		ClusterID:   cqlClusterID.String(),
		ClusterName: cqlClusterName, Endpoint: cqlEndpoint,
	}
	return &clusterDetails, nil
}
