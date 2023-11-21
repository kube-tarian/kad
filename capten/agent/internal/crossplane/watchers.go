package crossplane

import (
	"fmt"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
)

func RegisterK8SWatcher(log logging.Logger, dbStore *captenstore.Store) error {
	k8sclient, err := k8s.NewK8SClient(log)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	err = RegisterK8SClusterClaimWatcher(log, dbStore, k8sclient.DynamicClientInterface)
	if err != nil {
		return fmt.Errorf("failed to RegisterK8SClusterClaimWatcher: %v", err)
	}

	err = registerK8SProviderWatcher(log, dbStore, k8sclient.DynamicClientInterface)
	if err != nil {
		return fmt.Errorf("failed to RegisterK8SProviderWatcher: %v", err)
	}
	return nil
}
