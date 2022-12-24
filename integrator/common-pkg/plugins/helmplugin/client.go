package helmplugin

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/model"
)

type HelmCLient struct {
	logger logging.Logger
}

func NewClient(logger logging.Logger) (*HelmCLient, error) {
	return &HelmCLient{logger: logger}, nil
}

func (a *HelmCLient) DeployActivities(req interface{}) (json.RawMessage, error) {
	var payload model.RequestPayload
	switch p := req.(type) {
	case model.RequestPayload:
		payload = p
	default:
		return nil, fmt.Errorf("unexpected request data, type: %v", reflect.TypeOf(req))
	}

	// payload, ok := req.(model.RequestPayload)
	// if !ok {
	// 	return nil, fmt.Errorf("request is not proper: %v, type: %v", req, reflect.TypeOf(req))
	// }
	switch payload.Action {
	case "install":
		return a.Create(payload)
	case "delete":
		return a.Delete(payload)
	case "list":
		return a.List(payload)
	default:
		return nil, fmt.Errorf("unsupported action for helm plugin: %v", payload.Action)
	}
}
