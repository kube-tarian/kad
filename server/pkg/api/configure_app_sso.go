package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/pkg/errors"
)

func (s *Server) configureSSOForClusterApps(ctx context.Context, orgId, clusterID string) error {
	agentClient, err := s.agentHandeler.GetAgent(orgId, clusterID)
	if err != nil {
		return errors.WithMessagef(err, "failed to get agent for cluster %s", clusterID)
	}

	resp, err := agentClient.GetClient().GetClusterAppLaunches(ctx, &agentpb.GetClusterAppLaunchesRequest{})
	if err != nil || resp == nil || resp.Status != agentpb.StatusCode_OK {
		return fmt.Errorf("failed to get cluster app launches from cluster %s, err: %v", clusterID, resp)
	}

	for _, app := range resp.LaunchConfigList {
		err := s.configureSSOForClusterApp(ctx, agentClient, clusterID, app.ReleaseName, app.LaunchURL)
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *Server) configureSSOForClusterApp(ctx context.Context, agentClient *agent.Agent, clusterID, releaseName, launchURL string) error {
	s.log.Infof("Configuring app launch SSO for %s on cluster %s", releaseName, clusterID)
	appName := fmt.Sprintf("%s-%s", clusterID, releaseName)
	clientID, clientSecret, err := s.iam.RegisterAppClientSecrets(ctx, appName, launchURL)
	if err != nil {
		return errors.WithMessagef(err, "failed to register app %s on cluster %s with IAM", releaseName, clusterID)
	}

	ssoResp, err := agentClient.GetClient().ConfigureAppSSO(ctx, &agentpb.ConfigureAppSSORequest{
		ReleaseName:  releaseName,
		ClientId:     clientID,
		ClientSecret: clientSecret,
		OAuthBaseURL: s.cfg.CaptenOAuthURL,
	})

	if err != nil || ssoResp == nil || ssoResp.Status != agentpb.StatusCode_OK {
		return fmt.Errorf("failed to configure sso for app  %s on cluster %s, err: %v", releaseName, clusterID, ssoResp)
	}
	s.log.Infof("Configured app launch SSO for %s on cluster %s", releaseName, clusterID)
	return nil
}
