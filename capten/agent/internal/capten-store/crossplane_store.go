package captenstore

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

const (
	getAllCrossplaneProvidersQuery         = "SELECT id, cloud_type, provider_name, cloud_provider_id, status FROM %s.CrossplaneProviders;"
	insertCrossplaneProviderQuery          = "INSERT INTO %s.CrossplaneProviders(id, cloud_type, provider_name, cloud_provider_id, status) VALUES (?, ?, ?, ?, ?);"
	deleteCrossplaneProviderByIDQuery      = "DELETE FROM %s.CrossplaneProviders WHERE id=%s;"
	updateCrossplaneProviderQuery          = "UPDATE %s.CrossplaneProviders SET cloud_type=?, provider_name=?, cloud_provider_id=?, status=? WHERE id=?;"
	selectGetCrossplaneProviderByCloudType = "SELECT id, cloud_type, provider_name, cloud_provider_id, status FROM %s.CrossplaneProviders WHERE cloud_type='%s' ALLOW FILTERING;"
	selectGetCrossplaneProviderById        = "SELECT id, cloud_type, provider_name, cloud_provider_id, status FROM %s.CrossplaneProviders  WHERE id='%s';"
)

func (a *Store) InsertCrossplaneProvider(provider *model.CrossplaneProvider) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertCrossplaneProviderQuery, a.keyspace), provider.Id, provider.CloudType, provider.ProviderName, provider.CloudProviderId, provider.Status)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		return errors.Wrap(err, "failed to insert Crossplane provider")
	}
	return err
}

func (a *Store) DeleteCrossplaneProviderById(id string) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(deleteCrossplaneProviderByIDQuery, a.keyspace, id))
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		return errors.Wrap(err, "failed to delete Crossplane provider")
	}
	return nil
}

func (a *Store) GetCrossplaneProviders() ([]*captenpluginspb.CrossplaneProvider, error) {
	query := fmt.Sprintf(getAllCrossplaneProvidersQuery, a.keyspace)
	providers, err := a.executeCrossplaneProvidersSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf(objectNotFoundErrorMessage)
	}
	return providers, nil
}

func (a *Store) executeCrossplaneProvidersSelectQuery(query string) ([]*captenpluginspb.CrossplaneProvider, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	var provider captenpluginspb.CrossplaneProvider
	providers := make([]*captenpluginspb.CrossplaneProvider, 0)

	for iter.Scan(&provider.Id, &provider.CloudType, &provider.ProviderName, &provider.CloudProviderId, &provider.Status) {
		tmpProvider := &captenpluginspb.CrossplaneProvider{
			Id:              provider.Id,
			CloudType:       provider.CloudType,
			ProviderName:    provider.ProviderName,
			CloudProviderId: provider.CloudProviderId,
			Status:          provider.Status,
		}
		providers = append(providers, tmpProvider)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.Wrap(err, "error occured while iterating through results")
	}

	return providers, nil
}

func (a *Store) UpdateCrossplaneProvider(provider *model.CrossplaneProvider) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	query := fmt.Sprintf(updateCrossplaneProviderQuery, a.keyspace)
	batch.Query(query, provider.CloudType, provider.ProviderName, provider.CloudProviderId, provider.Status, provider.Id)
	err := a.client.Session().ExecuteBatch(batch)
	return err
}

func (a *Store) GetCrossplanProviderByCloudType(cloudType string) (*captenpluginspb.CrossplaneProvider, error) {
	query := fmt.Sprintf(selectGetCrossplaneProviderByCloudType, a.keyspace, cloudType)

	providers, err := a.executeCrossplaneProvidersSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(providers) <= 0 {
		return nil, nil
	}
	return providers[0], nil
}

func (a *Store) GetCrossplanProviderById(id string) (*captenpluginspb.CrossplaneProvider, error) {
	query := fmt.Sprintf(selectGetCrossplaneProviderById, a.keyspace, id)

	providers, err := a.executeCrossplaneProvidersSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(providers) <= 0 {
		return nil, nil
	}
	return providers[0], nil
}
