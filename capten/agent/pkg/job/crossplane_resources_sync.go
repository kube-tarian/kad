package job

import (
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"
)

type CrossplaneResourcesSync struct {
	dbStore   *captenstore.Store
	log       logging.Logger
	frequency string
}

func NewCrossplaneResourcesSync(log logging.Logger, frequency string, dbStore *captenstore.Store) (*CrossplaneResourcesSync, error) {
	return &CrossplaneResourcesSync{
		log:       log,
		frequency: frequency,
		dbStore:   dbStore,
	}, nil
}

func (v *CrossplaneResourcesSync) CronSpec() string {
	return v.frequency
}

func (v *CrossplaneResourcesSync) Run() {
	v.log.Debug("started crossplane resource sync job")

	v.log.Debug("crossplane resource sync job completed")
}
