package captenstore

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/pkg/errors"
)

const (
	insertGitProject             = "INSERT INTO %s.GitProjects(id, project_url, labels, last_update_time) VALUES (?)"
	insertGitProjectId           = "INSERT INTO %s.GitProjects(id) VALUES (?)"
	updateGitProjectById         = "UPDATE %s.GitProjects SET %s WHERE id = ?"
	deleteGitProjectById         = "DELETE FROM %s.GitProjects WHERE id= ?"
	selectAllGitProjects         = "SELECT id, project_url, labels, last_update_time FROM %s.GitProjects"
	selectAllGitProjectsByLabels = "SELECT id, project_url, labels, last_update_time FROM %s.GitProjects WHERE %s"

	getGitProjectsById = "SELECT id, project_url, labels, last_update_time FROM %s.GitProjects WHERE id='%s'"
	getTektonQuery     = "SELECT git_project_id, status, status FROM %s.tekton;"
	insertTektonQuery  = "INSERT INTO %s.tekton(git_project_id, status, lastUpdateTime) VALUES (?,?,?);"
	updateTektonQuery  = "UPDATE %s.tekton SET status='%s', lastUpdateTime='%s' WHERE git_project_id='%s';"
	deleteTektonQuery  = "DELETE FROM %s.tekton WHERE git_project_id='%s';"
)

func (a *Store) UpsertGitProject(config *captenpluginspb.GitProject) error {
	config.LastUpdateTime = time.Now().Format(time.RFC3339)
	kvPairs, isEmptyUpdate := formUpdateKvPairsForGitProject(config)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertGitProjectId, a.keyspace), config.Id)
	if !isEmptyUpdate {
		batch.Query(fmt.Sprintf(updateGitProjectById, a.keyspace, kvPairs), config.Id)
	}
	return a.client.Session().ExecuteBatch(batch)
}

func (a *Store) DeleteGitProjectById(id string) error {
	deleteAction := a.client.Session().Query(fmt.Sprintf(deleteGitProjectById,
		a.keyspace), id)
	err := deleteAction.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *Store) GetGitProject(id string) ([]*captenpluginspb.GitProject, error) {
	query := fmt.Sprintf(getGitProjectsById, a.keyspace, id)
	return a.executeSelectQuery(query)
}

func (a *Store) GetGitProjects() ([]*captenpluginspb.GitProject, error) {
	query := fmt.Sprintf(selectAllGitProjects, a.keyspace)
	return a.executeSelectQuery(query)
}

func (a *Store) GetGitProjectsByLabels(searchLabels []string) ([]*captenpluginspb.GitProject, error) {
	if len(searchLabels) == 0 {
		return nil, fmt.Errorf("searchLabels empty")
	}

	labelContains := []string{}
	for _, label := range searchLabels {
		labelContains = append(labelContains, fmt.Sprintf("labels CONTAINS '%s'", label))
	}
	whereLabelsClause := strings.Join(labelContains, " OR ")
	whereLabelsClause += " ALLOW FILTERING"
	query := fmt.Sprintf(selectAllGitProjectsByLabels, a.keyspace, whereLabelsClause)
	return a.executeSelectQuery(query)
}

func (a *Store) executeSelectQuery(query string) ([]*captenpluginspb.GitProject, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	config := captenpluginspb.GitProject{}
	var labels []string

	ret := make([]*captenpluginspb.GitProject, 0)
	for iter.Scan(
		&config.Id, &config.ProjectUrl,
		&labels, &config.LastUpdateTime,
	) {
		labelsTmp := make([]string, len(labels))
		copy(labelsTmp, labels)
		gitProject := &captenpluginspb.GitProject{
			Id:             config.Id,
			ProjectUrl:     config.ProjectUrl,
			Labels:         labelsTmp,
			LastUpdateTime: config.LastUpdateTime,
		}
		ret = append(ret, gitProject)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return ret, nil
}

func formUpdateKvPairsForGitProject(config *captenpluginspb.GitProject) (kvPairs string, emptyUpdate bool) {
	params := []string{}

	if config.ProjectUrl != "" {
		params = append(params,
			fmt.Sprintf("project_url = '%s'", config.ProjectUrl))
	}

	// comma separated labels, change this later
	if len(config.Labels) > 0 {
		labels := []string{}
		for _, label := range config.Labels {
			labels = append(labels, fmt.Sprintf("'%s'", label))
		}
		param := "{" + strings.Join(labels, ", ") + "}"
		params = append(params,
			fmt.Sprintf("labels = %v", param))
	}

	if (config.LastUpdateTime) != "" {
		params = append(params,
			fmt.Sprintf("last_update_time = '%v'", config.LastUpdateTime))
	}

	if len(params) == 0 {
		// query is empty there is nothing to update
		return "", true
	}
	return strings.Join(params, ", "), false
}

func (a *Store) AddTektonProject(payload *model.RegisterTekton) error {

	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertTektonQuery, a.keyspace), payload.Id, payload.Status, gocql.TimeUUID().String())

	err := a.client.Session().ExecuteBatch(batch)

	return err
}

func (a *Store) UpdateTektonProject(payload *model.RegisterTekton) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(updateTektonQuery, a.keyspace, payload.Status, gocql.TimeUUID().String(), payload.Id))
	err := a.client.Session().ExecuteBatch(batch)

	return err
}

func (a *Store) DeleteTektonProject(id string) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(deleteTektonQuery, a.keyspace, id))
	err := a.client.Session().ExecuteBatch(batch)

	return err
}

func (a *Store) GetTektonProjects() ([]*captenpluginspb.TektonProject, error) {
	selectAllQuery := a.client.Session().Query(fmt.Sprintf(getTektonQuery, a.keyspace))
	iter := selectAllQuery.Iter()

	config := model.RegisterTekton{}

	ret := make([]*captenpluginspb.TektonProject, 0)
	for iter.Scan(
		&config.Id, config.Status, &config.LastUpdateTime) {
		configCopy := config
		project, err := a.GetGitProject(configCopy.Id)
		if err != nil || len(project) != 1 {
			return nil, errors.WithMessage(err, "failed to get git projects:")
		}

		a := &captenpluginspb.TektonProject{
			Id: configCopy.Id, Status: configCopy.Status,
			GitProjectUrl: project[0].ProjectUrl,
		}

		ret = append(ret, a)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return ret, nil

}
