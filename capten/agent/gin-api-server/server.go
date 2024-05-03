package ginapiserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/capten/agent/gin-api-server/api"
	"github.com/kube-tarian/kad/capten/agent/internal/config"
)

func StartRestServer(rpcapi api.ServerInterface, cfg *config.SericeConfig, certFileName, keyFileName string) {
	r := gin.Default()
	api.RegisterHandlers(r, rpcapi)

	r.RunTLS(fmt.Sprintf("%s:%d", cfg.Host, cfg.RestPort), certFileName, keyFileName)
}
