package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
	"github.com/kube-tarian/kad/capten/model"
)

func (a *Store) UpsertCrossplaneProvider(crossplaneProvider *model.CrossplaneProvider) error {
	if crossplaneProvider.Id == "" {
		provider := CrossplaneProvider{
			ID:              uuid.New(),
			ProviderName:    crossplaneProvider.ProviderName,
			CloudProviderID: crossplaneProvider.CloudProviderId,
			CloudType:       crossplaneProvider.CloudType,
			Status:          crossplaneProvider.Status,
			LastUpdateTime:  time.Now(),
		}
		return a.dbClient.Create(&provider)
	}

	provider := CrossplaneProvider{CloudType: crossplaneProvider.CloudType,
		CloudProviderID: crossplaneProvider.CloudProviderId,
		Status:          crossplaneProvider.Status,
		LastUpdateTime:  time.Now()}

	return a.dbClient.Update(provider, CrossplaneProvider{ID: uuid.MustParse(crossplaneProvider.Id)})
}

func (a *Store) DeleteCrossplaneProviderById(id string) error {
	err := a.dbClient.Delete(CrossplaneProvider{}, CrossplaneProvider{ID: uuid.MustParse(id)})
	return err
}

func (a *Store) GetCrossplaneProviders() ([]*captenpluginspb.CrossplaneProvider, error) {
	providers := []CrossplaneProvider{}
	err := a.dbClient.Find(&providers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch providers: %v", err.Error())
	}

	crossplaneProviders := []*captenpluginspb.CrossplaneProvider{}
	for _, provider := range providers {
		crossplaneProviders = append(crossplaneProviders, &captenpluginspb.CrossplaneProvider{
			Id:              provider.ID.String(),
			CloudProviderId: provider.CloudProviderID,
			CloudType:       provider.CloudType,
			Status:          provider.Status,
		})
	}
	return crossplaneProviders, nil
}

func (a *Store) UpdateCrossplaneProvider(provider *model.CrossplaneProvider) error {
	crossplaneProvider := CrossplaneProvider{CloudType: provider.CloudType,
		CloudProviderID: provider.CloudProviderId,
		Status:          provider.Status,
		LastUpdateTime:  time.Now()}

	return a.dbClient.Update(crossplaneProvider, CrossplaneProvider{ID: uuid.MustParse(provider.Id)})
}

func (a *Store) GetCrossplanProviderByCloudType(cloudType string) (*captenpluginspb.CrossplaneProvider, error) {
	providers := []CrossplaneProvider{}
	err := a.dbClient.Find(&providers, "cloud_type = ?", cloudType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch providers: %v", err.Error())
	}

	if len(providers) == 1 {
		provider := providers[0]
		crossplaneProvider := &captenpluginspb.CrossplaneProvider{
			Id:              provider.ID.String(),
			ProviderName:    provider.ProviderName,
			CloudProviderId: provider.CloudProviderID,
			CloudType:       provider.CloudType,
			Status:          provider.Status,
		}
		return crossplaneProvider, err
	}
	return nil, gerrors.New(postgresdb.ObjectNotExist, "Crossplane provider not found")
}

func (a *Store) GetCrossplanProviderById(id string) (*captenpluginspb.CrossplaneProvider, error) {
	provider := CrossplaneProvider{}
	err := a.dbClient.FindFirst(&provider, CrossplaneProvider{ID: uuid.MustParse(id)})
	if err != nil {
		return nil, err
	}

	crossplaneProvider := &captenpluginspb.CrossplaneProvider{
		Id:              provider.ID.String(),
		ProviderName:    provider.ProviderName,
		CloudProviderId: provider.CloudProviderID,
		CloudType:       provider.CloudType,
		Status:          provider.Status,
	}
	return crossplaneProvider, err
}
