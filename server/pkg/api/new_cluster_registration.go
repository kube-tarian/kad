package api

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) NewClusterRegistration(ctx context.Context, request *serverpb.NewClusterRegistrationRequest) (
	*serverpb.NewClusterRegistrationResponse, error) {
	orgId, err := validateRequest(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	clusterID := gocql.TimeUUID().String()
	s.log.Infof("New cluster registration request for cluster %s recieved, clusterId: %s, [org: %s]",
		request.ClusterName, clusterID, orgId)

	caData, caDataErr := getBase64DecodedString(request.ClientCAChainData)
	clientKey, clientKeyErr := getBase64DecodedString(request.ClientKeyData)
	clientCrt, clientCrtErr := getBase64DecodedString(request.ClientCertData)
	if caDataErr != nil || clientKeyErr != nil || clientCrtErr != nil {
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "only base64 encoded certificates are allowed",
		}, nil
	}

	agentConfig := &agent.Config{
		ClusterName: request.ClusterName,
		Address:     request.AgentEndpoint,
		CaCert:      caData,
		Key:         clientKey,
		Cert:        clientCrt,
	}
	if err := s.agentHandeler.AddAgent(clusterID, agentConfig); err != nil {
		s.log.Errorf("failed to connect to agent on cluster %s, %v", request.ClusterName, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent",
		}, nil
	}

	err = credential.PutClusterCerts(ctx, clusterID, caData, clientKey, clientCrt)
	if err != nil {
		s.log.Errorf("failed to store cert in vault for cluster %s, %v", clusterID, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to store cluster credentials",
		}, nil
	}

	err = s.serverStore.AddCluster(orgId, clusterID, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		s.log.Errorf("failed to store cluster %s to db, %v", request.ClusterName, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to store cluster registration",
		}, nil
	}

	if s.cfg.RegisterLaunchAppsConifg {
		if err := s.configureSSOForClusterApps(ctx, orgId, clusterID); err != nil {
			s.log.Errorf("%v", err)
			return &serverpb.NewClusterRegistrationResponse{
				Status:        serverpb.StatusCode_INTERNRAL_ERROR,
				StatusMessage: "failed to configure SSO for cluster apps",
			}, nil
		}
	}

	s.log.Infof("New cluster registration request for cluster %s successful, clusterId: %s, [org: %s]",
		request.ClusterName, clusterID, orgId)
	return &serverpb.NewClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "register cluster success",
		ClusterID:     clusterID,
	}, nil
}
