package captenstore

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/pkg/errors"
)

const (
	getTektonProjectsQuery      = "SELECT id, git_project_id, status, lastUpdateTime FROM %s.tekton;"
	getTektonProjectsForIDQuery = "SELECT id, git_project_id, status, lastUpdateTime FROM %s.tekton WHERE id='%s';"
	insertTektonProjectQuery    = "INSERT INTO %s.tekton(id, git_project_id, status, lastUpdateTime) VALUES (?,?,?);"
	updateTektonProjectQuery    = "UPDATE %s.tekton SET status='%s', lastUpdateTime='%s' WHERE id='%s' and git_project_id='%s';"
	deleteTektonProjectQuery    = "DELETE FROM %s.tekton WHERE id='%s';"
)

func (a *Store) UpsertTektonProject(payload *model.TektonProject) error {
	payload.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertTektonProjectQuery, a.keyspace), payload.Id, payload.GitProjectId, payload.Status, payload.LastUpdateTime)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		batch.Query(fmt.Sprintf(updateTektonProjectQuery, a.keyspace, payload.Status, payload.LastUpdateTime, payload.Id, payload.GitProjectId))
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
		return nil, fmt.Errorf("project not found")
	}
	return projects[0], nil
}

func (a *Store) GetTektonProjects() ([]*model.TektonProject, error) {
	query := fmt.Sprintf(getTektonProjectsQuery, a.keyspace)
	return a.executeTektonProjectsSelectQuery(query)
}

func (a *Store) executeTektonProjectsSelectQuery(query string) ([]*model.TektonProject, error) {
	selectAllQuery := a.client.Session().Query(query)
	iter := selectAllQuery.Iter()
	project := model.TektonProject{}

	ret := make([]*model.TektonProject, 0)
	for iter.Scan(
		&project.Id, project.Status, &project.LastUpdateTime) {
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
