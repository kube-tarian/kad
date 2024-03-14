package activities

import (
	"context"

	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	pluginappstore "github.com/kube-tarian/kad/capten/common-pkg/pluginapp-store"
	dbstore "github.com/kube-tarian/kad/capten/deployment-worker/internal/db-store"
)

type PluginActivities struct {
	store *dbstore.Store
	pas   *pluginappstore.Store
}

func NewPluginActivities() (*PluginActivities, error) {
	store, err := dbstore.NewStore(logger)
	if err != nil {
		return nil, err
	}

	pas, err := pluginappstore.NewStore(logger)
	if err != nil {
		logger.Errorf("failed to initialize plugin app store, %v", err)
	}

	return &PluginActivities{
		store: store,
		pas:   pas,
	}, nil
}

func (p *PluginActivities) PluginDeployActivity(ctx context.Context, req *clusterpluginspb.DeployClusterPluginRequest) (*clusterpluginspb.DeployClusterPluginResponse, error) {
	return nil, nil
}

func (p *PluginActivities) PluginUndeployActivity(ctx context.Context, req *clusterpluginspb.DeployClusterPluginRequest) (*clusterpluginspb.DeployClusterPluginResponse, error) {
	return nil, nil
}

func (p *PluginActivities) PluginDeployPreActionPostgresStoreActivity(ctx context.Context, req *clusterpluginspb.DeployClusterPluginRequest) (*clusterpluginspb.DeployClusterPluginResponse, error) {
	return nil, nil
}

func (p *PluginActivities) PluginDeployPreActionVaultStoreActivity(ctx context.Context, req *clusterpluginspb.DeployClusterPluginRequest) (*clusterpluginspb.DeployClusterPluginResponse, error) {
	return nil, nil
}

func (p *PluginActivities) PluginDeployPreActionMTLSActivity(ctx context.Context, req *clusterpluginspb.DeployClusterPluginRequest) (*clusterpluginspb.DeployClusterPluginResponse, error) {
	return nil, nil
}

func (p *PluginActivities) PluginDeployPostActionActivity(ctx context.Context, req *clusterpluginspb.DeployClusterPluginRequest) (*clusterpluginspb.DeployClusterPluginResponse, error) {
	return nil, nil
}
