package captenstore

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/pkg/errors"
)

const (
	insertGitProject             = "INSERT INTO %s.GitProjects(id, project_url, labels, last_update_time) VALUES (?,?,?,?)"
	insertGitProjectId           = "INSERT INTO %s.GitProjects(id) VALUES (?)"
	updateGitProjectById         = "UPDATE %s.GitProjects SET %s WHERE id = ?"
	deleteGitProjectById         = "DELETE FROM %s.GitProjects WHERE id= ?"
	selectAllGitProjects         = "SELECT id, project_url, labels, last_update_time FROM %s.GitProjects"
	selectAllGitProjectsByLabels = "SELECT id, project_url, labels, last_update_time FROM %s.GitProjects WHERE %s"
	selectGetGitProjectById      = "SELECT id, project_url, labels, last_update_time FROM %s.GitProjects WHERE id=%s;"
)

func (a *Store) UpsertGitProject(config *captenpluginspb.GitProject) error {
	config.LastUpdateTime = time.Now().Format(time.RFC3339)
	kvPairs := formUpdateKvPairsForGitProject(config)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertGitProject, a.keyspace), config.Id, config.ProjectUrl, config.Labels, config.LastUpdateTime)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		batch.Query(fmt.Sprintf(updateGitProjectById, a.keyspace, kvPairs), config.Id)
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
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

func (a *Store) GetGitProjectForID(id string) (*captenpluginspb.GitProject, error) {
	query := fmt.Sprintf(selectGetGitProjectById, a.keyspace, id)
	projects, err := a.executeGitProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) != 1 {
		return nil, fmt.Errorf("project not found")
	}
	return projects[0], nil
}

func (a *Store) GetGitProjects() ([]*captenpluginspb.GitProject, error) {
	query := fmt.Sprintf(selectAllGitProjects, a.keyspace)
	return a.executeGitProjectsSelectQuery(query)
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
	return a.executeGitProjectsSelectQuery(query)
}

func (a *Store) executeGitProjectsSelectQuery(query string) ([]*captenpluginspb.GitProject, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	project := captenpluginspb.GitProject{}
	var labels []string

	ret := make([]*captenpluginspb.GitProject, 0)
	for iter.Scan(
		&project.Id, &project.ProjectUrl,
		&labels, &project.LastUpdateTime,
	) {
		labelsTmp := make([]string, len(labels))
		copy(labelsTmp, labels)
		gitProject := &captenpluginspb.GitProject{
			Id:             project.Id,
			ProjectUrl:     project.ProjectUrl,
			Labels:         labelsTmp,
			LastUpdateTime: project.LastUpdateTime,
		}
		ret = append(ret, gitProject)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return ret, nil
}

func formUpdateKvPairsForGitProject(config *captenpluginspb.GitProject) (kvPairs string) {
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
		return ""
	}
	return strings.Join(params, ", ")
}
