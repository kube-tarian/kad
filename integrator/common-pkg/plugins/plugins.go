package plugins

import (
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins/argocd"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins/helmplugin"
	workerframework "github.com/kube-tarian/kad/integrator/common-pkg/worker-framework"
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
