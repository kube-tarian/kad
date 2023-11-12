package job

import (
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"
	"github.com/kube-tarian/kad/capten/agent/pkg/crossplane"
)

type CrossplaneResourcesSync struct {
	log             logging.Logger
	frequency       string
	clusterHandler  *crossplane.ClusterClaimSyncHandler
	providerHandler *crossplane.ProvidersSyncHandler
}

func NewCrossplaneResourcesSync(log logging.Logger, frequency string, dbStore *captenstore.Store) (*CrossplaneResourcesSync, error) {
	return &CrossplaneResourcesSync{
		log:             log,
		frequency:       frequency,
		clusterHandler:  crossplane.NewClusterClaimSyncHandler(log, dbStore),
		providerHandler: crossplane.NewProvidersSyncHandler(log, dbStore),
	}, nil
}

func (s *CrossplaneResourcesSync) CronSpec() string {
	return s.frequency
}

func (s *CrossplaneResourcesSync) Run() {
	s.log.Debug("started crossplane resource sync job")
	if err := s.providerHandler.Sync(); err != nil {
		s.log.Errorf("failed to synch providers, %v", err)
	}

	if err := s.clusterHandler.Sync(); err != nil {
		s.log.Errorf("failed to synch managed clusters, %v", err)
	}
	s.log.Debug("crossplane resource sync job completed")
}
