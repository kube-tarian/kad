package ginapiserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/gin-api-server/api"
	"github.com/kube-tarian/kad/capten/agent/internal/clusterissuer"
	"github.com/kube-tarian/kad/capten/agent/internal/config"
)

func StartRestServer(rpcapi api.ServerInterface, cfg *config.SericeConfig, log logging.Logger) error {
	err := clusterissuer.GenerateServerCertificates(cfg.ClusterCAIssuerName, log)
	if err != nil {
		log.Errorf("Failed to generate Server certificate, %v", err)
		return err
	}

	r := gin.Default()
	api.RegisterHandlers(r, rpcapi)

	return r.RunTLS(fmt.Sprintf("%s:%d", cfg.Host, cfg.RestPort), clusterissuer.CertFileName, clusterissuer.KeyFileName)
}
