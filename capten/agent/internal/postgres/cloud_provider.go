package postgres

import (
	"time"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
)

func (handler *Postgres) UpsertCloudProvider(config *captenpluginspb.CloudProvider) error {

	if config.Id == "" {
		gp := CloudProviders{
			CloudType:      config.CloudType,
			Labels:         config.Labels,
			LastUpdateTime: time.Now(),
		}
		return handler.db.Create(&gp).Error
	}
	return handler.db.Where("id = ", config.Id).Updates(CloudProviders{CloudType: config.CloudType, Labels: config.Labels, LastUpdateTime: time.Now()}).Error

}

func (handler *Postgres) DeleteCloudProviderById(id string) error {

	err := handler.db.Where("id = ", id).Delete(&CloudProviders{}).Error
	return err

}

func (handler *Postgres) GetCloudProviderForID(id string) (*captenpluginspb.CloudProvider, error) {

	cp := CloudProviders{}

	err := handler.db.Select("*").Where("id = ", id).Scan(&cp).Error

	result := &captenpluginspb.CloudProvider{
		Id:             cp.ID.String(),
		CloudType:      cp.CloudType,
		Labels:         cp.Labels,
		LastUpdateTime: cp.LastUpdateTime.String(),
	}

	return result, err
}

func (handler *Postgres) GetCloudProviders() ([]*captenpluginspb.CloudProvider, error) {

	cp := []CloudProviders{}

	err := handler.db.Select("*").Scan(&cp).Error

	result := make([]*captenpluginspb.CloudProvider, 0)
	for _, v := range cp {
		result = append(result, &captenpluginspb.CloudProvider{
			Id:             v.ID.String(),
			CloudType:      v.CloudType,
			Labels:         v.Labels,
			LastUpdateTime: v.LastUpdateTime.String(),
		})
	}

	return result, err
}

func (handler *Postgres) GetCloudProvidersByLabelsAndCloudType(searchLabels []string, cloudType string) ([]*captenpluginspb.CloudProvider, error) {

	cps := []CloudProviders{}

	err := handler.db.Select("cloud_type = ? AND labels @> ? ", cloudType, searchLabels).Scan(&cps).Error

	result := make([]*captenpluginspb.CloudProvider, 0)
	for _, v := range cps {
		result = append(result, &captenpluginspb.CloudProvider{
			Id:             v.ID.String(),
			CloudType:      v.CloudType,
			Labels:         v.Labels,
			LastUpdateTime: v.LastUpdateTime.String(),
		})
	}

	return result, err

}
