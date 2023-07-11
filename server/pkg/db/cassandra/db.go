package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/types"
	"log"
	"strings"
	"sync"
)

type cassandra struct {
	session *gocql.Session
}

var (
	cassandraSession *cassandra
	once             sync.Once
)

func New() (*cassandra, error) {
	var err error
	once.Do(func() {
		cassandraSession = &cassandra{}
		cassandraSession.session, err = connect()
		if err != nil {
			log.Println("failed to connect to cassandra")
		}

		err = cassandraSession.initializeDb()
		if err != nil {
			log.Println("failed to initialize db")
		}
	})

	return cassandraSession, err
}

func connect() (*gocql.Session, error) {
	cfg := config.GetConfig()
	host := cfg.GetString("server.dbHost")
	user := cfg.GetString("server.dbUsername")
	password := cfg.GetString("server.dbPassword")
	cluster := gocql.NewCluster(host)
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: user,
		Password: password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create c")
	}

	return session, nil
}

func (c *cassandra) initializeDb() error {
	if err := c.session.Query(createKeyspaceQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace, %w", err)
	}

	if err := c.session.Query(createClusterEndpointTableQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create cluster_endpoint table, %w", err)
	}

	if err := c.session.Query(createOrgClusterTableQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create cluster_endpoint table, %w", err)
	}

	return nil
}

func (c *cassandra) GetClusterEndpoint(orgID, clusterName string) (string, error) {
	clusterId, err := c.getClusterID(orgID, clusterName)
	if err != nil {
		return "", err
	}

	iter := c.session.Query(fmt.Sprintf("Select endpoint FROM %s.cluster_endpoint WHERE cluster_id=%s;",
		keyspace, clusterId)).Iter()
	var endpoint string
	iter.Scan(&endpoint)
	return endpoint, nil
}

func (c *cassandra) GetClusters(orgId string) ([]types.ClusterDetails, error) {
	iter := c.session.Query(fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
		keyspace, orgId)).Iter()
	var clusterUUIds []gocql.UUID
	iter.Scan(&clusterUUIds)
	var clusterIds []string
	for _, id := range clusterUUIds {
		clusterIds = append(clusterIds, id.String())
	}

	iter = c.session.Query(fmt.Sprintf("Select cluster_name, endpoint FROM %s.cluster_endpoint WHERE cluster_id in (%s);",
		keyspace, strings.Join(clusterIds, ","))).Iter()

	var cqlClusterName string
	var cqlClusterEndpoint string
	cqlScanner := iter.Scanner()
	var clusterDetails []types.ClusterDetails
	for cqlScanner.Next() {
		if err := cqlScanner.Scan(&cqlClusterName, &cqlClusterEndpoint); err != nil {
			return nil, err
		}

		clusterDetails = append(clusterDetails, types.ClusterDetails{
			ClusterName: cqlClusterName,
			Endpoint:    cqlClusterEndpoint,
		})
	}

	return clusterDetails, nil
}

func (c *cassandra) RegisterCluster(orgId, clusterName, endpoint string) error {
	clusterExists := c.clusterEntryExists(orgId)
	clusterId := gocql.TimeUUID()
	batch := c.session.NewBatch(gocql.LoggedBatch)
	if clusterExists {
		batch.Query(
			fmt.Sprintf(
				"UPDATE %s.org_cluster SET cluster_ids= cluster_ids + {%s} WHERE org_id=%s;",
				keyspace, clusterId.String(), orgId))
	} else {
		batch.Query(
			fmt.Sprintf("INSERT INTO %s.org_cluster(org_id, cluster_ids) VALUES (%s, {%s});",
				keyspace, orgId, clusterId),
		)
	}

	batch.Query(fmt.Sprintf("INSERT INTO %s.cluster_endpoint (cluster_id, org_id, cluster_name, endpoint) VALUES (%s, %s, '%s', '%s');",
		keyspace, clusterId, orgId, clusterName, endpoint))
	err := c.session.ExecuteBatch(batch)
	if err != nil {
		return fmt.Errorf("failed insert cluster details %w", err)
	}

	return nil
}

func (c *cassandra) clusterEntryExists(orgID string) bool {
	iter := c.session.Query(fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
		keyspace, orgID)).Iter()
	var clusterIds []gocql.UUID
	iter.Scan(&clusterIds)
	if len(clusterIds) == 0 {
		return false
	}

	return true
}

func (c *cassandra) UpdateCluster(orgID, clusterName, endpoint string) error {
	clusterId, err := c.getClusterID(orgID, clusterName)
	if err != nil {
		return err
	}

	err = c.session.Query(fmt.Sprintf(
		"UPDATE %s.cluster_endpoint set endpoint='%s' WHERE cluster_id=%s AND org_id=%s",
		keyspace, endpoint, clusterId, orgID)).Exec()
	return err
}

func (c *cassandra) getClusterID(orgID, clusterName string) (string, error) {
	iter := c.session.Query(fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s;",
		keyspace, orgID)).Iter()

	var clusterUUIds []gocql.UUID
	iter.Scan(&clusterUUIds)
	var clusterIds []string
	for _, id := range clusterUUIds {
		clusterIds = append(clusterIds, id.String())
	}

	iter = c.session.Query(fmt.Sprintf("Select cluster_id, cluster_name FROM %s.cluster_endpoint WHERE cluster_id in (%s);",
		keyspace, strings.Join(clusterIds, ","))).Iter()

	var cqlClusterId gocql.UUID
	var cqlClusterName string
	cqlScanner := iter.Scanner()
	for cqlScanner.Next() {
		if err := cqlScanner.Scan(&cqlClusterId, &cqlClusterName); err != nil {
			return "", err
		}

		if cqlClusterName == clusterName {
			return cqlClusterId.String(), nil
		}
	}

	return "", fmt.Errorf("cluster not found")
}

func (c *cassandra) DeleteCluster(orgID, clusterName string) error {
	clusterId, err := c.getClusterID(orgID, clusterName)
	if err != nil {
		return err
	}

	batch := c.session.NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(
		"DELETE FROM %s.cluster_endpoint WHERE cluster_id=%s ;",
		keyspace, clusterId))
	batch.Query(fmt.Sprintf(
		"UPDATE %s.org_cluster set cluster_ids = cluster_ids - {%s} WHERE org_id=%s ;",
		keyspace, clusterId, orgID))
	return c.session.ExecuteBatch(batch)
}
