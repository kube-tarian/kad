package ginapiserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/gin-api-server/api"
	"github.com/kube-tarian/kad/capten/agent/internal/config"
)

var log = logging.NewLogger()

func StartRestServer(rpcapi api.ServerInterface, cfg *config.SericeConfig) {
	r := gin.Default()
	api.RegisterHandlers(r, rpcapi)

	r.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.RestPort))
}
