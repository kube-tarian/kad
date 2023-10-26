package api

import (
	"context"
	"encoding/hex"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"gopkg.in/yaml.v2"
)

func (s *Server) DeployStoreApp(ctx context.Context, request *serverpb.DeployStoreAppRequest) (
	*serverpb.DeployStoreAppResponse, error) {
	orgId, err := validateOrgWithArgs(ctx, request.ClusterID, request.AppName, request.Version)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Deploy store app [%s:%s] request for cluster %s recieved, [org: %s]",
		request.AppName, request.Version, request.ClusterID, orgId)

	config, err := s.serverStore.GetAppFromStore(request.AppName, request.Version)
	if err != nil {
		s.log.Errorf("failed to get store app values, %v", err)
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to find store app values",
		}, nil
	}

	marshaledOverride, err := yaml.Marshal(config.OverrideValues)
	if err != nil {
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to find store app values",
		}, nil
	}

	marshaledLaunchUi, err := yaml.Marshal(config.LaunchUIValues)
	if err != nil {
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to find store app values",
		}, nil
	}

	decodedIconBytes, _ := hex.DecodeString(config.Icon)
	req := &agentpb.InstallAppRequest{
		AppConfig: &agentpb.AppConfig{
			AppName:             config.Name,
			Version:             config.Version,
			ReleaseName:         config.ReleaseName,
			Category:            config.Category,
			Description:         config.Description,
			ChartName:           config.ChartName,
			RepoName:            config.RepoName,
			RepoURL:             config.RepoURL,
			Namespace:           config.Namespace,
			CreateNamespace:     config.CreateNamespace,
			PrivilegedNamespace: config.PrivilegedNamespace,
			Icon:                decodedIconBytes,
			LaunchURL:           config.LaunchURL,
			LaunchUIDescription: config.LaunchUIDescription,
			DefualtApp:          false,
			PluginName:          config.PluginName,
			PluginDescription:   config.PluginDescription,
		},
		AppValues: &agentpb.AppValues{
			OverrideValues: request.OverrideValues,
			LaunchUIValues: decodeBase64StringToBytes(string(marshaledLaunchUi)),
			TemplateValues: decodeBase64StringToBytes(string(marshaledOverride)),
		},
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to deploy the app",
		}, nil
	}

	_, err = agent.GetClient().InstallApp(ctx, req)
	if err != nil {
		s.log.Errorf("failed to deploy app, %v", err)
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to deploy the app",
		}, nil
	}

	s.log.Infof("Deploy Store app [%s:%s] request request triggered for cluster %s, [org: %s]",
		request.AppName, request.Version, request.ClusterID, orgId)

	return &serverpb.DeployStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app is successfully deployed",
	}, nil
}
