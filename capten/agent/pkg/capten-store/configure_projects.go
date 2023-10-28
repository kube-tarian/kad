package captenstore

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/pkg/errors"
)

const (
	TektonTableName = "TektonProjects"
)

const (
	getConfigProjectsQuery      = "SELECT id, git_project_id, status, last_update_time, workflow_id FROM %s.%s;"
	getConfigProjectsForIDQuery = "SELECT id, git_project_id, status, last_update_time, workflow_id FROM %s.%s WHERE id=%s;"
	insertConfigProjectQuery    = "INSERT INTO %s.%s(id, git_project_id, status, last_update_time, workflow_id) VALUES (?,?,?,?,?);"
	updateConfigProjectQuery    = "UPDATE %s.%s SET status='%s', last_update_time='%s', workflow_id='%s' WHERE id=%s;"
	deleteConfigProjectQuery    = "DELETE FROM %s.%s WHERE id=%s;"
)

func (a *Store) UpsertConfigProject(project *model.ConfigureProject, tablename string) error {
	project.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertConfigProjectQuery, a.keyspace, tablename), project.Id, project.GitProjectId, project.Status, project.LastUpdateTime)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		batch.Query(fmt.Sprintf(updateConfigProjectQuery, a.keyspace, tablename, project.Status, project.LastUpdateTime, project.Id, project.WorkflowId))
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) DeleteConfigProject(id string, tablename string) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(deleteConfigProjectQuery, a.keyspace, tablename, id))
	err := a.client.Session().ExecuteBatch(batch)
	return err
}

func (a *Store) GetConfigProjectForID(id string, tablename string) (*model.ConfigureProject, error) {
	query := fmt.Sprintf(getConfigProjectsForIDQuery, a.keyspace, tablename, id)
	projects, err := a.executeConfigProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) != 1 {
		return nil, fmt.Errorf(objectNotFoundErrorMessage)
	}
	return projects[0], nil
}

func (a *Store) GetConfigProjects(tablename string) ([]*model.ConfigureProject, error) {
	return a.updateConfigProjects(tablename)
}

// tektonProjects, err := a.GetGitProjectsByLabels([]string{"tekton"})
func (a *Store) updateConfigProjects(tablename string) ([]*model.ConfigureProject, error) {
	allTektonProjects, err := a.GetGitProjectsByLabels([]string{"tekton"})
	if err != nil {
		a.log.Errorf("failed to fetch all tekton projects, :%v", err)
		return nil, err
	}

	query := fmt.Sprintf(getConfigProjectsQuery, a.keyspace, tablename)
	regTektonProjects, err := a.executeConfigProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	regTektonProjectId := make(map[string]*model.ConfigureProject)
	for _, tekPro := range regTektonProjects {
		regTektonProjectId[tekPro.Id] = tekPro
	}

	ret := make([]*model.ConfigureProject, 0)
	for _, allTekProject := range allTektonProjects {
		project := &model.ConfigureProject{Id: allTekProject.Id, GitProjectId: allTekProject.Id,
			GitProjectUrl: allTekProject.ProjectUrl}
		if _, ok := regTektonProjectId[allTekProject.Id]; !ok {
			project.Status = "available"
			project.WorkflowId = "NA"
			if err := a.UpsertConfigProject(project, tablename); err != nil {
				return nil, err
			}
		} else {
			project.Status = regTektonProjectId[allTekProject.Id].Status
		}
		ret = append(ret, project)
	}

	return ret, nil
}

func (a *Store) executeConfigProjectsSelectQuery(query string) ([]*model.ConfigureProject, error) {
	selectAllQuery := a.client.Session().Query(query)
	iter := selectAllQuery.Iter()
	project := model.ConfigureProject{}

	ret := make([]*model.ConfigureProject, 0)
	for iter.Scan(
		&project.Id, &project.GitProjectId, &project.Status, &project.LastUpdateTime) {
		gitProject, err := a.GetGitProjectForID(project.Id)
		if err != nil {
			a.log.Errorf("tekton project %s not exist in git projects", project.Id)
			continue
		}

		a := &model.ConfigureProject{
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
