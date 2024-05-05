package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
)

func (a *Store) UpsertContainerRegistry(config *captenpluginspb.ContainerRegistry) error {
	if config.Id == "" {
		registry := ContainerRegistry{
			ID:             uuid.New(),
			RegistryURL:    config.RegistryUrl,
			RegistryType:   config.RegistryType,
			Labels:         config.Labels,
			LastUpdateTime: time.Now(),
		}
		return a.dbClient.Create(&registry)
	}

	registry := ContainerRegistry{RegistryURL: config.RegistryUrl,
		ID:             uuid.MustParse(config.Id),
		RegistryType:   config.RegistryType,
		Labels:         config.Labels,
		LastUpdateTime: time.Now()}
	return a.dbClient.Update(&registry, ContainerRegistry{ID: registry.ID})
}

func (a *Store) GetContainerRegistryForID(id string) (*captenpluginspb.ContainerRegistry, error) {
	registry := ContainerRegistry{}
	err := a.dbClient.Find(&registry, ContainerRegistry{ID: uuid.MustParse(id)})
	if err != nil {
		return nil, err
	}

	result := &captenpluginspb.ContainerRegistry{
		Id:             registry.ID.String(),
		RegistryUrl:    registry.RegistryURL,
		RegistryType:   registry.RegistryType,
		Labels:         registry.Labels,
		LastUpdateTime: registry.LastUpdateTime.String(),
	}
	return result, err
}

func (a *Store) GetContainerRegistries() ([]*captenpluginspb.ContainerRegistry, error) {
	registries := []ContainerRegistry{}
	err := a.dbClient.Find(&registries, nil)
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		return nil, fmt.Errorf("failed to fetch registries: %v", err.Error())
	}

	result := make([]*captenpluginspb.ContainerRegistry, 0)
	for _, registry := range registries {
		result = append(result, &captenpluginspb.ContainerRegistry{
			Id:             registry.ID.String(),
			RegistryUrl:    registry.RegistryURL,
			RegistryType:   registry.RegistryType,
			Labels:         registry.Labels,
			LastUpdateTime: registry.LastUpdateTime.String(),
		})
	}
	return result, err
}

func (a *Store) GetContainerRegistriesByLabels(searchLabels []string) ([]*captenpluginspb.ContainerRegistry, error) {
	registries := []ContainerRegistry{}
	err := a.dbClient.Find(&registries, "labels @> ?", fmt.Sprintf("{%s}", searchLabels[0]))
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		return nil, fmt.Errorf("failed to fetch registries: %v", err.Error())
	}

	result := make([]*captenpluginspb.ContainerRegistry, 0)
	for _, registry := range registries {
		result = append(result, &captenpluginspb.ContainerRegistry{
			Id:             registry.ID.String(),
			RegistryUrl:    registry.RegistryURL,
			RegistryType:   registry.RegistryType,
			Labels:         registry.Labels,
			LastUpdateTime: registry.LastUpdateTime.String(),
		})
	}
	return result, err
}

func (a *Store) DeleteContainerRegistryById(id string) error {
	err := a.dbClient.Delete(ContainerRegistry{}, ContainerRegistry{ID: uuid.MustParse(id)})
	if err != nil {
		err = prepareError(err, id, "Delete")
	}
	return err
}
