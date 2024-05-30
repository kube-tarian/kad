package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
)

func (a *Store) AddGitProject(config *captenpluginspb.GitProject) error {
	project := GitProject{
		ID:             uuid.MustParse(config.Id),
		ProjectURL:     config.ProjectUrl,
		Labels:         config.Labels,
		LastUpdateTime: time.Now(),
	}
	return a.dbClient.Create(&project)
}

func (a *Store) UpsertGitProject(config *captenpluginspb.GitProject) error {
	if config.Id == "" {
		project := GitProject{
			ID:             uuid.New(),
			ProjectURL:     config.ProjectUrl,
			Labels:         config.Labels,
			LastUpdateTime: time.Now(),
		}
		return a.dbClient.Create(&project)
	}

	project := GitProject{
		ProjectURL:     config.ProjectUrl,
		Labels:         config.Labels,
		LastUpdateTime: time.Now()}
	return a.dbClient.Update(project, GitProject{ID: uuid.MustParse(config.Id)})
}

func (a *Store) DeleteGitProjectById(id string) error {
	err := a.dbClient.Delete(GitProject{}, GitProject{ID: uuid.MustParse(id)})
	if err != nil {
		err = prepareError(err, id, "Delete")
	}
	return err
}

func (a *Store) GetGitProjectForID(id string) (*captenpluginspb.GitProject, error) {
	project := GitProject{}
	err := a.dbClient.FindFirst(&project, GitProject{ID: uuid.MustParse(id)})
	if err != nil {
		return nil, err
	}

	gitProject := &captenpluginspb.GitProject{
		Id:             project.ID.String(),
		ProjectUrl:     project.ProjectURL,
		Labels:         project.Labels,
		LastUpdateTime: project.LastUpdateTime.String(),
	}
	return gitProject, err
}

func (a *Store) GetGitProjects() ([]*captenpluginspb.GitProject, error) {
	projects := []GitProject{}
	err := a.dbClient.Find(&projects, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %v", err.Error())
	}

	gitProjects := make([]*captenpluginspb.GitProject, 0)
	for _, project := range projects {
		gitProjects = append(gitProjects, &captenpluginspb.GitProject{
			Id:             project.ID.String(),
			ProjectUrl:     project.ProjectURL,
			Labels:         project.Labels,
			LastUpdateTime: project.LastUpdateTime.String(),
		})
	}
	return gitProjects, err
}

func (a *Store) GetGitProjectsByLabels(searchLabels []string) ([]*captenpluginspb.GitProject, error) {
	projects := []GitProject{}
	err := a.dbClient.Find(&projects, "labels @> ?", fmt.Sprintf("{%s}", searchLabels[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %v", err.Error())
	}

	gitProjects := make([]*captenpluginspb.GitProject, 0)
	for _, project := range projects {
		gitProjects = append(gitProjects, &captenpluginspb.GitProject{
			Id:             project.ID.String(),
			ProjectUrl:     project.ProjectURL,
			Labels:         project.Labels,
			LastUpdateTime: project.LastUpdateTime.String(),
		})
	}

	return gitProjects, err
}
