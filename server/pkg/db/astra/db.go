package astra

import (
	"fmt"
	"sync"
	"time"

	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"

	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/types"

	"go.uber.org/zap"
)

var (
	once     sync.Once
	astraObj *astra
)

type astra struct {
	session *gocql.Session
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

func connect() (*gocql.Session, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var cluster *gocql.ClusterConfig
	// Astra DB configuration
	//const astraUri = "0d175de3-c592-43f7-adf5-1bdda2761385-us-east1.apps.astra.datastax.com:443"
	//const bearerToken = "AstraCS:kYZPvIeLpthElpvKXQZUWHZF:32613fec5fe0be7f3cff755c2a09c5a411f0b0516d5521fc1fe8f3cbb3bf74ef"
	cfg := config.GetConfig()
	host := cfg.GetString("server.dbHost")
	user := cfg.GetString("server.dbUsername")
	password := cfg.GetString("server.dbPassword")

	// Create connection with authentication
	// For Astra DB:
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: false,
	//}

	//conn, err := grpc.Dial(host, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	//	grpc.WithBlock(),
	//	grpc.WithPerRPCCredentials(
	//		auth.NewStaticTokenProvider(bearerToken),
	//	),
	//)
	//
	//stargateClient, err := client.NewStargateClientWithConn(conn)
	//if err != nil {
	//	logger.Error("error creating stargate client", zap.Error(err))
	//	return nil, err
	//}
	cluster, err := gocqlastra.NewClusterFromURL(host,
		user, password, 50*time.Second)

	if err != nil {
		return nil, err
	}

	cluster.Timeout = 80 * time.Second
	session, err := gocql.NewSession(*cluster)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (a *astra) initializeDb() error {
	if err := a.session.Query(createKeyspaceQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace, %w", err)
	}

	if err := a.session.Query(createClusterEndpointTableQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create cluster_endpoint table, %w", err)
	}

	if err := a.session.Query(createClusterEndpointTableQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create cluster_endpoint table, %w", err)
	}

	return nil
}

func (a *astra) GetClusterEndpoint(orgID, clusterName string) (string, error) {
	//selectQuery := &pb.Query{
	//	Cql: fmt.Sprintf("SELECT endpoint FROM %s.cluster WHERE org_id = %s AND name = %s", keyspace, orgID, clusterName),
	//}

	return "", nil
}

func (a *astra) RegisterCluster(orgId, clusterName, endpoint string) error {
	clusterId, err := a.createCluster(orgId)
	if err != nil {
		return err
	}

	err = a.session.Query(fmt.Sprintf("INSERT INTO %s.cluster_endpoint (cluster_id, org_id, cluster_name, endpoint) VALUES (%s, %s, '%s', '%s');",
		keyspace, clusterId, orgId, clusterName, endpoint)).Exec()
	if err != nil {
		return fmt.Errorf("failed insert cluster details %w", err)
	}

	return nil
}

func (a *astra) createCluster(orgID string) (string, error) {
	clusterId := gocql.TimeUUID()
	iter := a.session.Query(fmt.Sprintf("Select * FROM %s.org_cluster WHERE org_id=%s;",
		keyspace, orgID)).Iter()

	clusterIds := make([]string, 0)
	iter.Scan(&clusterIds)
	if len(clusterIds) == 0 {
		err := a.session.Query(
			fmt.Sprintf("INSERT INTO %s.org_cluster(org_id, cluster_ids) VALUES (%s, {%s});",
				keyspace, orgID, clusterId),
		).Exec()
		return clusterId.String(), err
	}

	err := a.session.Query(
		fmt.Sprintf("UPDATE %s.org_cluster SET cluster_ids= cluster_ids + {%s};", keyspace, clusterId.String())).
		Exec()

	if err != nil {
		return "", err
	}

	return clusterId.String(), nil
}

func (a *astra) UpdateCluster(orgID, clusterID, endpoint string) error {

	return nil
}

func (a *astra) DeleteCluster(orgID, clusterID string) error {

	return nil
}

func (a *astra) GetClusters(orgID string) ([]types.ClusterDetails, error) {

	return nil, nil
}
