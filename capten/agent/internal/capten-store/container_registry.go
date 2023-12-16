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
	insertContainerRegistry        = "INSERT INTO %s.ContainerRegistry(id, registry_url, labels, last_update_time, token, username, password, vendor) VALUES (?,?,?,?,?,?,?,?)"
	updateContainerRegistryById    = "UPDATE %s.ContainerRegistry SET %s WHERE id=?"
	deleteContainerRegistryById    = "DELETE FROM %s.ContainerRegistry WHERE id= ?"
	selectAllContainerRegistrys    = "SELECT id, registry_url, labels, last_update_time, token, username, password, vendor FROM %s.ContainerRegistry"
	selectGetContainerRegistryById = "SELECT id, registry_url, labels, last_update_time, token, username, password, vendor FROM %s.ContainerRegistry WHERE id=%s;"
)

func (a *Store) UpsertContainerRegistry(config *captenpluginspb.ContainerRegistry) error {
	config.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertContainerRegistry, a.keyspace), config.Id,
		config.RegistryUrl, config.Labels, config.LastUpdateTime, config.Token, config.UserName, config.Password, config.Vendor)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		updatePlaceholders, values := formUpdateKvPairsForContainerRegistry(config)
		if updatePlaceholders == "" {
			return err
		}
		query := fmt.Sprintf(updateContainerRegistryById, a.keyspace, updatePlaceholders)
		args := append(values, config.Id)
		batch = a.client.Session().NewBatch(gocql.LoggedBatch)
		batch.Query(query, args...)
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) DeleteContainerRegistryById(id string) error {
	deleteAction := a.client.Session().Query(fmt.Sprintf(deleteContainerRegistryById,
		a.keyspace), id)
	err := deleteAction.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *Store) GetContainerRegistryForID(id string) (*captenpluginspb.ContainerRegistry, error) {
	query := fmt.Sprintf(selectGetContainerRegistryById, a.keyspace, id)
	projects, err := a.executeContainerRegistrysSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) != 1 {
		return nil, fmt.Errorf("project not found")
	}
	return projects[0], nil
}

func (a *Store) GetContainerRegistrys() ([]*captenpluginspb.ContainerRegistry, error) {
	query := fmt.Sprintf(selectAllContainerRegistrys, a.keyspace)
	return a.executeContainerRegistrysSelectQuery(query)
}

func (a *Store) executeContainerRegistrysSelectQuery(query string) ([]*captenpluginspb.ContainerRegistry, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	project := captenpluginspb.ContainerRegistry{}

	ret := make([]*captenpluginspb.ContainerRegistry, 0)
	for iter.Scan(
		&project.Id, &project.RegistryUrl,
		&project.Labels, &project.LastUpdateTime, &project.Token, &project.UserName,
		&project.Password, &project.Vendor,
	) {
		ContainerRegistry := &captenpluginspb.ContainerRegistry{
			Id:             project.Id,
			RegistryUrl:    project.RegistryUrl,
			Labels:         project.Labels,
			LastUpdateTime: project.LastUpdateTime,
			Token:          project.Token,
			UserName:       project.UserName,
			Password:       project.Password,
			Vendor:         project.Vendor,
		}
		ret = append(ret, ContainerRegistry)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return ret, nil
}

func formUpdateKvPairsForContainerRegistry(config *captenpluginspb.ContainerRegistry) (updatePlaceholders string, values []interface{}) {
	params := []string{}

	if config.RegistryUrl != "" {
		params = append(params, "registry_url = ?")
		values = append(values, config.RegistryUrl)
	}

	if len(config.Labels) != 0 {
		params = append(params, "labels = ?")
		values = append(values, config.Labels)
	}

	if config.LastUpdateTime != "" {
		params = append(params, "last_update_time = ?")
		values = append(values, config.LastUpdateTime)
	}

	if config.Token != "" {
		params = append(params, "token = ?")
		values = append(values, config.Token)
	}
	if config.UserName != "" {
		params = append(params, "username = ?")
		values = append(values, config.UserName)
	}
	if config.Password != "" {
		params = append(params, "password = ?")
		values = append(values, config.Password)
	}
	if config.Vendor != "" {
		params = append(params, "vendor = ?")
		values = append(values, config.Vendor)
	}

	if len(params) == 0 {
		return "", nil
	}
	return strings.Join(params, ", "), values
}
