package captenstore

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/pkg/errors"
)

const (
	getArgocdProjectsQuery      = "SELECT id, git_project_id, status, last_update_time FROM %s.ArgocdProjects;"
	getArgocdProjectsForIdQuery = "SELECT id, git_project_id, status, last_update_time FROM %s.ArgocdProjects WHERE id='%s';"
	insertArgocdProjectQuery    = "INSERT INTO %s.ArgocdProjects(id, git_project_id, status, last_update_time) VALUES (?,?,?);"
	updateArgocdProjectQuery    = "UPDATE %s.ArgocdProjects SET status='%s', lastUpdateTime='%s' WHERE id='%s' and git_project_id='%s';"
	deleteArgocdProjectQuery    = "DELETE FROM %s.ArgocdProjects WHERE id='%s';"
)

func (a *Store) UpsertArgoCDProject(project *model.ArgoCDProject) error {
	project.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertArgocdProjectQuery, a.keyspace), project.Id, project.GitProjectId, project.Status, project.LastUpdateTime)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		batch.Query(fmt.Sprintf(updateArgocdProjectQuery, a.keyspace, project.Status, project.LastUpdateTime, project.Id, project.GitProjectId))
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) GetArgoCDProjectForID(id string) (*model.ArgoCDProject, error) {
	query := fmt.Sprintf(getArgocdProjectsForIdQuery, a.keyspace, id)
	projects, err := a.executeArgoCDProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) != 1 {
		return nil, fmt.Errorf("project not found")
	}
	return projects[0], nil
}

func (a *Store) GetArgoCDProjects() ([]*model.ArgoCDProject, error) {
	return a.updateArgoCDProjects()
}

func (a *Store) DeleteArgoCDProjectsData(id string) error {
	deleteQuery := a.client.Session().Query(fmt.Sprintf(deleteArgocdProjectQuery,
		a.keyspace, id))
	err := deleteQuery.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *Store) executeArgoCDProjectsSelectQuery(query string) ([]*model.ArgoCDProject, error) {
	selectAllQuery := a.client.Session().Query(query)
	iter := selectAllQuery.Iter()
	project := model.ArgoCDProject{}

	ret := make([]*model.ArgoCDProject, 0)
	for iter.Scan(
		&project.Id, project.GitProjectId, project.Status, &project.LastUpdateTime) {
		gitProject, err := a.GetGitProjectForID(project.GitProjectId)
		if err != nil {
			a.log.Errorf("argocd project %s not exist in git projects", project.GitProjectId)
			continue
		}

		a := &model.ArgoCDProject{
			Id:             project.Id,
			GitProjectId:   gitProject.Id,
			GitProjectUrl:  gitProject.ProjectUrl,
			Status:         project.Status,
			LastUpdateTime: project.LastUpdateTime,
		}
		ret = append(ret, a)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}
	return ret, nil
}

func (a *Store) updateArgoCDProjects() ([]*model.ArgoCDProject, error) {
	gitProjects, err := a.GetGitProjectsByLabels([]string{"argocd"})
	if err != nil {
		a.log.Errorf("failed to fetch all argocd projects, :%v", err)
		return nil, err
	}

	query := fmt.Sprintf(getArgocdProjectsQuery, a.keyspace)
	regArgoCDProjects, err := a.executeArgoCDProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	argoCDProjects := make(map[string]*model.ArgoCDProject)
	for _, tekPro := range regArgoCDProjects {
		argoCDProjects[tekPro.Id] = tekPro
	}

	ret := make([]*model.ArgoCDProject, 0)
	for _, gitProject := range gitProjects {
		project := &model.ArgoCDProject{Id: gitProject.Id, GitProjectId: gitProject.Id,
			GitProjectUrl: gitProject.ProjectUrl}
		if _, ok := argoCDProjects[gitProject.Id]; !ok {
			project.Status = "available"
			if err := a.UpsertArgoCDProject(project); err != nil {
				return nil, err
			}
		} else {
			project.Status = argoCDProjects[gitProject.Id].Status
		}
		ret = append(ret, project)
	}
	return ret, nil
}
