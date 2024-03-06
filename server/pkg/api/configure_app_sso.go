package api

import (
	"context"
	"fmt"

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
		appClientName := fmt.Sprintf("%s-%s", clusterID, app.ReleaseName)
		s.log.Infof("Register app %s as app-client %s with IAM, clusterId: %s, [org: %s]",
			app.ReleaseName, appClientName, clusterID, orgId)
		clientID, clientSecret, err := s.iam.RegisterAppClientSecrets(ctx, appClientName, app.LaunchURL, orgId)
		if err != nil {
			return errors.WithMessagef(err, "failed to register app %s on cluster %s with IAM", app.ReleaseName, clusterID)
		}

		s.log.Infof("Configuring SSO for app %s, clusterId: %s, [org: %s]", app.ReleaseName, clusterID, orgId)
		ssoResp, err := agentClient.GetClient().ConfigureAppSSO(ctx, &agentpb.ConfigureAppSSORequest{
			ReleaseName:  app.ReleaseName,
			ClientId:     clientID,
			ClientSecret: clientSecret,
			OAuthBaseURL: s.cfg.CaptenOAuthURL,
		})

		if err != nil || ssoResp == nil || ssoResp.Status != agentpb.StatusCode_OK {
			s.log.Errorf("failed to configure sso for app  %s on cluster %s, err: %v", app.ReleaseName, clusterID, err)
			continue
		}
		s.log.Infof("Configure SSO for app %s triggerred, clusterId: %s, [org: %s]",
			app.ReleaseName, appClientName, clusterID, orgId)
	}
	return nil
}
