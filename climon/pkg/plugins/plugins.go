package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/climon/pkg/plugins/helm"
)

type Plugin interface {
	Run(ctx context.Context, payload interface{}) error
	Status() string
}

func GetPlugin(plugin string) (Plugin, error) {
	switch strings.ToLower(plugin) {
	case "helm":
		return helm.NewHelm(), nil
	}

	return nil, fmt.Errorf("plugin %s not found", plugin)
}
