package cassandra

import (
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	cassandraclient "github.com/kube-tarian/kad/server/pkg/cassandra-client"
	"github.com/kube-tarian/kad/server/pkg/types"
)

type CassandraServerStore struct {
	c        *cassandraclient.Client
	keyspace string
}

func NewStore() (*CassandraServerStore, error) {
	cs := &CassandraServerStore{}
	err := cs.initSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cassandra db, %w", err)
	}
	cs.keyspace = cs.c.Keyspace()
	return cs, err
}

func (c *CassandraServerStore) initSession() error {
	var err error
	c.c, err = cassandraclient.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create cassandra session")
	}
	return nil
}

func (c *CassandraServerStore) InitializeDb() error {
	if err := c.c.Session().Query(fmt.Sprintf(createKeyspaceQuery, c.keyspace)).Exec(); err != nil {
		return fmt.Errorf("failed to create c.keyspace, %w", err)
	}

	if err := c.c.Session().Query(fmt.Sprintf(createClusterEndpointTableQuery, c.keyspace)).Exec(); err != nil {
		return fmt.Errorf("failed to create cluster_endpoint table, %w", err)
	}

	if err := c.c.Session().Query(fmt.Sprintf(createOrgClusterTableQuery, c.keyspace)).Exec(); err != nil {
		return fmt.Errorf("failed to create cluster_endpoint table, %w", err)
	}

	return nil
}

func (c *CassandraServerStore) GetClusterEndpoint(orgID, clusterName string) (string, error) {
	clusterId, err := c.getClusterID(orgID, clusterName)
	if err != nil {
		return "", err
	}

	iter := c.c.Session().Query(fmt.Sprintf("Select endpoint FROM %s.cluster_endpoint WHERE cluster_id=%s;",
		c.keyspace, clusterId)).Iter()
	var endpoint string
	iter.Scan(&endpoint)
	return endpoint, nil
}

func (c *CassandraServerStore) GetClusters(orgId string) ([]types.ClusterDetails, error) {
	iter := c.c.Session().Query(fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
		c.keyspace, orgId)).Iter()
	var clusterUUIds []gocql.UUID
	iter.Scan(&clusterUUIds)
	var clusterIds []string
	for _, id := range clusterUUIds {
		clusterIds = append(clusterIds, id.String())
	}

	iter = c.c.Session().Query(fmt.Sprintf("Select cluster_name, endpoint FROM %s.cluster_endpoint WHERE cluster_id in (%s);",
		c.keyspace, strings.Join(clusterIds, ","))).Iter()

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

func (c *CassandraServerStore) AddCluster(orgId, clusterName, endpoint string) error {
	clusterExists := c.clusterEntryExists(orgId)
	clusterId := gocql.TimeUUID()
	batch := c.c.Session().NewBatch(gocql.LoggedBatch)
	if clusterExists {
		batch.Query(
			fmt.Sprintf(
				"UPDATE %s.org_cluster SET cluster_ids= cluster_ids + {%s} WHERE org_id=%s;",
				c.keyspace, clusterId.String(), orgId))
	} else {
		batch.Query(
			fmt.Sprintf("INSERT INTO %s.org_cluster(org_id, cluster_ids) VALUES (%s, {%s});",
				c.keyspace, orgId, clusterId),
		)
	}

	batch.Query(fmt.Sprintf("INSERT INTO %s.cluster_endpoint (cluster_id, org_id, cluster_name, endpoint) VALUES (%s, %s, '%s', '%s');",
		c.keyspace, clusterId, orgId, clusterName, endpoint))
	err := c.c.Session().ExecuteBatch(batch)
	if err != nil {
		return fmt.Errorf("failed insert cluster details %w", err)
	}

	return nil
}

func (c *CassandraServerStore) clusterEntryExists(orgID string) bool {
	iter := c.c.Session().Query(fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
		c.keyspace, orgID)).Iter()
	var clusterIds []gocql.UUID
	iter.Scan(&clusterIds)
	if len(clusterIds) == 0 {
		return false
	}

	return true
}

func (c *CassandraServerStore) UpdateCluster(orgID, clusterName, endpoint string) error {
	clusterId, err := c.getClusterID(orgID, clusterName)
	if err != nil {
		return err
	}

	err = c.c.Session().Query(fmt.Sprintf(
		"UPDATE %s.cluster_endpoint set endpoint='%s' WHERE cluster_id=%s AND org_id=%s",
		c.keyspace, endpoint, clusterId, orgID)).Exec()
	return err
}

func (c *CassandraServerStore) getClusterID(orgID, clusterName string) (string, error) {
	iter := c.c.Session().Query(fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s;",
		c.keyspace, orgID)).Iter()

	var clusterUUIds []gocql.UUID
	iter.Scan(&clusterUUIds)
	var clusterIds []string
	for _, id := range clusterUUIds {
		clusterIds = append(clusterIds, id.String())
	}

	iter = c.c.Session().Query(fmt.Sprintf("Select cluster_id, cluster_name FROM %s.cluster_endpoint WHERE cluster_id in (%s);",
		c.keyspace, strings.Join(clusterIds, ","))).Iter()

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

func (c *CassandraServerStore) DeleteCluster(orgID, clusterName string) error {
	clusterId, err := c.getClusterID(orgID, clusterName)
	if err != nil {
		return err
	}

	batch := c.c.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(
		"DELETE FROM %s.cluster_endpoint WHERE cluster_id=%s ;",
		c.keyspace, clusterId))
	batch.Query(fmt.Sprintf(
		"UPDATE %s.org_cluster set cluster_ids = cluster_ids - {%s} WHERE org_id=%s ;",
		c.keyspace, clusterId, orgID))
	return c.c.Session().ExecuteBatch(batch)
}
