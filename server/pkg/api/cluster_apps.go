package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func mapAgentAppsToServerResp(appDataList []*agentpb.AppData) []*serverpb.ClusterAppConfig {
	clusterAppConfigs := make([]*serverpb.ClusterAppConfig, len(appDataList))
	for index, appConfig := range appDataList {
		var clusterAppConfig serverpb.ClusterAppConfig
		clusterAppConfig.AppName = appConfig.Config.AppName
		clusterAppConfig.Version = appConfig.Config.Version
		clusterAppConfig.Category = appConfig.Config.Category
		clusterAppConfig.Description = appConfig.Config.Description
		clusterAppConfig.ChartName = appConfig.Config.ChartName
		clusterAppConfig.RepoName = appConfig.Config.RepoName
		clusterAppConfig.RepoURL = appConfig.Config.RepoURL
		clusterAppConfig.Namespace = appConfig.Config.Namespace
		clusterAppConfig.CreateNamespace = appConfig.Config.CreateNamespace
		clusterAppConfig.PrivilegedNamespace = appConfig.Config.PrivilegedNamespace
		clusterAppConfig.Icon = appConfig.Config.Icon
		clusterAppConfig.LaunchURL = appConfig.Config.LaunchURL
		clusterAppConfig.InstallStatus = appConfig.Config.InstallStatus
		clusterAppConfig.RuntimeStatus = ""

		clusterAppConfigs[index] = &clusterAppConfig
	}

	return clusterAppConfigs

}

func mapAgentAppLauncesToServerResp(appLaunchCfgs []*agentpb.AppLaunchConfig) []*serverpb.AppLaunchConfig {
	svrAppLaunchCfg := make([]*serverpb.AppLaunchConfig, len(appLaunchCfgs))

	for index, cfg := range appLaunchCfgs {
		var launchCfg serverpb.AppLaunchConfig
		launchCfg.ReleaseName = cfg.ReleaseName
		launchCfg.Category = cfg.Category
		launchCfg.LaunchUIDescription = cfg.Description
		launchCfg.Icon = cfg.Icon
		launchCfg.LaunchURL = cfg.LaunchURL

		svrAppLaunchCfg[index] = &launchCfg
	}

	return svrAppLaunchCfg
}

func (s *Server) GetClusterApps(ctx context.Context, request *serverpb.GetClusterAppsRequest) (
	*serverpb.GetClusterAppsResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 || request.ClusterID == "" {
		s.log.Error("organizationID/ClusterID is missing in the request")
		return &serverpb.GetClusterAppsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}
	s.log.Infof("[org: %s] GetClusterApps request recieved for cluster %s", orgId, request.ClusterID)

	a, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		s.log.Error("failed to connect to agent", err)
		return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent"}, nil
	}

	resp, err := a.GetClient().GetClusterApps(ctx, &agentpb.GetClusterAppsRequest{})
	if err != nil {
		s.log.Error("failed to get cluster application from agent", err)
		return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get cluster application from agent"}, nil
	}

	s.log.Infof("[org: %s] Fetched %d installed apps from the cluster %s", orgId, len(resp.AppData), request.ClusterID)
	return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_OK,
		StatusMessage: "successfully fetched the data from agent",
		AppConfigs:    mapAgentAppsToServerResp(resp.AppData)}, nil
}

func (s *Server) GetClusterAppLaunchConfigs(ctx context.Context, request *serverpb.GetClusterAppLaunchConfigsRequest) (
	*serverpb.GetClusterAppLaunchConfigsResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 || request.ClusterID == "" {
		s.log.Error("organizationID/ClusterID is missing in the request")
		return &serverpb.GetClusterAppLaunchConfigsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[org: %s] GetClusterAppLaunchConfigs request recieved for cluster %s", orgId, request.ClusterID)
	a, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		s.log.Error("failed to connect to agent", err)
		return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent"}, nil
	}

	resp, err := a.GetClient().GetClusterAppLaunches(ctx, &agentpb.GetClusterAppLaunchesRequest{})
	if err != nil {
		s.log.Error("failed to get cluster application launches from agent", err)
		return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get cluster application launches from agent"}, err
	}

	s.log.Infof("[org: %s] Fetched %d app launch UIs from the cluster %s", orgId, len(resp.LaunchConfigList), request.ClusterID)
	return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_OK,
		StatusMessage:   "successfully fetched the data from agent",
		AppLaunchConfig: mapAgentAppLauncesToServerResp(resp.LaunchConfigList)}, nil
}

func (s *Server) GetClusterApp(ctx context.Context, request *serverpb.GetClusterAppRequest) (
	*serverpb.GetClusterAppResponse, error) {
	return &serverpb.GetClusterAppResponse{}, nil
}

func (s *Server) configureSSOForApp(ctx context.Context, agentClient agentpb.AgentClient, app *agentpb.AppLaunchConfig) error {
	iamURL := s.iam.GetURL()

	if app.LaunchURL == "" {
		return nil
	}

	clientID, clientSecret, err := s.iam.RegisterAppClientSecrets(ctx, app.ReleaseName, app.LaunchURL)
	if err != nil {
		return err
	}

	// assume the agent store the creds as part of this.
	ssoResp, err := agentClient.ConfigureAppSSO(ctx, &agentpb.ConfigureAppSSORequest{
		ReleaseName:  app.ReleaseName,
		ClientId:     clientID,
		ClientSecret: clientSecret,
		OAuthBaseURL: iamURL,
	})

	if err != nil || ssoResp.Status != 0 {
		return err
	}

	return nil
}
