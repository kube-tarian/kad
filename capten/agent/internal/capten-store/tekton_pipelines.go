package captenstore

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

const (
	insertTektonPipelines        = "INSERT INTO %s.TektonPipelines(id, pipeline_name, git_org_id, container_registry_id, managed_cluster_id, crossplane_git_project_id, status, last_update_time, workflow_id, workflow_status) VALUES (?,?,?,?,?,?,?,?,?,?)"
	updateTektonPipelinesById    = "UPDATE %s.TektonPipelines SET %s WHERE id=?"
	deleteTektonPipelinesById    = "DELETE FROM %s.TektonPipelines WHERE id= ?"
	selectAllTektonPipelines     = "SELECT id, pipeline_name, git_org_id, container_registry_id, managed_cluster_id, crossplane_git_project_id, status, last_update_time FROM %s.TektonPipelines"
	selectGetTektonPipelinesById = "SELECT id, pipeline_name, git_org_id, container_registry_id, managed_cluster_id, crossplane_git_project_id, status, last_update_time FROM %s.TektonPipelines WHERE id=%s;"
)

func (a *Store) UpsertTektonPipelines(config *model.TektonPipeline) error {
	config.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertTektonPipelines, a.keyspace), config.Id,
		config.PipelineName, config.GitProjectId, config.ContainerRegId, config.Status,
		config.LastUpdateTime, config.WorkflowId, config.WorkflowStatus)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		updatePlaceholders, values := formUpdateKvPairsForTektonPipelines(config)
		if updatePlaceholders == "" {
			return err
		}
		query := fmt.Sprintf(updateTektonPipelinesById, a.keyspace, updatePlaceholders)
		args := append(values, config.Id)
		batch = a.client.Session().NewBatch(gocql.LoggedBatch)
		batch.Query(query, args...)
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) GetTektonPipelinesForID(id string) (*model.TektonPipeline, error) {
	query := fmt.Sprintf(selectGetTektonPipelinesById, a.keyspace, id)
	projects, err := a.executeTektonPipelinessSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) != 1 {
		return nil, fmt.Errorf("pipelines not found")
	}
	return projects[0], nil
}

func (a *Store) DeleteTektonPipelinesById(id string) error {
	deleteAction := a.client.Session().Query(fmt.Sprintf(deleteTektonPipelinesById,
		a.keyspace), id)
	err := deleteAction.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *Store) GetTektonPipeliness() ([]*model.TektonPipeline, error) {
	query := fmt.Sprintf(selectAllTektonPipelines, a.keyspace)
	return a.executeTektonPipelinessSelectQuery(query)
}

func (a *Store) executeTektonPipelinessSelectQuery(query string) ([]*model.TektonPipeline, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	project := &model.TektonPipeline{}

	ret := make([]*model.TektonPipeline, 0)
	for iter.Scan(
		&project.Id, &project.PipelineName,
		&project.GitProjectId, &project.ContainerRegId, &project.ManagedClusterId, &project.CrossplaneGitProjectId,
		&project.Status, &project.LastUpdateTime,
	) {
		TektonPipelines := &model.TektonPipeline{
			Id:                     project.Id,
			PipelineName:           project.PipelineName,
			LastUpdateTime:         project.LastUpdateTime,
			GitProjectId:           project.GitProjectId,
			ContainerRegId:         project.ContainerRegId,
			ManagedClusterId:       project.ManagedClusterId,
			CrossplaneGitProjectId: project.CrossplaneGitProjectId,
			Status:                 project.Status,
		}
		ret = append(ret, TektonPipelines)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return ret, nil
}

func formUpdateKvPairsForTektonPipelines(config *model.TektonPipeline) (updatePlaceholders string, values []interface{}) {
	params := []string{}

	if config.GitProjectId != "" {
		params = append(params, "git_org_id = ?")
		values = append(values, config.GitProjectId)
	}

	if config.LastUpdateTime != "" {
		params = append(params, "last_update_time = ?")
		values = append(values, config.LastUpdateTime)
	}

	if len(config.ContainerRegId) != 0 {
		params = append(params, "container_registry_id = ?")
		values = append(values, config.ContainerRegId)
	}

	if config.ManagedClusterId != "" {
		params = append(params, "managed_cluster_id = ?")
		values = append(values, config.Status)
	}

	if config.CrossplaneGitProjectId != "" {
		params = append(params, "crossplane_git_project_id = ?")
		values = append(values, config.Status)
	}

	if config.Status != "" {
		params = append(params, "status = ?")
		values = append(values, config.Status)
	}

	if config.WorkflowStatus != "" {
		params = append(params, "workflow_status = ?")
		values = append(values, config.Status)
	}

	if config.WorkflowId != "" {
		params = append(params, "workflow_id = ?")
		values = append(values, config.Status)
	}

	if len(params) == 0 {
		return "", nil
	}
	return strings.Join(params, ", "), values
}
