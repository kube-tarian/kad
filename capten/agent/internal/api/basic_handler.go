package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/kube-tarian/kad/capten/agent/gin-api-server/api"
)

// open api swagger documentation
func (a *Agent) GetApiDocs(c *gin.Context) {
	apiDocs, err := api.GetSwagger()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.IndentedJSON(http.StatusOK, apiDocs)
}

// readiness and liveness probe health status endpoint
func (a *Agent) GetStatus(c *gin.Context) {
	c.String(http.StatusOK, "OK\n")
}
