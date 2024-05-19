package crossplane

import (
	"fmt"
	"time"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/common-pkg/capten-store"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
)

func RegisterK8SWatcher(log logging.Logger, dbStore *captenstore.Store) error {
	k8sclient, err := k8s.NewK8SClient(log)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	clusterHandler, err := NewClusterClaimSyncHandler(log, dbStore)
	if err != nil {
		return err
	}

	provider, err := NewProvidersSyncHandler(log, dbStore)
	if err != nil {
		return err
	}

	go retryForEver(60*time.Second, func() (err error) {
		err = registerK8SClusterClaimWatcher(log, clusterHandler, k8sclient)
		if err != nil {
			return fmt.Errorf("failed to RegisterK8SClusterClaimWatcher: %v", err)
		}
		return nil
	})

	go retryForEver(60*time.Second, func() (err error) {
		err = registerK8SProviderWatcher(log, provider, k8sclient)
		if err != nil {
			return fmt.Errorf("failed to RegisterK8SProviderWatcher: %v", err)
		}
		return nil
	})
	return nil
}

func retryForEver(sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}
		time.Sleep(sleep)
	}
	return
}
