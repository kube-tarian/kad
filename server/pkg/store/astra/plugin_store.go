package astra

import (
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/types"
)

func (a *AstraServerStore) AddOrUpdatePlugin(config *types.Plugin) error {
	fmt.Println("Plugin added, %+v", config)
	return nil
}
