package captenstore

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/pkg/errors"
)

const (
	getCrossplaneProjectsQuery      = "SELECT id, git_project_id, status, last_update_time, workflow_id, workflow_status FROM %s.CrossplaneProjects;"
	getCrossplaneProjectsForIDQuery = "SELECT id, git_project_id, status, last_update_time, workflow_id, workflow_status FROM %s.CrossplaneProjects WHERE id=%s;"
	insertCrossplaneProjectQuery    = "INSERT INTO %s.CrossplaneProjects(id, git_project_id, status, last_update_time) VALUES (?,?,?,?);"
	updateCrossplaneProjectQuery    = "UPDATE %s.CrossplaneProjects SET status='%s', last_update_time='%s', workflow_id='%s', workflow_status='%s' WHERE id=%s;"
	deleteCrossplaneProjectQuery    = "DELETE FROM %s.CrossplaneProjects WHERE id=%s;"
)

func (a *Store) UpsertCrossplaneProject(project *model.CrossplaneProject) error {
	project.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertCrossplaneProjectQuery, a.keyspace), project.Id, project.GitProjectId, project.Status, project.LastUpdateTime)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		batch.Query(fmt.Sprintf(updateCrossplaneProjectQuery, a.keyspace, project.Status, project.LastUpdateTime, project.WorkflowId, project.WorkflowStatus, project.Id))
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) DeleteCrossplaneProject(id string) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(deleteCrossplaneProjectQuery, a.keyspace, id))
	err := a.client.Session().ExecuteBatch(batch)
	return err
}

func (a *Store) GetCrossplaneProjectForID(id string) (*model.CrossplaneProject, error) {
	query := fmt.Sprintf(getCrossplaneProjectsForIDQuery, a.keyspace, id)
	projects, err := a.executeCrossplaneProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf(objectNotFoundErrorMessage)
	}
	return projects[0], nil
}

func (a *Store) GetCrossplaneProject() (*model.CrossplaneProject, error) {
	return a.updateCrossplaneProject()
}

// CrossplaneProjects, err := a.GetGitProjectsByLabels([]string{"Crossplane"})
func (a *Store) updateCrossplaneProject() (*model.CrossplaneProject, error) {
	allCrossplaneGitProjects, err := a.GetGitProjectsByLabels([]string{"crossplane"})
	if err != nil {
		a.log.Errorf("failed to fetch all Crossplane projects, :%v", err)
		return nil, err
	}

	query := fmt.Sprintf(getCrossplaneProjectsQuery, a.keyspace)
	regCrossplaneProjects, err := a.executeCrossplaneProjectsSelectQuery(query)
	if err != nil {
		return nil, err
	}

	regCrossplaneProjectForGitProjectId := make(map[string]*model.CrossplaneProject)
	for _, crossplanePro := range regCrossplaneProjects {
		regCrossplaneProjectForGitProjectId[crossplanePro.GitProjectId] = crossplanePro
	}

	ret := make([]*model.CrossplaneProject, 0)

	for _, crossplaneGitProject := range allCrossplaneGitProjects {

		project := &model.CrossplaneProject{
			GitProjectId:  crossplaneGitProject.Id,
			GitProjectUrl: crossplaneGitProject.ProjectUrl}

		if _, found := regCrossplaneProjectForGitProjectId[crossplaneGitProject.Id]; !found {
			project.Status = "available"
			project.WorkflowId = "NA"
			project.WorkflowStatus = "NA"
			project.LastUpdateTime = time.Now().Format(time.RFC3339)
			project.Id = uuid.New().String()
			if err := a.UpsertCrossplaneProject(project); err != nil {
				return nil, err
			}
		} else {
			crossplaneProject := regCrossplaneProjectForGitProjectId[crossplaneGitProject.Id]

			project.Id = crossplaneProject.Id
			project.WorkflowId = crossplaneProject.WorkflowId
			project.WorkflowStatus = crossplaneProject.WorkflowStatus
			project.Status = crossplaneProject.Status
		}
		ret = append(ret, project)
	}
	return ret[0], nil
}

func (a *Store) executeCrossplaneProjectsSelectQuery(query string) ([]*model.CrossplaneProject, error) {
	selectAllQuery := a.client.Session().Query(query)
	iter := selectAllQuery.Iter()
	project := model.CrossplaneProject{}

	ret := make([]*model.CrossplaneProject, 0)
	for iter.Scan(
		&project.Id, &project.GitProjectId, &project.Status,
		&project.LastUpdateTime, &project.WorkflowId, &project.WorkflowStatus) {

		gitProject, err := a.GetGitProjectForID(project.Id)
		if err != nil {
			a.log.Errorf("Crossplane project %s not exist in git projects", project.Id)
			continue
		}

		a := &model.CrossplaneProject{
			Id:             project.Id,
			GitProjectId:   gitProject.Id,
			GitProjectUrl:  gitProject.ProjectUrl,
			Status:         project.Status,
			LastUpdateTime: project.LastUpdateTime,
			WorkflowId:     project.WorkflowId,
			WorkflowStatus: project.WorkflowStatus,
		}
		ret = append(ret, a)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}
	return ret, nil
}
