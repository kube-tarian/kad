package helmplugin

import (
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
	"github.com/kube-tarian/kad/integrator/pkg/logging"
)

type HelmCLient struct {
	logger logging.Logger
}

func NewClient(logger logging.Logger) (*HelmCLient, error) {
	return &HelmCLient{logger: logger}, nil
}

func (a *HelmCLient) Exec(payload model.RequestPayload) (json.RawMessage, error) {
	var err error

	switch payload.Action {
	case "install":
		err = a.Install(payload)
	case "delete":
		err = a.Delete(payload)
	case "list":
		err = a.List()
	default:
		err = fmt.Errorf("unsupported action for helm plugin: %v", payload.Action)
	}
	if err != nil {
		a.logger.Errorf("helm %v of application failed, %v", payload.Action, err)
		return nil, err
	}

	return nil, nil
}
