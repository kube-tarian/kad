package postgres

import (
	"time"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
)

func (handler *Postgres) UpsertContainerRegistry(config *captenpluginspb.ContainerRegistry) error {
	if config.Id == "" {
		gp := ContainerRegistry{
			RegistryURL:    config.RegistryUrl,
			RegistryType:   config.RegistryType,
			Labels:         config.Labels,
			LastUpdateTime: time.Now(),
		}
		return handler.db.Create(&gp).Error
	}
	return handler.db.Where("id = ", config.Id).Updates(ContainerRegistry{RegistryURL: config.RegistryUrl, RegistryType: config.RegistryType, Labels: config.Labels, LastUpdateTime: time.Now()}).Error
}

func (handler *Postgres) DeleteContainerRegistryById(id string) error {

	err := handler.db.Where("id = ", id).Delete(&ContainerRegistry{}).Error
	return err
}

func (handler *Postgres) GetContainerRegistryForID(id string) (*captenpluginspb.ContainerRegistry, error) {

	cr := ContainerRegistry{}

	err := handler.db.Select("*").Where("id = ", id).Scan(&cr).Error

	result := &captenpluginspb.ContainerRegistry{
		Id:             cr.ID.String(),
		RegistryUrl:    cr.RegistryURL,
		RegistryType:   cr.RegistryType,
		Labels:         cr.Labels,
		LastUpdateTime: cr.LastUpdateTime.String(),
		UsedPlugins:    cr.UsedPlugins,
	}

	return result, err
}

func (handler *Postgres) GetContainerRegistries() ([]*captenpluginspb.ContainerRegistry, error) {

	crs := []ContainerRegistry{}

	err := handler.db.Select("*").Scan(&crs).Error

	result := make([]*captenpluginspb.ContainerRegistry, 0)
	for _, cr := range crs {
		result = append(result, &captenpluginspb.ContainerRegistry{
			Id:             cr.ID.String(),
			RegistryUrl:    cr.RegistryURL,
			RegistryType:   cr.RegistryType,
			Labels:         cr.Labels,
			LastUpdateTime: cr.LastUpdateTime.String(),
			UsedPlugins:    cr.UsedPlugins,
		})
	}

	return result, err
}

func (handler *Postgres) GetContainerRegistriesByLabels(searchLabels []string) ([]*captenpluginspb.ContainerRegistry, error) {

	cps := []ContainerRegistry{}

	err := handler.db.Select("labels @> ? ", searchLabels).Scan(&cps).Error

	result := make([]*captenpluginspb.ContainerRegistry, 0)
	for _, cr := range cps {
		result = append(result, &captenpluginspb.ContainerRegistry{
			Id:             cr.ID.String(),
			RegistryUrl:    cr.RegistryURL,
			RegistryType:   cr.RegistryType,
			Labels:         cr.Labels,
			LastUpdateTime: cr.LastUpdateTime.String(),
			UsedPlugins:    cr.UsedPlugins,
		})
	}

	return result, err
}
