package agent

import (
	"context"

	"github.com/gogo/status"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"google.golang.org/grpc/codes"
)

func (a *Agent) SyncApp(ctx context.Context, request *agentpb.SyncAppRequest) (*agentpb.SyncAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncApp not implemented")
}
func (a *Agent) GetClusterApps(ctx context.Context, request *agentpb.GetClusterAppsRequest) (*agentpb.GetClusterAppsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterApps not implemented")
}
func (a *Agent) GetClusterAppLaunches(ctx context.Context, request *agentpb.GetClusterAppLaunchesRequest) (*agentpb.GetClusterAppLaunchesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterAppLaunches not implemented")
}
func (a *Agent) GetClusterAppConfig(ctx context.Context, request *agentpb.GetClusterAppConfigRequest) (*agentpb.GetClusterAppConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterAppConfig not implemented")
}
func (a *Agent) GetClusterAppValues(ctx context.Context, request *agentpb.GetClusterAppValuesRequest) (*agentpb.GetClusterAppValuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterAppValues not implemented")
}
