package captenstore

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/pkg/errors"
)

const (
	insertTektonPipelines        = "INSERT INTO %s.TektonPipelines(id, pipeline_name, git_org_id, container_registry_id, last_update_time) VALUES (?,?,?,?,?)"
	updateTektonPipelinesById    = "UPDATE %s.TektonPipelines SET %s WHERE id=?"
	selectAllTektonPipelines     = "SELECT id, pipeline_name, git_org_id, container_registry_id, last_update_time FROM %s.TektonPipelines"
	selectGetTektonPipelinesById = "SELECT id, pipeline_name, git_org_id, container_registry_id, last_update_time FROM %s.TektonPipelines WHERE id=%s;"
)

func (a *Store) UpsertTektonPipelines(config *captenpluginspb.TektonPipelines) error {
	config.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertTektonPipelines, a.keyspace), config.Id,
		config.PipelineName, config.GitOrgId, config.ContainerRegistryId, config.LastUpdateTime)
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

func (a *Store) GetTektonPipelinesForID(id string) (*captenpluginspb.TektonPipelines, error) {
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

func (a *Store) GetTektonPipeliness() ([]*captenpluginspb.TektonPipelines, error) {
	query := fmt.Sprintf(selectAllTektonPipelines, a.keyspace)
	return a.executeTektonPipelinessSelectQuery(query)
}

func (a *Store) executeTektonPipelinessSelectQuery(query string) ([]*captenpluginspb.TektonPipelines, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	project := captenpluginspb.TektonPipelines{}

	ret := make([]*captenpluginspb.TektonPipelines, 0)
	for iter.Scan(
		&project.Id, &project.PipelineName,
		&project.GitOrgId, &project.ContainerRegistryId, &project.Status, &project.LastUpdateTime,
	) {
		TektonPipelines := &captenpluginspb.TektonPipelines{
			Id:                  project.Id,
			PipelineName:        project.PipelineName,
			LastUpdateTime:      project.LastUpdateTime,
			GitOrgId:            project.GitOrgId,
			ContainerRegistryId: project.ContainerRegistryId,
			Status:              project.Status,
		}
		ret = append(ret, TektonPipelines)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return ret, nil
}

func formUpdateKvPairsForTektonPipelines(config *captenpluginspb.TektonPipelines) (updatePlaceholders string, values []interface{}) {
	params := []string{}

	if config.GitOrgId != "" {
		params = append(params, "git_org_id = ?")
		values = append(values, config.GitOrgId)
	}

	if config.LastUpdateTime != "" {
		params = append(params, "last_update_time = ?")
		values = append(values, config.LastUpdateTime)
	}

	if len(config.ContainerRegistryId) != 0 {
		params = append(params, "container_registry_id = ?")
		values = append(values, config.ContainerRegistryId)
	}

	if config.Status != "" {
		params = append(params, "status = ?")
		values = append(values, config.Status)
	}

	if len(params) == 0 {
		return "", nil
	}
	return strings.Join(params, ", "), values
}
