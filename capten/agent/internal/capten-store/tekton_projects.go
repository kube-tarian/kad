package captenstore

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

const (
	getTektonProjectsQuery      = "SELECT id, git_project_id, status, last_update_time, workflow_id, workflow_status FROM %s.TektonProjects;"
	getTektonProjectsForIDQuery = "SELECT id, git_project_id, status, last_update_time, workflow_id, workflow_status FROM %s.TektonProjects WHERE id=%s;"
	insertTektonProjectQuery    = "INSERT INTO %s.TektonProjects(id, git_project_id, status, last_update_time, workflow_id, workflow_status) VALUES (?,?,?,?,?,?);"
	updateTektonProjectQuery    = "UPDATE %s.TektonProjects SET status=?, last_update_time=?, workflow_id=?, workflow_status=? WHERE id=?;"
	deleteTektonProjectQuery    = "DELETE FROM %s.TektonProjects WHERE id=%s;"
)

func (a *Store) UpsertTektonProject(project *model.TektonProject) error {
	project.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertTektonProjectQuery, a.keyspace), project.Id, project.GitProjectId, project.Status, project.LastUpdateTime, project.WorkflowId, project.WorkflowStatus)
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		batch = a.client.Session().NewBatch(gocql.LoggedBatch)
		query := fmt.Sprintf(updateTektonProjectQuery, a.keyspace)
		batch.Query(query, project.Status, project.LastUpdateTime, project.WorkflowId, project.WorkflowStatus, project.Id)
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

func (a *Store) GetTektonProject() (*model.TektonProject, error) {
	return a.updateTektonProject()
}

func (a *Store) updateTektonProject() (*model.TektonProject, error) {
	allTektonProjects, err := a.GetGitProjectsByLabels([]string{"tekton"})
	if err != nil {
		a.log.Errorf("failed to fetch all tekton projects, :%v", err)
		return nil, err
	}

	if len(allTektonProjects) == 0 {
		return nil, fmt.Errorf("no git project found with tekton tag")
	}
	tektonGitProject := allTektonProjects[0]

	query := fmt.Sprintf(getTektonProjectsQuery, a.keyspace)
	regTektonProjects, err := a.executeTektonProjectsSelectQuery(query)
	if err != nil {
		a.log.Errorf("failed to execute select tekton projects, :%v", err)
		return nil, err
	}

	regTektonProjectId := make(map[string]*model.TektonProject)
	for _, tekPro := range regTektonProjects {
		var deleteRecord = true
		for _, gitProject := range allTektonProjects {
			if gitProject.Id == tekPro.GitProjectId {
				deleteRecord = false
				break
			}
		}

		if deleteRecord {
			if err := a.DeleteTektonProject(tekPro.Id); err != nil {
				return nil, err
			}

			for _, gitProject := range allTektonProjects {
				var deleteRecord = true
				if tekPro.GitProjectId == gitProject.Id {
					deleteRecord = false
					break
				}

				if deleteRecord {
					tags := []string{}
					for _, v := range gitProject.Tags {
						if v != "tekton" {
							tags = append(tags, v)
						}
					}

					if len(tags) > 0 {
						gitProject.Tags = tags
						if err := a.UpsertGitProject(gitProject); err != nil {
							return nil, err
						}
					}
				}
			}
		} else {
			regTektonProjectId[tekPro.Id] = tekPro
		}
	}

	if len(regTektonProjectId) == 0 {
		// no project was registered, register the git project
		project := &model.TektonProject{
			Id:             uuid.New().String(),
			GitProjectId:   tektonGitProject.Id,
			GitProjectUrl:  tektonGitProject.ProjectUrl,
			Status:         "available",
			WorkflowId:     "NA",
			WorkflowStatus: "NA",
			LastUpdateTime: time.Now().Format(time.RFC3339),
		}
		if err := a.UpsertTektonProject(project); err != nil {
			return nil, err
		}

		tektonGitProject.Tags = append(tektonGitProject.Tags, "tekton")
		if err := a.UpsertGitProject(tektonGitProject); err != nil {
			return nil, err
		}

		return project, nil
	}
	return regTektonProjects[0], nil
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
			a.log.Debugf("tekton project %s not exist in git projects, %v", project.Id, err)
			continue
		}

		a := &model.TektonProject{
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
