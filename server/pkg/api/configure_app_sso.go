package api

import (
	"context"
	"fmt"
	"time"

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

	if err := s.serverStore.InsertClusterAppLaunches(orgId, clusterID, resp.LaunchConfigList); err != nil {
		return fmt.Errorf("failed to store cluster app launches on server db %s, err: %v", clusterID, err)
	}

	s.mutex.Lock()
	s.orgClusterIDCache[orgId+"-"+clusterID] = time.Now().Add(delayTimeinMin * time.Minute).Unix()
	s.mutex.Unlock()

	for _, app := range resp.LaunchConfigList {
		appName := fmt.Sprintf("%s-%s", clusterID, app.ReleaseName)
		clientID, clientSecret, err := s.iam.RegisterAppClientSecrets(ctx, appName, app.LaunchURL)
		if err != nil {
			return errors.WithMessagef(err, "failed to register app %s on cluster %s with IAM", app.ReleaseName, clusterID)
		}

		ssoResp, err := agentClient.GetClient().ConfigureAppSSO(ctx, &agentpb.ConfigureAppSSORequest{
			ReleaseName:  app.ReleaseName,
			ClientId:     clientID,
			ClientSecret: clientSecret,
			OAuthBaseURL: s.cfg.CaptenOAuthURL,
		})

		if err != nil || ssoResp == nil || ssoResp.Status != agentpb.StatusCode_OK {
			return fmt.Errorf("failed to configure sso for app  %s on cluster %s, err: %v", app.ReleaseName, clusterID, ssoResp)
		}
	}
	return nil
}
