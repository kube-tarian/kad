package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/logging"
)

type APIHanlder struct {
	log    logging.Logger
	client *client.Agent
}

func NewAPIHandler(log logging.Logger) (*APIHanlder, error) {
	client, err := client.NewAgent()
	if err != nil {
		log.Errorf("failed to connect agent internal error", err)
		return nil, err
	}

	return &APIHanlder{
		log:    log,
		client: client,
	}, nil
}

func (s *APIHanlder) Close(c *gin.Context) {
	s.client.Close()
}

func (ah *APIHanlder) GetApiDocs(c *gin.Context) {
	swagger, err := api.GetSwagger()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.IndentedJSON(http.StatusOK, swagger)
}

func (ah *APIHanlder) GetStatus(c *gin.Context) {
	c.String(http.StatusOK, "")
}
