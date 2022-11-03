package plugins

import (
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/plugins/argocd"
	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/plugins/helmplugin"
	"github.com/kube-tarian/kad/integrator/pkg/logging"
	workerframework "github.com/kube-tarian/kad/integrator/pkg/worker-framework"
)

func GetPlugin(plugin string, logger logging.Logger) (workerframework.Plugin, error) {
	switch strings.ToLower(plugin) {
	case "helm":
		return helmplugin.NewClient(logger)
	case "argocd":
		return argocd.NewClient(logger)
	default:
		return nil, fmt.Errorf("plugin %s not found", plugin)
	}
}
