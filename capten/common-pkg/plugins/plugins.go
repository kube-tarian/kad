package plugins

import (
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/capten/common-pkg/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/helm"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
)

func GetPlugin(plugin string, logger logging.Logger) (workerframework.Plugin, error) {
	switch strings.ToLower(plugin) {
	case "helm":
		return helm.NewClient(logger)
	case "argocd":
		return argocd.NewClient(logger)
	default:
		return nil, fmt.Errorf("plugin %s not found", plugin)
	}
}
