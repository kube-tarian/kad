package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/logging"
)

type APIHanlder struct {
	log logging.Logger
}

func NewAPIHandler(log logging.Logger) (*APIHanlder, error) {
	return &APIHanlder{
		log: log,
	}, nil
}

func (s *APIHanlder) Close(c *gin.Context) {
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
