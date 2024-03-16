package postgres

import (
	"time"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
)

func (handler *Postgres) UpsertGitProject(config *captenpluginspb.GitProject) error {

	if config.Id == "" {
		gp := GitProjects{
			ProjectURL:     config.ProjectUrl,
			Labels:         config.Labels,
			LastUpdateTime: time.Now(),
		}
		return handler.db.Create(&gp).Error
	}
	return handler.db.Where("id = ", config.Id).Updates(GitProjects{ProjectURL: config.ProjectUrl, Labels: config.Labels, LastUpdateTime: time.Now()}).Error
}

func (handler *Postgres) DeleteGitProjectById(id string) error {

	err := handler.db.Where("id = ", id).Delete(&GitProjects{}).Error
	return err
}

func (handler *Postgres) GetGitProjectForID(id string) (*captenpluginspb.GitProject, error) {

	gp := GitProjects{}

	err := handler.db.Select("*").Where("id = ", id).Scan(&gp).Error

	project := &captenpluginspb.GitProject{
		Id:             gp.ID.String(),
		ProjectUrl:     gp.ProjectURL,
		Labels:         gp.Labels,
		LastUpdateTime: gp.LastUpdateTime.String(),
		UsedPlugins:    gp.UsedPlugins,
	}

	return project, err
}

func (handler *Postgres) GetGitProjects() ([]*captenpluginspb.GitProject, error) {

	gp := []GitProjects{}

	err := handler.db.Select("*").Scan(&gp).Error

	result := make([]*captenpluginspb.GitProject, 0)
	for _, v := range gp {
		result = append(result, &captenpluginspb.GitProject{
			Id:             v.ID.String(),
			ProjectUrl:     v.ProjectURL,
			Labels:         v.Labels,
			LastUpdateTime: v.LastUpdateTime.String(),
			UsedPlugins:    v.UsedPlugins,
		})
	}

	return result, err
}

func (handler *Postgres) GetGitProjectsByLabels(searchLabels []string) ([]*captenpluginspb.GitProject, error) {
	gp := []GitProjects{}

	err := handler.db.Select("labels @> ?", searchLabels).Scan(&gp).Error

	result := make([]*captenpluginspb.GitProject, 0)
	for _, v := range gp {
		result = append(result, &captenpluginspb.GitProject{
			Id:             v.ID.String(),
			ProjectUrl:     v.ProjectURL,
			Labels:         v.Labels,
			LastUpdateTime: v.LastUpdateTime.String(),
			UsedPlugins:    v.UsedPlugins,
		})
	}

	return result, err
}
