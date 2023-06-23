package helm

import (
	"github.com/intelops/go-common/logging"
)

type HelmCLient struct {
	logger logging.Logger
}

func NewClient(logger logging.Logger) (*HelmCLient, error) {
	return &HelmCLient{logger: logger}, nil
}
