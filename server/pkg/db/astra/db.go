package astra

import (
	"crypto/tls"
	"fmt"
	"strings"
	"sync"

	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/types"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.uber.org/zap"
)

var (
	once     sync.Once
	astraObj *astra
)

type astra struct {
	session *client.StargateClient
}

func New() (*astra, error) {
	var err error
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	once.Do(func() {
		astraObj = &astra{}
		astraObj.session, err = connect()
		if err != nil {
			logger.Error("failed to connect to astra db", zap.Error(err))
		}
		err = astraObj.initializeDb()
		if err != nil {
			logger.Error("failed to initialize db")
		}
	})

	return astraObj, err
}

func connect() (*client.StargateClient, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Astra DB configuration example (dummy entries):
	//const astraUri = "0d175de3-c592-43f7-adf5-1bdda2761385-us-east1.apps.astra.datastax.com:443"
	//const bearerToken = "AstraCS:kYZPvIeLpthElpvKXQZUWHZF:32613fec5fe0be7f3cff755c2a09c5a411f0b0516d5521fc1fe8f3cbb3bf74ef"
	cfg := config.GetConfig()
	host := cfg.GetString("server.dbHost")
	password := cfg.GetString("server.dbPassword")
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(password),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to astra db: %w ", err)
	}

	stargateClient, err := client.NewStargateClientWithConn(conn)
	if err != nil {
		logger.Error("error creating stargate client", zap.Error(err))
		return nil, err
	}

	return stargateClient, nil
}

func (a *astra) initializeDb() error {
	initDbQueries := []string{
		createClusterEndpointTableQuery,
		createOrgClusterTableQuery,
	}

	for _, query := range initDbQueries {
		createQuery := &pb.Query{
			Cql: query,
		}

		_, err := a.session.ExecuteQuery(createQuery)
		if err != nil {
			return fmt.Errorf("failed to initialise db: %w", err)
		}
	}

	return nil
}

func (a *astra) GetClusterEndpoint(orgID, clusterName string) (string, error) {
	clusterId, err := a.getClusterID(orgID, clusterName)
	if err != nil {
		return "", fmt.Errorf("failed to get cluster endpoint: %w", err)
	}

	endpointQuery := &pb.Query{
		Cql: fmt.Sprintf("Select endpoint FROM %s.cluster_endpoint WHERE cluster_id=%s;",
			keyspace, clusterId),
	}

	response, err := a.session.ExecuteQuery(endpointQuery)
	result := response.GetResultSet()

	if len(result.Rows) == 0 {
		return "", fmt.Errorf("cluster: %s not found", clusterName)
	}

	endpoint, err := client.ToString(result.Rows[0].Values[0])
	if err != nil {
		return "", fmt.Errorf("cluster: %s unable to convert endpoint to string", clusterName)
	}

	return endpoint, nil
}

func (a *astra) RegisterCluster(orgId, clusterName, endpoint string) error {
	clusterExists, err := a.clusterEntryExists(orgId)
	if err != nil {
		return fmt.Errorf("failed to store cluster details: %w", err)
	}

	clusterId := gocql.TimeUUID()
	var batchQueries []*pb.BatchQuery
	if clusterExists {
		batchQueries = append(batchQueries,
			&pb.BatchQuery{
				Cql: fmt.Sprintf(
					"UPDATE %s.org_cluster SET cluster_ids= cluster_ids + {%s} WHERE org_id=%s;",
					keyspace, clusterId.String(), orgId),
			})
	} else {
		batchQueries = append(batchQueries,
			&pb.BatchQuery{
				Cql: fmt.Sprintf("INSERT INTO %s.org_cluster(org_id, cluster_ids) VALUES (%s, {%s});",
					keyspace, orgId, clusterId),
			})
	}

	batchQueries = append(batchQueries,
		&pb.BatchQuery{
			Cql: fmt.Sprintf("INSERT INTO %s.cluster_endpoint (cluster_id, org_id, cluster_name, endpoint) VALUES (%s, %s, '%s', '%s');",
				keyspace, clusterId, orgId, clusterName, endpoint),
		})

	batch := &pb.Batch{
		Type:    pb.Batch_LOGGED,
		Queries: batchQueries,
	}

	_, err = a.session.ExecuteBatch(batch)
	if err != nil {
		return fmt.Errorf("failed store cluster details %w", err)
	}

	return nil
}

func (a *astra) clusterEntryExists(orgID string) (bool, error) {
	selectClusterQuery := &pb.Query{
		Cql: fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
			keyspace, orgID),
	}

	response, err := a.session.ExecuteQuery(selectClusterQuery)
	if err != nil {
		return false, fmt.Errorf("failed to initialise db: %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) > 0 {
		return true, nil
	}

	return false, nil
}

func (a *astra) UpdateCluster(orgID, clusterName, endpoint string) error {
	clusterId, err := a.getClusterID(orgID, clusterName)
	if err != nil {
		return fmt.Errorf("failed to update the cluster info: %w", err)
	}

	updateQuery := &pb.Query{
		Cql: fmt.Sprintf(
			"UPDATE %s.cluster_endpoint set endpoint='%s' WHERE cluster_id=%s AND org_id=%s",
			keyspace, endpoint, clusterId, orgID),
	}

	_, err = a.session.ExecuteQuery(updateQuery)
	if err != nil {
		return fmt.Errorf("failed to update cluster info: %w", err)
	}

	return nil
}

func (a *astra) DeleteCluster(orgID, clusterName string) error {
	clusterId, err := a.getClusterID(orgID, clusterName)
	if err != nil {
		return fmt.Errorf("failed to delete cluster info: %w", err)
	}

	batch := &pb.Batch{
		Type: pb.Batch_LOGGED,
		Queries: []*pb.BatchQuery{
			{
				Cql: fmt.Sprintf(
					"DELETE FROM %s.cluster_endpoint WHERE cluster_id=%s ;",
					keyspace, clusterId),
			},
			{
				Cql: fmt.Sprintf(
					"UPDATE %s.org_cluster set cluster_ids = cluster_ids - {%s} WHERE org_id=%s ;",
					keyspace, clusterId, orgID),
			},
		},
	}

	_, err = a.session.ExecuteBatch(batch)
	if err != nil {
		return fmt.Errorf("failed delete cluster details %w", err)
	}

	return nil
}

func (a *astra) GetClusters(orgID string) ([]types.ClusterDetails, error) {
	selectClusterIdsQuery := &pb.Query{
		Cql: fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
			keyspace, orgID),
	}

	response, err := a.session.ExecuteQuery(selectClusterIdsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info from db: %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) == 0 {
		return []types.ClusterDetails{}, nil
	}

	clusterIds, err := client.ToSet(result.Rows[0].Values[0], UuidSetSpec)
	if err != nil {
		return nil, err
	}

	var clusterIdStrs []string
	for _, clusterId := range clusterIds.([]interface{}) {
		clusterIdStrs = append(clusterIdStrs, clusterId.(*uuid.UUID).String())
	}

	selectClustersQuery := &pb.Query{
		Cql: fmt.Sprintf("Select cluster_name, endpoint FROM %s.cluster_endpoint WHERE cluster_id in (%s);",
			keyspace, strings.Join(clusterIdStrs, ",")),
	}

	response, err = a.session.ExecuteQuery(selectClustersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster endpoint info from db: %w", err)
	}

	result = response.GetResultSet()
	var clusterDetails []types.ClusterDetails
	for _, row := range result.Rows {
		cqlClusterName, err := client.ToString(row.Values[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster name: %w", err)
		}

		cqlEndpoint, err := client.ToString(row.Values[1])
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster endpoint: %w", err)
		}

		clusterDetails = append(clusterDetails,
			types.ClusterDetails{
				ClusterName: cqlClusterName,
				Endpoint:    cqlEndpoint,
			})
	}

	return clusterDetails, nil
}

func (a *astra) getClusterID(orgID, clusterName string) (string, error) {
	selectClusterIdsQuery := &pb.Query{
		Cql: fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
			keyspace, orgID),
	}

	response, err := a.session.ExecuteQuery(selectClusterIdsQuery)
	if err != nil {
		return "", fmt.Errorf("failed to get cluster info from db: %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) == 0 {
		return "", fmt.Errorf("no cluster found for orgId %s", orgID)
	}

	clusterIds, err := client.ToSet(result.Rows[0].Values[0], UuidSetSpec)
	if err != nil {
		return "", err
	}

	var clusterIdStrs []string
	for _, clusterId := range clusterIds.([]interface{}) {
		clusterIdStrs = append(clusterIdStrs, clusterId.(*uuid.UUID).String())
	}

	selectClustersQuery := &pb.Query{
		Cql: fmt.Sprintf("Select cluster_id, cluster_name FROM %s.cluster_endpoint WHERE cluster_id in (%s);",
			keyspace, strings.Join(clusterIdStrs, ",")),
	}

	response, err = a.session.ExecuteQuery(selectClustersQuery)
	if err != nil {
		return "", fmt.Errorf("failed to get cluster endpoint info from db: %w", err)
	}

	result = response.GetResultSet()
	for _, row := range result.Rows {
		cqlClusterId, err := client.ToUUID(row.Values[0])
		if err != nil {
			return "", fmt.Errorf("failed to get cluster uuid: %w", err)
		}

		cqlClusterName, err := client.ToString(row.Values[1])
		if err != nil {
			return "", fmt.Errorf("failed to get cluster name: %w", err)
		}

		if cqlClusterName == clusterName {
			return cqlClusterId.String(), nil
		}
	}

	return "", fmt.Errorf("cluster not found")
}
