package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (a *Store) UpsertCrossplaneProject(crossplaneProject *model.CrossplaneProject) error {
	gitProjectUUID, err := uuid.Parse(crossplaneProject.GitProjectId)
	if err != nil {
		return err
	}

	if crossplaneProject.Id == "" {
		project := CrossplaneProject{
			ID:             1,
			GitProjectID:   gitProjectUUID,
			GitProjectURL:  crossplaneProject.GitProjectUrl,
			Status:         crossplaneProject.Status,
			LastUpdateTime: time.Now(),
		}
		return a.dbClient.Create(&project)
	}

	project := CrossplaneProject{
		GitProjectID:   gitProjectUUID,
		GitProjectURL:  crossplaneProject.GitProjectUrl,
		Status:         crossplaneProject.Status,
		LastUpdateTime: time.Now()}
	return a.dbClient.Update(project, CrossplaneProject{ID: 1})
}

func (a *Store) DeleteCrossplaneProject(id string) error {
	err := a.dbClient.Delete(CrossplaneProject{}, CrossplaneProject{ID: 1})
	if err != nil {
		err = prepareError(err, id, "Delete")
	}
	return err
}

func (a *Store) GetCrossplaneProjectForID(id string) (*model.CrossplaneProject, error) {
	project := CrossplaneProject{}
	err := a.dbClient.Find(&project, CrossplaneProject{ID: 1})
	if err != nil {
		return nil, err
	} else if project.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	crossplaneProject := &model.CrossplaneProject{
		Id:             "1",
		GitProjectId:   project.GitProjectID.String(),
		GitProjectUrl:  project.GitProjectURL,
		Status:         project.Status,
		LastUpdateTime: project.LastUpdateTime.String(),
	}
	return crossplaneProject, err
}

func (a *Store) GetCrossplaneProject() (*model.CrossplaneProject, error) {
	return a.updateCrossplaneProject()
}

func (a *Store) updateCrossplaneProject() (*model.CrossplaneProject, error) {
	allCrossplaneGitProjects, err := a.GetGitProjectsByLabels([]string{"crossplane"})
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		return nil, fmt.Errorf("failed to fetch projects: %v", err.Error())
	}

	if len(allCrossplaneGitProjects) == 0 {
		return nil, fmt.Errorf("no git project found with crossplane tag")
	}
	crosplaneGitProject := allCrossplaneGitProjects[0]
	gitProjectUUID, err := uuid.Parse(crosplaneGitProject.Id)
	if err != nil {
		return nil, err
	}

	crossplaneProject, err := a.GetCrossplaneProjectForID("0")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			project := CrossplaneProject{
				ID:             1,
				GitProjectID:   gitProjectUUID,
				GitProjectURL:  crosplaneGitProject.ProjectUrl,
				LastUpdateTime: time.Now(),
			}
			err = a.dbClient.Create(&project)
			if err != nil {
				return nil, err
			}
			return a.GetCrossplaneProjectForID("1")
		} else {
			return nil, err
		}

	}

	if crossplaneProject.GitProjectId == crosplaneGitProject.Id &&
		crossplaneProject.GitProjectUrl == crosplaneGitProject.ProjectUrl {
		return crossplaneProject, nil
	}

	project := CrossplaneProject{GitProjectID: gitProjectUUID,
		GitProjectURL:  crosplaneGitProject.ProjectUrl,
		Status:         crossplaneProject.Status,
		LastUpdateTime: time.Now()}
	err = a.dbClient.Update(&project, CrossplaneProject{ID: 1})
	if err != nil {
		return nil, err
	}

	// project already registered, return that
	return a.GetCrossplaneProjectForID("1")
}
