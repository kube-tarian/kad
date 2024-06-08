package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
	"github.com/kube-tarian/kad/capten/model"
)

func (a *Store) UpsertTektonProject(tektonProject *model.TektonProject) error {
	gitProjectUUID, err := uuid.Parse(tektonProject.GitProjectId)
	if err != nil {
		return err
	}

	if tektonProject.Id == "" {
		project := TektonProject{
			ID:             1,
			GitProjectID:   gitProjectUUID,
			GitProjectURL:  tektonProject.GitProjectUrl,
			Status:         tektonProject.Status,
			LastUpdateTime: time.Now(),
		}
		return a.dbClient.Create(&project)
	}

	project := TektonProject{GitProjectID: gitProjectUUID,
		GitProjectURL:  tektonProject.GitProjectUrl,
		Status:         tektonProject.Status,
		LastUpdateTime: time.Now()}
	return a.dbClient.Update(project, TektonProject{ID: 1})
}

func (a *Store) DeleteTektonProject(id string) error {
	err := a.dbClient.Delete(TektonProject{}, TektonProject{ID: 1})
	if err != nil {
		err = prepareError(err, id, "Delete")
	}
	return err
}

func (a *Store) GetTektonProjectForID(id string) (*model.TektonProject, error) {
	project := TektonProject{}
	err := a.dbClient.FindFirst(&project, TektonProject{ID: 1})
	if err != nil {
		return nil, err
	}

	tektonProject := &model.TektonProject{
		Id:             "1",
		GitProjectId:   project.GitProjectID.String(),
		GitProjectUrl:  project.GitProjectURL,
		Status:         project.Status,
		LastUpdateTime: project.LastUpdateTime.String(),
	}
	return tektonProject, err
}

func (a *Store) GetTektonProject() (*model.TektonProject, error) {
	return a.updateTektonProject()
}

func (a *Store) updateTektonProject() (*model.TektonProject, error) {
	allTektonGitProjects, err := a.GetGitProjectsByLabels([]string{"tekton"})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %v", err.Error())
	}

	if len(allTektonGitProjects) == 0 {
		return nil, fmt.Errorf("no git project found with tekton tag")
	}
	tektonGitProject := allTektonGitProjects[0]
	gitProjectUUID, err := uuid.Parse(tektonGitProject.Id)
	if err != nil {
		return nil, err
	}

	tektonProject, err := a.GetTektonProjectForID("1")
	if err != nil {
		if gerrors.GetErrorType(err) == postgresdb.ObjectNotExist {
			project := TektonProject{
				ID:             1,
				GitProjectID:   gitProjectUUID,
				GitProjectURL:  tektonGitProject.ProjectUrl,
				Status:         string(model.TektonProjectAvailable),
				LastUpdateTime: time.Now(),
			}
			err = a.dbClient.Create(&project)
			if err != nil {
				return nil, err
			}
			return a.GetTektonProjectForID("1")
		} else {
			return nil, err
		}
	}

	if tektonProject.GitProjectId == tektonGitProject.Id &&
		tektonProject.GitProjectUrl == tektonGitProject.ProjectUrl {
		return tektonProject, nil
	}

	project := TektonProject{GitProjectID: gitProjectUUID,
		GitProjectURL:  tektonGitProject.ProjectUrl,
		Status:         string(model.TektonProjectAvailable),
		LastUpdateTime: time.Now()}
	err = a.dbClient.Update(&project, TektonProject{ID: 1})
	if err != nil {
		return nil, err
	}

	// project already registered, return that
	return a.GetTektonProjectForID("1")
}
