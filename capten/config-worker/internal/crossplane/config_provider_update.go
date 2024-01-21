package crossplane

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

const (
	CrossPlaneResource  = "crossplane"
	CrossplaneNamespace = "crossplane-system"
)

func (cp *CrossPlaneApp) configureConfigProviderUpdate(ctx context.Context, req *model.CrossplaneProviderUpdate) (status string, err error) {
	logger.Infof("configuring config provider %s update", req.CloudType)

	fmt.Println("provider")
	x, _ := json.Marshal(req)
	fmt.Println(string(x))
	// x := strings.TrimPrefix(req.ManagedClusterName, "provider-")

	err = cp.helper.SyncArgoCDApp(ctx, CrossPlaneResource, CrossPlaneResource)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to sync config providers")
	}

	logger.Infof("synched config providers %s", CrossPlaneResource)

	err = cp.helper.WaitForArgoCDToSync(ctx, CrossPlaneResource, CrossPlaneResource)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to fetch config providers")
	}

	return string(agentmodel.WorkFlowStatusCompleted), nil
}
