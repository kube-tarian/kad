package captenstore

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/pkg/errors"
)

const (
	insertCloudProvider               = "INSERT INTO %s.CloudProviders(id, cloud_type, labels, last_update_time) VALUES (?,?,?,?)"
	insertCloudProviderId             = "INSERT INTO %s.CloudProviders(id) VALUES (?)"
	updateCloudProviderById           = "UPDATE %s.CloudProviders SET %s WHERE id=?"
	deleteCloudProviderById           = "DELETE FROM %s.CloudProviders WHERE id= ?"
	selectAllCloudProviders           = "SELECT id, cloud_type, labels, last_update_time FROM %s.CloudProviders"
	selectAllCloudProvidersByLabels   = "SELECT id, cloud_type, labels, last_update_time FROM %s.CloudProviders WHERE %s"
	selectGetCloudProviderById        = "SELECT id, cloud_type, labels, last_update_time FROM %s.CloudProviders WHERE id=%s;"
	selectGetCloudProviderByCloudType = "SELECT id, cloud_type, labels, last_update_time FROM %s.CloudProviders WHERE cloud_type=%s;"
)

func (a *Store) UpsertCloudProvider(config *captenpluginspb.CloudProvider) error {
	config.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertCloudProvider, a.keyspace), config.Id, config.CloudType, config.Labels, config.LastUpdateTime)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		updatePlaceholders, values := formUpdateKvPairsForCloudProvider(config)
		if updatePlaceholders == "" {
			return err
		}
		query := fmt.Sprintf(updateCloudProviderById, a.keyspace, updatePlaceholders)
		args := append(values, config.Id)
		batch = a.client.Session().NewBatch(gocql.LoggedBatch)
		batch.Query(query, args...)
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) DeleteCloudProviderById(id string) error {
	deleteAction := a.client.Session().Query(fmt.Sprintf(deleteCloudProviderById,
		a.keyspace), id)
	err := deleteAction.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *Store) GetCloudProviderForID(id string) (*captenpluginspb.CloudProvider, error) {
	query := fmt.Sprintf(selectGetCloudProviderById, a.keyspace, id)
	projects, err := a.executeCloudProvidersSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) != 1 {
		return nil, fmt.Errorf("project not found")
	}
	return projects[0], nil
}

func (a *Store) GetCloudProviders() ([]*captenpluginspb.CloudProvider, error) {
	query := fmt.Sprintf(selectAllCloudProviders, a.keyspace)
	return a.executeCloudProvidersSelectQuery(query)
}

func (a *Store) GetCloudProvidersByLabelsAndCloudType(searchLabels []string, cloudType string) ([]*captenpluginspb.CloudProvider, error) {

	whereLabelsClause := ""
	if cloudType != "" {
		whereLabelsClause += fmt.Sprintf("cloud_type = '%s'", cloudType)
	}

	if len(searchLabels) != 0 {
		if whereLabelsClause != "" {
			whereLabelsClause += " AND "
		}
		labelContains := []string{}
		for _, label := range searchLabels {
			labelContains = append(labelContains, fmt.Sprintf("labels CONTAINS '%s'", label))
		}
		whereLabelsClause += "(" + strings.Join(labelContains, " OR ") + ")"
		whereLabelsClause += " ALLOW FILTERING"
	}

	query := fmt.Sprintf(selectAllCloudProvidersByLabels, a.keyspace, whereLabelsClause)
	return a.executeCloudProvidersSelectQuery(query)

}

func (a *Store) GetCloudProviderByCloudType(cloudType string) (*captenpluginspb.CloudProvider, error) {
	query := fmt.Sprintf(selectGetCloudProviderByCloudType, a.keyspace, cloudType)
	fmt.Println("Query => ", query)
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	provider := captenpluginspb.CloudProvider{}
	var labels []string

	ret := make([]*captenpluginspb.CloudProvider, 0)
	for iter.Scan(
		&provider.Id, &provider.CloudType,
		&labels, &provider.LastUpdateTime,
	) {
		labelsTmp := make([]string, len(labels))
		copy(labelsTmp, labels)
		CloudProvider := &captenpluginspb.CloudProvider{
			Id:             provider.Id,
			CloudType:      provider.CloudType,
			Labels:         labelsTmp,
			LastUpdateTime: provider.LastUpdateTime,
		}
		ret = append(ret, CloudProvider)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	v, _ := json.Marshal(ret)
	fmt.Println("Cloud Provider => \n" + string(v))

	if len(ret) <= 0 {
		return nil, nil
	}
	return ret[0], nil
}

func (a *Store) executeCloudProvidersSelectQuery(query string) ([]*captenpluginspb.CloudProvider, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	provider := captenpluginspb.CloudProvider{}
	var labels []string

	ret := make([]*captenpluginspb.CloudProvider, 0)
	for iter.Scan(
		&provider.Id, &provider.CloudType,
		&labels, &provider.LastUpdateTime,
	) {
		labelsTmp := make([]string, len(labels))
		copy(labelsTmp, labels)
		CloudProvider := &captenpluginspb.CloudProvider{
			Id:             provider.Id,
			CloudType:      provider.CloudType,
			Labels:         labelsTmp,
			LastUpdateTime: provider.LastUpdateTime,
		}
		ret = append(ret, CloudProvider)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return ret, nil
}

func formUpdateKvPairsForCloudProvider(config *captenpluginspb.CloudProvider) (updatePlaceholders string, values []interface{}) {
	params := []string{}
	values = []interface{}{}
	if config.CloudType != "" {
		params = append(params, "cloud_type = ?")
		values = append(values, config.CloudType)
	}

	if len(config.Labels) > 0 {
		labels := []string{}
		for _, label := range config.Labels {
			labels = append(labels, fmt.Sprintf("'%s'", label))
		}
		params = append(params, "labels = ?")
		labelsStr := "{" + strings.Join(labels, ", ") + "}"
		values = append(values, labelsStr)
	}

	if config.LastUpdateTime != "" {
		params = append(params, "last_update_time = ?")
		values = append(values, config.LastUpdateTime)
	}

	if len(params) == 0 {
		return "", nil
	}

	return strings.Join(params, ", "), values
}
