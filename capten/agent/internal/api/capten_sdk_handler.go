package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/kube-tarian/kad/capten/agent/gin-api-server/api"
)

func (a *Agent) GetGitProjectById(c *gin.Context, projectID string) {
	if projectID == "" {
		a.log.Error("Project Id is not provided")
		c.String(http.StatusBadRequest, "project Id is not provided")
		return
	}

	a.log.Infof("Get Git project By Id request recieved for Id: %s", projectID)

	res, err := a.as.GetGitProjectForID(projectID)
	if err != nil {
		a.log.Errorf("failed to get gitProject from db for project Id: %s, %v", projectID, err)
		c.String(http.StatusInternalServerError, "failed to fetch git project for "+projectID)
		return
	}

	accessToken, _, _, _, err := a.getGitProjectCredential(context.TODO(), res.Id)
	if err != nil {
		a.log.Errorf("failed to get git credential for project Id: %s, %v", projectID, err)
		c.String(http.StatusInternalServerError, "failed to fetch git project for "+projectID)
		return
	}

	a.log.Infof("Fetched %s git project", res.Id)

	c.IndentedJSON(http.StatusOK, &api.GitProjectResponse{
		Project: api.GitProject{
			AccessToken:    accessToken,
			Id:             res.Id,
			Labels:         res.Labels,
			LastUpdateTime: res.LastUpdateTime,
			ProjectUrl:     res.ProjectUrl,
		},
		Status:        api.OK,
		StatusMessage: "successfully fetched git project for " + projectID,
	})
}

func (a *Agent) GetContainerRegistryById(c *gin.Context, id string) {
	if id == "" {
		a.log.Error("Container registry Id is not provided")
		c.String(http.StatusBadRequest, "container registry Id is not provided")
		return
	}

	a.log.Infof("Get Container registry By Id request recieved for Id: %s", id)

	res, err := a.as.GetContainerRegistryForID(id)
	if err != nil {
		a.log.Errorf("failed to get ContainerRegistry from db, %v", err)
		c.String(http.StatusInternalServerError, "failed to fetch container registry for "+id)
	}

	cred, _, _, err := a.getContainerRegCredential(context.TODO(), res.Id)
	if err != nil {
		a.log.Errorf("failed to get container registry credential for %s, %v", id, err)
		c.String(http.StatusInternalServerError, "failed to fetch container registry for "+id)
		return
	}

	a.log.Infof("Fetched %s container registry", id)
	c.IndentedJSON(http.StatusOK,
		&api.ContainerRegistryResponse{
			Registry: &api.ContainerRegistry{
				Id:                 res.Id,
				RegistryUrl:        res.RegistryUrl,
				Labels:             res.Labels,
				LastUpdateTime:     res.LastUpdateTime,
				RegistryType:       res.RegistryType,
				RegistryAttributes: cred,
			},
		})
}
