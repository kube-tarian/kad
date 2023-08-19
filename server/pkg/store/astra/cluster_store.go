package astra

import (
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

const (
	insertClusterQuery     = "INSERT INTO %s.capten_clusters (cluster_id, org_id, cluster_name, endpoint) VALUES (%s, %s, '%s', '%s');"
	updateClusterQuery     = "UPDATE %s.capten_clusters SET cluster_name='%s', endpoint='%s' WHERE cluster_id=%s AND org_id=%s;"
	deleteClusterQuery     = "DELETE FROM %s.capten_clusters WHERE cluster_id=%s AND org_id=%s;"
	getClusterDetailsQuery = "SELECT cluster_name, endpoint FROM %s.capten_clusters WHERE cluster_id=%s;"
	getClustersForOrgQuery = "SELECT org_id, cluster_id, cluster_name, endpoint FROM %s.capten_clusters WHERE org_id=%s ALLOW FILTERING;"
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
		Cql: fmt.Sprintf(updateClusterQuery, a.keyspace, clusterName, endpoint, clusterID, orgID),
	}

	_, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return fmt.Errorf("failed to update cluster info: %w", err)
	}
	return nil
}

func (a *AstraServerStore) DeleteCluster(orgID, clusterID string) error {
	q := &pb.Query{
		Cql: fmt.Sprintf(deleteClusterQuery, a.keyspace, clusterID, orgID),
	}

	_, err := a.c.Session().ExecuteQuery(q)
	if err != nil {
		return fmt.Errorf("failed to update cluster info: %w", err)
	}
	return nil
}

func (a *AstraServerStore) GetClusterDetails(clusterID string) (*types.ClusterDetail, error) {
	q := &pb.Query{
		Cql: fmt.Sprintf(getClusterDetailsQuery, a.keyspace, clusterID),
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

	return &types.ClusterDetail{ClusterName: clusterName, Endpoint: clusterEndpoint}, nil
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
		cqlOrgID, err := client.ToUUID(row.Values[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get orgID: %w", err)
		}

		cqlClusterID, err := client.ToUUID(row.Values[1])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterID: %w", err)
		}

		cqlClusterName, err := client.ToString(row.Values[2])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterName: %w", err)
		}

		cqlEndpoint, err := client.ToString(row.Values[3])
		if err != nil {
			return nil, fmt.Errorf("failed to get clusterEndpoint: %w", err)
		}

		clusterDetails = append(clusterDetails,
			types.ClusterDetails{
				OrgID:         cqlOrgID.String(),
				ClusterID:     cqlClusterID.String(),
				ClusterDetail: types.ClusterDetail{ClusterName: cqlClusterName, Endpoint: cqlEndpoint},
			})
	}

	return clusterDetails, nil
}
