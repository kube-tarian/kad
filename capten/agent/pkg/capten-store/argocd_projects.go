package captenstore

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
)

const (
	getArgocdProjectsDataQuery    = "SELECT id, last_update_time, git_project_id, status FROM %s.argocd_projects_data WHERE git_project_id='%s';"
	insertArgocdProjectsDataQuery = "INSERT INTO %s.argocd_projects_data(last_update_time, git_project_id, status) VALUES (?,?,?);"
	updateArgocdProjectsDataQuery = "UPDATE %s.argocd_projects_data SET %s WHERE id='%s';"
	deleteArgocdProjectsDataQuery = "DELETE FROM %s.argocd_projects_data WHERE git_project_id='%s';"
)

func (a *Store) AddArgoCDProjectsData(gitProjectId, status string) error {

	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertArgocdProjectsDataQuery, a.keyspace), time.Now(), gitProjectId, status)

	err := a.client.Session().ExecuteBatch(batch)

	return err
}

func (a *Store) GetArgoCDProjectsData(id string) (*model.ArgoCDProjectsData, error) {

	selectQuery := a.client.Session().Query(fmt.Sprintf(getArgocdProjectsDataQuery,
		a.keyspace, id))

	data := model.ArgoCDProjectsData{}

	if err := selectQuery.Scan(
		&data.Id, &data.GitProjectID, &data.LastUpdateTime, &data.Status,
	); err != nil {
		return nil, err
	}

	return &data, nil
}

func (a *Store) DeleteArgoCDProjectsData(id string) error {

	deleteQuery := a.client.Session().Query(fmt.Sprintf(deleteArgocdProjectsDataQuery,
		a.keyspace, id))

	err := deleteQuery.Exec()
	if err != nil {
		return err
	}

	return nil
}
