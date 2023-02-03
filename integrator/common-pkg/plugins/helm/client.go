package helm

import (
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
)

type HelmCLient struct {
	logger logging.Logger
}

func NewClient(logger logging.Logger) (*HelmCLient, error) {
	return &HelmCLient{logger: logger}, nil
}
