package activities

import (
	"context"

	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	dbstore "github.com/kube-tarian/kad/capten/deployment-worker/internal/db-store"
)

type PluginActivities struct {
	store *dbstore.Store
}

func NewPluginActivities() (*PluginActivities, error) {
	store, err := dbstore.NewStore(logger)
	if err != nil {
		return nil, err
	}

	return &PluginActivities{
		store: store,
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

func (p *PluginActivities) PluginDeployInstallActivity(ctx context.Context, req *clusterpluginspb.DeployClusterPluginRequest) (*clusterpluginspb.DeployClusterPluginResponse, error) {
	return nil, nil
}

func (p *PluginActivities) PluginDeployPostActionActivity(ctx context.Context, req *clusterpluginspb.DeployClusterPluginRequest) (*clusterpluginspb.DeployClusterPluginResponse, error) {
	return nil, nil
}
