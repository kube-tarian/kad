package ginapiserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/gin-api-server/api"
	"github.com/kube-tarian/kad/capten/agent/internal/config"
)

func StartRestServer(rpcapi api.ServerInterface, cfg *config.SericeConfig, log logging.Logger) error {
	r := gin.Default()
	api.RegisterHandlers(r, rpcapi)

	return r.RunTLS(fmt.Sprintf("%s:%d", cfg.Host, cfg.RestPort), cfg.CertFileName, cfg.KeyFileName)
}
