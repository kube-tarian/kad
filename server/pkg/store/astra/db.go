package astra

import (
	"fmt"
	"strings"
	"time"

	astraclient "github.com/kube-tarian/kad/server/pkg/astra-client"
	"github.com/kube-tarian/kad/server/pkg/types"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

const (
	createAppConfigQuery         string = "INSERT INTO %s.app_config (name, chart_name, repo_name,release_name, repo_url, namespace, version, create_namespace, privileged_namespace, launch_ui_url, launch_ui_redirect_url, category, icon, description, launch_ui_values, override_values, created_time, id) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s','%v', '%s' );"
	updateAppConfigQuery         string = "UPDATE %s.app_config SET chart_name = '%s', repo_name = '%s', repo_url = '%s', namespace = '%s', create_namespace = %t, privileged_namespace = %t, launch_ui_url = '%s', launch_ui_redirect_url = '%s', category = '%s', icon = '%s', description = '%s', launch_ui_values = '%s', override_values = '%s',last_updated_time='%v' WHERE name = '%s' AND version = '%s';"
	deleteAppConfigQuery         string = "DELETE FROM %s.app_config WHERE name='%s' AND version='%s';"
	getAppConfigQuery            string = "Select name,chart_name,repo_name,repo_url,namespace,version,create_namespace,privileged_namespace,launch_ui_url,launch_ui_redirect_url,category,icon,description,launch_ui_values,override_values, release_name FROM %s.app_config WHERE name='%s' AND version='%s';"
	getAllAppConfigsQuery        string = "Select name,chart_name,repo_name,repo_url,namespace,version,create_namespace,privileged_namespace,launch_ui_url,launch_ui_redirect_url,category,icon,description,launch_ui_values,override_values, release_name FROM %s.app_config;"
	appConfigExistanceCheckQuery string = "Select name, version FROM %s.app_config WHERE name='%s' AND version ='%s';"
)

type AstraServerStore struct {
	c        *astraclient.Client
	keyspace string
}

func NewStore() (*AstraServerStore, error) {
	a := &AstraServerStore{}
	err := a.initClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to astra db, %w", err)
	}
	a.keyspace = a.c.Keyspace()
	a.keyspace = "capten_server"
	return a, nil
}

func (a *AstraServerStore) initClient() error {
	var err error
	a.c, err = astraclient.NewClient()
	return err
}

func (a *AstraServerStore) InitializeDb() error {
	initDbQueries := []string{
		fmt.Sprintf(createClusterEndpointTableQuery, a.keyspace),
		fmt.Sprintf(createOrgClusterTableQuery, a.keyspace),
	}

	for _, query := range initDbQueries {
		createQuery := &pb.Query{
			Cql: query,
		}

		_, err := a.c.Session().ExecuteQuery(createQuery)
		if err != nil {
			return fmt.Errorf("failed to initialise db: %w", err)
		}
	}

	return nil
}

func (a *AstraServerStore) GetClusterEndpoint(orgID, clusterName string) (string, error) {
	clusterId, err := a.getClusterID(orgID, clusterName)
	if err != nil {
		return "", fmt.Errorf("failed to get cluster endpoint: %w", err)
	}

	endpointQuery := &pb.Query{
		Cql: fmt.Sprintf("Select endpoint FROM %s.cluster_endpoint WHERE cluster_id=%s;",
			a.keyspace, clusterId),
	}

	response, err := a.c.Session().ExecuteQuery(endpointQuery)
	if err != nil {
		return "", err
	}
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

func (a *AstraServerStore) AddCluster(orgId, clusterName, endpoint string) error {
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
					a.keyspace, clusterId.String(), orgId),
			})
	} else {
		batchQueries = append(batchQueries,
			&pb.BatchQuery{
				Cql: fmt.Sprintf("INSERT INTO %s.org_cluster(org_id, cluster_ids) VALUES (%s, {%s});",
					a.keyspace, orgId, clusterId),
			})
	}

	batchQueries = append(batchQueries,
		&pb.BatchQuery{
			Cql: fmt.Sprintf("INSERT INTO %s.cluster_endpoint (cluster_id, org_id, cluster_name, endpoint) VALUES (%s, %s, '%s', '%s');",
				a.keyspace, clusterId, orgId, clusterName, endpoint),
		})

	batch := &pb.Batch{
		Type:    pb.Batch_LOGGED,
		Queries: batchQueries,
	}

	_, err = a.c.Session().ExecuteBatch(batch)
	if err != nil {
		return fmt.Errorf("failed store cluster details %w", err)
	}

	return nil
}

func (a *AstraServerStore) clusterEntryExists(orgID string) (bool, error) {
	selectClusterQuery := &pb.Query{
		Cql: fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
			a.keyspace, orgID),
	}

	response, err := a.c.Session().ExecuteQuery(selectClusterQuery)
	if err != nil {
		return false, fmt.Errorf("failed to initialise db: %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) > 0 {
		return true, nil
	}

	return false, nil
}

func (a *AstraServerStore) UpdateCluster(orgID, clusterName, endpoint string) error {
	clusterId, err := a.getClusterID(orgID, clusterName)
	if err != nil {
		return fmt.Errorf("failed to update the cluster info: %w", err)
	}

	updateQuery := &pb.Query{
		Cql: fmt.Sprintf(
			"UPDATE %s.cluster_endpoint set endpoint='%s' WHERE cluster_id=%s AND org_id=%s",
			a.keyspace, endpoint, clusterId, orgID),
	}

	_, err = a.c.Session().ExecuteQuery(updateQuery)
	if err != nil {
		return fmt.Errorf("failed to update cluster info: %w", err)
	}

	return nil
}

func (a *AstraServerStore) DeleteCluster(orgID, clusterName string) error {
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
					a.keyspace, clusterId),
			},
			{
				Cql: fmt.Sprintf(
					"UPDATE %s.org_cluster set cluster_ids = cluster_ids - {%s} WHERE org_id=%s ;",
					a.keyspace, clusterId, orgID),
			},
		},
	}

	_, err = a.c.Session().ExecuteBatch(batch)
	if err != nil {
		return fmt.Errorf("failed delete cluster details %w", err)
	}

	return nil
}

func (a *AstraServerStore) GetClusters(orgID string) ([]types.ClusterDetails, error) {
	selectClusterIdsQuery := &pb.Query{
		Cql: fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
			a.keyspace, orgID),
	}

	response, err := a.c.Session().ExecuteQuery(selectClusterIdsQuery)
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
			a.keyspace, strings.Join(clusterIdStrs, ",")),
	}

	response, err = a.c.Session().ExecuteQuery(selectClustersQuery)
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

func (a *AstraServerStore) getClusterID(orgID, clusterName string) (string, error) {
	selectClusterIdsQuery := &pb.Query{
		Cql: fmt.Sprintf("Select cluster_ids FROM %s.org_cluster WHERE org_id=%s ;",
			a.keyspace, orgID),
	}

	response, err := a.c.Session().ExecuteQuery(selectClusterIdsQuery)
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
			a.keyspace, strings.Join(clusterIdStrs, ",")),
	}

	response, err = a.c.Session().ExecuteQuery(selectClustersQuery)
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

func (a *AstraServerStore) isAppExistsInStore(name, version string) (bool, error) {
	selectClusterQuery := &pb.Query{
		Cql: fmt.Sprintf(appConfigExistanceCheckQuery,
			a.keyspace, name, version),
	}

	response, err := a.c.Session().ExecuteQuery(selectClusterQuery)
	if err != nil {
		return false, fmt.Errorf("failed to initialise db: %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) > 0 {
		return true, nil
	}

	return false, nil
}

func (a *AstraServerStore) AddOrUpdateApp(config *types.StoreAppConfig) error {
	appExists, err := a.isAppExistsInStore(config.AppName, config.Version)
	if err != nil {
		return fmt.Errorf("failed to check app config existance : %w", err)
	}

	var query *pb.Query
	if appExists {
		query = &pb.Query{
			Cql: fmt.Sprintf(updateAppConfigQuery,
				a.keyspace, config.ChartName, config.RepoName, config.RepoURL, config.Namespace, config.CreateNamespace, config.PrivilegedNamespace, config.LaunchURL, config.LaunchRedirectURL, config.Category, config.Icon, config.Description, config.LaunchUIValues, config.OverrideValues, time.Now().Format("2006-01-02 15:04:05"), config.AppName, config.Version),
		}
	} else {
		query = &pb.Query{
			Cql: fmt.Sprintf(createAppConfigQuery,
				a.keyspace, config.AppName, config.ChartName, config.RepoName, config.ReleaseName, config.RepoURL, config.Namespace, config.Version, config.CreateNamespace, config.PrivilegedNamespace, config.LaunchURL, config.LaunchRedirectURL, config.Category, config.Icon, config.Description, config.LaunchUIValues, config.OverrideValues, time.Now().Format("2006-01-02 15:04:05"), uuid.New().String()),
		}
	}

	_, err = a.c.Session().ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to insert/update the app config into the app_config table : %w", err)
	}

	return nil
}

func (a *AstraServerStore) DeleteAppInStore(name, version string) error {

	deleteQuery := &pb.Query{
		Cql: fmt.Sprintf(
			deleteAppConfigQuery,
			a.keyspace, name, version),
	}

	_, err := a.c.Session().ExecuteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to initialise db: %w", err)
	}

	return nil
}

func (a *AstraServerStore) GetAppFromStore(name, version string) (*types.AppConfig, error) {

	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(getAppConfigQuery,
			a.keyspace, name, version),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise db: %w", err)
	}

	result := response.GetResultSet()

	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("app: %s not found", name)
	}

	config, err := toAppConfig(result.Rows[0])
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (a *AstraServerStore) GetAppsFromStore() (*[]types.AppConfig, error) {

	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(getAllAppConfigsQuery,
			a.keyspace),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise db: %w", err)
	}

	result := response.GetResultSet()

	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("app config's not found")
	}

	var appConfigs []types.AppConfig
	for _, row := range result.Rows {
		config, err := toAppConfig(row)
		if err != nil {
			return nil, err
		}
		appConfigs = append(appConfigs, *config)
	}

	return &appConfigs, nil
}

func (a *AstraServerStore) GetStoreAppValues(name, version string) (*types.AppConfig, error) {

	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(getAppConfigQuery,
			a.keyspace, name, version),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app store values: %w", err)
	}

	result := response.GetResultSet()

	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("app: %s not found", name)
	}

	config, err := toAppConfig(result.Rows[0])
	if err != nil {
		return nil, err
	}

	return config, nil
}

func toAppConfig(row *pb.Row) (*types.AppConfig, error) {

	cqlAppName, err := client.ToString(row.Values[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get app name: %w", err)
	}
	cqlChartName, err := client.ToString(row.Values[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get chart name: %w", err)
	}
	cqlRepoName, err := client.ToString(row.Values[2])
	if err != nil {
		return nil, fmt.Errorf("failed to get repo name: %w", err)
	}
	cqlRepoURL, err := client.ToString(row.Values[3])
	if err != nil {
		return nil, fmt.Errorf("failed to get repo url: %w", err)
	}
	cqlNamespace, err := client.ToString(row.Values[4])
	if err != nil {
		return nil, fmt.Errorf("failed to get Namespace: %w", err)
	}
	cqlVersion, err := client.ToString(row.Values[5])
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	cqlCreateNamespace, err := client.ToBoolean(row.Values[6])
	if err != nil {
		return nil, fmt.Errorf("failed to get Create Namespace: %w", err)
	}
	cqlPrivilegedNamespace, err := client.ToBoolean(row.Values[7])
	if err != nil {
		return nil, fmt.Errorf("failed to get Privileged Namespace: %w", err)
	}
	cqlLaunchUiUrl, err := client.ToString(row.Values[8])
	if err != nil {
		return nil, fmt.Errorf("failed to get launch ui url: %w", err)
	}
	cqlLaunchUiRedirectUrl, err := client.ToString(row.Values[9])
	if err != nil {
		return nil, fmt.Errorf("failed to get launch ui redirect url: %w", err)
	}
	cqlCategory, err := client.ToString(row.Values[10])
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	cqlIcon, err := client.ToString(row.Values[11])
	if err != nil {
		return nil, fmt.Errorf("failed to get icon: %w", err)
	}
	cqlDescription, err := client.ToString(row.Values[12])
	if err != nil {
		return nil, fmt.Errorf("failed to get launch ui redirect url: %w", err)
	}
	cqlLaunchUiValues, err := client.ToString(row.Values[13])
	if err != nil {
		return nil, fmt.Errorf("failed to get launch ui values: %w", err)
	}
	cqlOverrideValues, err := client.ToString(row.Values[14])
	if err != nil {
		return nil, fmt.Errorf("failed to get override values: %w", err)
	}
	cqlReleaseNameValues, err := client.ToString(row.Values[15])
	if err != nil {
		return nil, fmt.Errorf("failed to get override values: %w", err)
	}

	config := &types.AppConfig{
		Name:                cqlAppName,
		ChartName:           cqlChartName,
		RepoName:            cqlRepoName,
		RepoURL:             cqlRepoURL,
		Namespace:           cqlNamespace,
		Version:             cqlVersion,
		CreateNamespace:     cqlCreateNamespace,
		PrivilegedNamespace: cqlPrivilegedNamespace,
		LaunchUIURL:         cqlLaunchUiUrl,
		LaunchUIRedirectURL: cqlLaunchUiRedirectUrl,
		Category:            cqlCategory,
		Icon:                cqlIcon,
		Description:         cqlDescription,
		LaunchUIValues:      cqlLaunchUiValues,
		OverrideValues:      cqlOverrideValues,
		ReleaseName:         cqlReleaseNameValues,
	}
	return config, nil
}
