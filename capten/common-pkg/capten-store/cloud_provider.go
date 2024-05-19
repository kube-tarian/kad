package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
	"gorm.io/gorm"
)

func (a *Store) UpsertCloudProvider(config *captenpluginspb.CloudProvider) error {
	if config.Id == "" {
		provider := CloudProvider{
			ID:             uuid.New(),
			CloudType:      config.CloudType,
			Labels:         config.Labels,
			LastUpdateTime: time.Now(),
		}
		return a.dbClient.Create(&provider)
	}

	provider := CloudProvider{
		ID:             uuid.MustParse(config.Id),
		CloudType:      config.CloudType,
		Labels:         config.Labels,
		LastUpdateTime: time.Now()}

	return a.dbClient.Update(provider, CloudProvider{ID: provider.ID})
}

func (a *Store) GetCloudProviderForID(id string) (*captenpluginspb.CloudProvider, error) {
	provider := CloudProvider{}
	err := a.dbClient.Find(&provider, CloudProvider{ID: uuid.MustParse(id)})
	if err != nil {
		return nil, err
	}

	result := &captenpluginspb.CloudProvider{
		Id:             provider.ID.String(),
		CloudType:      provider.CloudType,
		Labels:         provider.Labels,
		LastUpdateTime: provider.LastUpdateTime.String(),
	}

	return result, err
}

func (a *Store) GetCloudProviders() ([]*captenpluginspb.CloudProvider, error) {
	providers := []CloudProvider{}
	err := a.dbClient.Find(&providers, nil)
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		return nil, fmt.Errorf("failed to fetch providers: %v", err.Error())
	}

	cloudProviders := make([]*captenpluginspb.CloudProvider, 0)
	for _, provider := range providers {
		cloudProviders = append(cloudProviders, &captenpluginspb.CloudProvider{
			Id:             provider.ID.String(),
			CloudType:      provider.CloudType,
			Labels:         provider.Labels,
			LastUpdateTime: provider.LastUpdateTime.String(),
		})
	}
	return cloudProviders, err
}

func (a *Store) GetCloudProvidersByLabelsAndCloudType(searchLabels []string, cloudType string) ([]*captenpluginspb.CloudProvider, error) {
	providers := []CloudProvider{}
	err := a.dbClient.Session().Where("cloud_type = ?", cloudType).Where("labels @> ?", fmt.Sprintf("{%s}", searchLabels[0])).Find(&providers).Error
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		if gorm.ErrRecordNotFound != err {
			return nil, fmt.Errorf("failed to fetch providers: %v", err.Error())
		}
		err = nil
	}

	cloudProviders := make([]*captenpluginspb.CloudProvider, 0)
	for _, provider := range providers {
		cloudProviders = append(cloudProviders, &captenpluginspb.CloudProvider{
			Id:             provider.ID.String(),
			CloudType:      provider.CloudType,
			Labels:         provider.Labels,
			LastUpdateTime: provider.LastUpdateTime.String(),
		})
	}
	return cloudProviders, err
}

func (a *Store) GetCloudProvidersByLabels(searchLabels []string) ([]*captenpluginspb.CloudProvider, error) {
	providers := []CloudProvider{}
	err := a.dbClient.Find(&providers, "labels @> ?", fmt.Sprintf("{%s}", searchLabels[0]))
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		return nil, fmt.Errorf("failed to fetch providers: %v", err.Error())
	}

	cloudProviders := make([]*captenpluginspb.CloudProvider, 0)
	for _, provider := range providers {
		cloudProviders = append(cloudProviders, &captenpluginspb.CloudProvider{
			Id:             provider.ID.String(),
			CloudType:      provider.CloudType,
			Labels:         provider.Labels,
			LastUpdateTime: provider.LastUpdateTime.String(),
		})
	}
	return cloudProviders, err
}

func (a *Store) DeleteCloudProviderById(id string) error {
	err := a.dbClient.Delete(CloudProvider{}, CloudProvider{ID: uuid.MustParse(id)})
	if err != nil {
		err = prepareError(err, id, "Delete")
	}
	return err
}
