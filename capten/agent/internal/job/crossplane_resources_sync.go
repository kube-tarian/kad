package job

import (
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/internal/crossplane"
	captenstore "github.com/kube-tarian/kad/capten/common-pkg/capten-store"
)

type CrossplaneResourcesSync struct {
	dbStore         *captenstore.Store
	clusterHandler  *crossplane.ClusterClaimSyncHandler
	providerHandler *crossplane.ProvidersSyncHandler
	log             logging.Logger
	frequency       string
}

func NewCrossplaneResourcesSync(log logging.Logger, frequency string, dbStore *captenstore.Store) (*CrossplaneResourcesSync, error) {
	ccObj, err := crossplane.NewClusterClaimSyncHandler(log, dbStore)
	if err != nil {
		return nil, err
	}
	providerObj, err := crossplane.NewProvidersSyncHandler(log, dbStore)
	if err != nil {
		return nil, err
	}
	return &CrossplaneResourcesSync{
		log:             log,
		frequency:       frequency,
		dbStore:         dbStore,
		clusterHandler:  ccObj,
		providerHandler: providerObj,
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
