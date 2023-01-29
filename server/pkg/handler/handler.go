package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/logging"
)

type APIHanlder struct {
	log    logging.Logger
	client *client.Agent
}

var (
	apiOnce sync.Once
)

func NewAPIHandler(log logging.Logger) (*APIHanlder, error) {
	return &APIHanlder{
		log:    log,
		client: nil,
	}, nil
}

func (s *APIHanlder) ConnectClient() error {
	var err error
	apiOnce.Do(func() {
		s.client, err = client.NewAgent(s.log)
		if err != nil {
			s.log.Errorf("failed to connect agent internal error", err)
		}
	})

	return err
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
