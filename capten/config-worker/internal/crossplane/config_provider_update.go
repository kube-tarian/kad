package crossplane

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/capten/model"
)

func (cp *CrossPlaneApp) configureConfigProviderUpdate(ctx context.Context, req *model.CrossplaneClusterUpdate) (status string, err error) {
	logger.Infof("configuring config provider %s update", req.ManagedClusterName)

	x, _ := json.Marshal(req)

	fmt.Println("configureConfigProviderUpdate request")

	fmt.Println(string(x))

	return "", nil

}
