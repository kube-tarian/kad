package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) GetCaptenPlugins(ctx context.Context, request *captenpluginspb.GetCaptenPluginsRequest) (
	*captenpluginspb.GetCaptenPluginsResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetCaptenPluginsResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get Capten Plugins request for cluster %s recieved, [org: %s]", clusterId, orgId)

	a, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Error("failed to connect to agent", err)
		return &captenpluginspb.GetCaptenPluginsResponse{Status: captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to connect to agent"}, nil
	}

	resp, err := a.GetCaptenPluginsClient().GetCaptenPlugins(ctx, &captenpluginspb.GetCaptenPluginsRequest{})
	if err != nil {
		s.log.Error("failed to get cluster capten plugins from agent", err)
		return &captenpluginspb.GetCaptenPluginsResponse{Status: captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get cluster application from agent"}, nil
	}

	s.log.Infof("Fetched %d capten plugins from the cluster %s, [org: %s]", len(resp.Plugins), clusterId, orgId)
	return &captenpluginspb.GetCaptenPluginsResponse{Status: captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the data from agent",
		Plugins:       resp.Plugins}, nil
}
