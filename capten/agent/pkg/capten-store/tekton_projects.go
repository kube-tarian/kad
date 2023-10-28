package captenstore

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/pkg/errors"
)

const (
	getTektonProjectsQuery      = "SELECT id, git_project_id, status, last_update_time FROM %s.TektonProjects;"
	getTektonProjectsForIDQuery = "SELECT id, git_project_id, status, last_update_time FROM %s.TektonProjects WHERE id=%s;"
	insertTektonProjectQuery    = "INSERT INTO %s.TektonProjects(id, git_project_id, status, last_update_time, workflow_id, workflow_status) VALUES (?,?,?,?,?,?);"
	updateTektonProjectQuery    = "UPDATE %s.TektonProjects SET status='%s', last_update_time='%s', workflow_id='%s', workflow_status='%s' WHERE id=%s;"
	deleteTektonProjectQuery    = "DELETE FROM %s.TektonProjects WHERE id=%s;"
)

func (a *Store) UpsertTektonProject(project *model.TektonProject) error {
	project.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertTektonProjectQuery, a.keyspace), project.Id, project.GitProjectId, project.Status, project.LastUpdateTime)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		batch.Query(fmt.Sprintf(updateTektonProjectQuery, a.keyspace, project.Status, project.LastUpdateTime, project.WorkflowId, project.WorkflowStatus, project.Id))
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) DeleteTektonProject(id string) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(deleteTektonProjectQuery, a.keyspace, id))
	err := a.client.Session().ExecuteBatch(batch)
	return err
}

func (a *Store) GetTektonProjectForID(id string) (*model.TektonProject, error) {
	query := fmt.Sprintf(getTektonProjectsForIDQuery, a.keyspace, id)
	projects, err := a.executeTektonProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) != 1 {
		return nil, fmt.Errorf(objectNotFoundErrorMessage)
	}
	return projects[0], nil
}

func (a *Store) GetTektonProjects() ([]*model.TektonProject, error) {
	return a.updateTektonProjects()
}

// tektonProjects, err := a.GetGitProjectsByLabels([]string{"tekton"})
func (a *Store) updateTektonProjects() ([]*model.TektonProject, error) {
	allTektonProjects, err := a.GetGitProjectsByLabels([]string{"tekton"})
	if err != nil {
		a.log.Errorf("failed to fetch all tekton projects, :%v", err)
		return nil, err
	}

	query := fmt.Sprintf(getTektonProjectsQuery, a.keyspace)
	regTektonProjects, err := a.executeTektonProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	regTektonProjectId := make(map[string]*model.TektonProject)
	for _, tekPro := range regTektonProjects {
		regTektonProjectId[tekPro.Id] = tekPro
	}

	ret := make([]*model.TektonProject, 0)
	for _, allTekProject := range allTektonProjects {
		project := &model.TektonProject{Id: allTekProject.Id, GitProjectId: allTekProject.Id,
			GitProjectUrl: allTekProject.ProjectUrl}
		if _, ok := regTektonProjectId[allTekProject.Id]; !ok {
			project.Status = "available"
			project.WorkflowId = "NA"
			project.Status = "NA"
			if err := a.UpsertTektonProject(project); err != nil {
				return nil, err
			}
		} else {
			project.Status = regTektonProjectId[allTekProject.Id].Status
		}
		ret = append(ret, project)
	}

	return ret, nil
}

func (a *Store) executeTektonProjectsSelectQuery(query string) ([]*model.TektonProject, error) {
	selectAllQuery := a.client.Session().Query(query)
	iter := selectAllQuery.Iter()
	project := model.TektonProject{}

	ret := make([]*model.TektonProject, 0)
	for iter.Scan(
		&project.Id, &project.GitProjectId, &project.Status, &project.LastUpdateTime, &project.WorkflowId, &project.WorkflowStatus) {
		gitProject, err := a.GetGitProjectForID(project.Id)
		if err != nil {
			a.log.Errorf("tekton project %s not exist in git projects", project.Id)
			continue
		}

		a := &model.TektonProject{
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
