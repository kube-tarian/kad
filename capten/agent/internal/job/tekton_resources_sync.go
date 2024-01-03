package job

import (
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"
	"github.com/kube-tarian/kad/capten/agent/internal/tekton"
)

type TektonResourcesSync struct {
	dbStore       *captenstore.Store
	eventlistener *tekton.TektonPipelineSyncHandler
	log           logging.Logger
	frequency     string
}

func NewTektonResourcesSync(log logging.Logger, frequency string, dbStore *captenstore.Store) (*TektonResourcesSync, error) {
	ccObj := tekton.NewTektonPipelineSyncHandler(log, dbStore)
	return &TektonResourcesSync{
		log:           log,
		frequency:     frequency,
		dbStore:       dbStore,
		eventlistener: ccObj,
	}, nil
}

func (s *TektonResourcesSync) CronSpec() string {
	return s.frequency
}

func (s *TektonResourcesSync) Run() {
	s.log.Debug("started Tekton resource sync job")
	if err := s.eventlistener.Sync(); err != nil {
		s.log.Errorf("failed to synch eventlisteneres, %v", err)
	}
	s.log.Debug("Tekton resource sync job completed")
}
