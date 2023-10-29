package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/logging"
	agentmodel "github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/model"
)

type Activities struct {
	config Config
	hg     HandleGit
}

func NewActivity() (*Activities, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	handleGit, err := NewHandleGit(config)
	if err != nil {
		return nil, err
	}

	return &Activities{config: config, hg: handleGit}, nil
}

var logger = logging.NewLogger()

func (a *Activities) ConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", params.Resource)

	switch params.Resource {
	case "cluster":
		return handleCluster(ctx, params, payload)
	case "repository":
		return handleRepository(ctx, params, payload)
	case "project":
		return handleProject(ctx, params, payload)
	case Tekton, CrossPlane:
		return a.hg.handleGit(ctx, params, payload)
	default:
		logger.Errorf("unknown resource type: %s in configuration", params.Resource)
		return model.ResponsePayload{
			Status:  string(agentmodel.WorkFlowStatusFailed),
			Message: json.RawMessage("{\"error\": \"unknown resource type in configuration\"}"),
		}, fmt.Errorf("unknown resource type in configuration")
	}
}
