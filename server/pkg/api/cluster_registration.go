package api

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) NewClusterRegistration(ctx context.Context, request *serverpb.NewClusterRegistrationRequest) (
	*serverpb.NewClusterRegistrationResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[org: %s] New cluster registration request for cluster %s recieved", orgId, request.ClusterName)
	clusterID := gocql.TimeUUID().String()

	caData, caDataErr := s.getBase64DecodedString(request.ClientCAChainData)
	clientKey, clientKeyErr := s.getBase64DecodedString(request.ClientKeyData)
	clientCrt, clientCrtErr := s.getBase64DecodedString(request.ClientCertData)
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
		s.log.Errorf("[org: %s] failed to connect to agent on cluster %s, %v", orgId, request.ClusterName, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent",
		}, nil
	}

	err := credential.PutClusterCerts(ctx, clusterID,
		caData, clientKey, clientCrt)
	if err != nil {
		s.log.Errorf("[org: %s] failed to store cert in vault for cluster %s, %v", orgId, clusterID, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to store cluster credentials",
		}, nil
	}

	err = s.serverStore.AddCluster(orgId, clusterID, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		s.log.Errorf("[org: %s] failed to store cluster %s to db, %v", orgId, request.ClusterName, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to store cluster registration",
		}, nil
	}

	if s.cfg.RegisterLaunchAppsConifg {
		if err := s.configureSSOForClusterApps(ctx, orgId, clusterID); err != nil {
			s.log.Errorf("[org: %s] %v", orgId, err)
			return &serverpb.NewClusterRegistrationResponse{
				Status:        serverpb.StatusCode_INTERNRAL_ERROR,
				StatusMessage: "failed to configure SSO for cluster apps",
			}, nil
		}
	}

	s.log.Infof("[org: %s] New cluster registration successful for %s/%s cluster", orgId, request.ClusterName, clusterID)
	return &serverpb.NewClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "register cluster success",
		ClusterID:     clusterID,
	}, nil
}

func (s *Server) UpdateClusterRegistration(ctx context.Context, request *serverpb.UpdateClusterRegistrationRequest) (
	*serverpb.UpdateClusterRegistrationResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[org: %s] Update cluster registration request for cluster %s recieved", orgId, request.ClusterName)

	caData, caDataErr := s.getBase64DecodedString(request.ClientCAChainData)
	clientKey, clientKeyErr := s.getBase64DecodedString(request.ClientKeyData)
	clientCrt, clientCrtErr := s.getBase64DecodedString(request.ClientCertData)
	if caDataErr != nil || clientKeyErr != nil || clientCrtErr != nil {
		return &serverpb.UpdateClusterRegistrationResponse{
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

	if err := s.agentHandeler.UpdateAgent(request.ClusterID, agentConfig); err != nil {
		s.log.Errorf("[org: %s] failed to connect to agent on cluster %s, %v", orgId, request.ClusterName, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent",
		}, nil
	}

	err := credential.PutClusterCerts(ctx, request.ClusterID,
		caData, clientKey, clientCrt)
	if err != nil {
		s.log.Errorf("[org: %s] failed to update cert in vault for cluster %s, %v", orgId, request.ClusterID, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update register cluster",
		}, nil
	}

	err = s.serverStore.UpdateCluster(orgId, request.ClusterID, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		s.log.Errorf("[org: %s] failed to update cluster %s in db, %v", orgId, request.ClusterName, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update register cluster",
		}, nil
	}

	s.log.Infof("[org: %s] Update cluster registration successful for %s/%s cluster", orgId, request.ClusterName, request.ClusterID)
	return &serverpb.UpdateClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster register update success",
	}, nil
}

func (s *Server) DeleteClusterRegistration(ctx context.Context, request *serverpb.DeleteClusterRegistrationRequest) (
	*serverpb.DeleteClusterRegistrationResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[org: %s] Delete cluster registration request for cluster %s recieved", orgId, request.ClusterID)
	s.agentHandeler.RemoveAgent(request.ClusterID)
	err := credential.DeleteClusterCerts(ctx, request.ClusterID)
	if err != nil {
		s.log.Errorf("[org: %s] failed to delete cert in vault for cluster %s, %v", orgId, request.ClusterID, err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	err = s.serverStore.DeleteFullCacheAppLaunches(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("[org: %s] failed to delete clusterappLaunches %s from db, %v", orgId, request.ClusterID, err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	err = s.serverStore.DeleteCluster(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("[org: %s] failed to delete cluster %s from db, %v", orgId, request.ClusterID, err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	s.log.Infof("[org: %s] Delete cluster registration request for cluster %s successful", orgId, request.ClusterID)
	return &serverpb.DeleteClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster deletion success",
	}, nil
}

func (s *Server) GetClusters(ctx context.Context, request *serverpb.GetClustersRequest) (
	*serverpb.GetClustersResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[org: %s] GetClusters request recieved", orgId)
	clusterDetails, err := s.serverStore.GetClusters(orgId)
	if err != nil {
		s.log.Errorf("[org: %s] failed to get clusters, %v", orgId, err)
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	var data []*serverpb.ClusterInfo
	for _, cluster := range clusterDetails {
		resp, err := s.GetClusterAppLaunchesFromCacheOrAgent(ctx, orgId, cluster.ClusterID)
		if err != nil {
			s.log.Errorf("failed to get cluster appLaunches for cluster: %s, %v", cluster.ClusterID, err)
			continue
		}

		attributes := []*serverpb.ClusterAttribute{}
		data = append(data, &serverpb.ClusterInfo{
			ClusterID:        cluster.ClusterID,
			ClusterName:      cluster.ClusterName,
			AgentEndpoint:    cluster.Endpoint,
			Attributes:       attributes,
			AppLaunchConfigs: mapAgentAppLauncesToServerResp(resp.LaunchConfigList),
		})
	}

	s.log.Infof("[org: %s] Found %d clusters", orgId, len(data))
	return &serverpb.GetClustersResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "get cluster details success",
		Data:          data,
	}, nil
}

func (s *Server) GetCluster(ctx context.Context, request *serverpb.GetClusterRequest) (
	*serverpb.GetClusterResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.GetClusterResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[org: %s] GetCluster request recieved for cluster %s", orgId, request.ClusterID)
	clusterDetails, err := s.serverStore.GetClusterDetails(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("[org: %s] failed to get cluster %s, %v", orgId, request.ClusterID, err)
		return &serverpb.GetClusterResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	resp, err := s.GetClusterAppLaunchesFromCacheOrAgent(ctx, orgId, request.ClusterID)
	if err != nil || resp == nil || resp.Status != agentpb.StatusCode_OK {
		s.log.Error("failed to get cluster application launches from cache/agent: %v", resp)
		return &serverpb.GetClusterResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster app lauches",
		}, err
	}

	attributes := []*serverpb.ClusterAttribute{}
	data := &serverpb.ClusterInfo{
		ClusterID:        request.ClusterID,
		ClusterName:      clusterDetails.ClusterName,
		AgentEndpoint:    clusterDetails.Endpoint,
		Attributes:       attributes,
		AppLaunchConfigs: mapAgentAppLauncesToServerResp(resp.LaunchConfigList),
	}

	s.log.Infof("[org: %s] GetCluster request processed for cluster %s", orgId, request.ClusterID)
	return &serverpb.GetClusterResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "get cluster details success",
		Data:          data,
	}, nil
}

func (s *Server) getBase64DecodedString(encodedString string) (string, error) {
	decodedByte, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		// This will assume the string is not encoded and returns the original string.
		s.log.Errorf("Failed to decode the string: %v", err)
		return "", err
	}

	return string(decodedByte), nil
}

func (s *Server) GetClusterAppLaunchesFromCacheOrAgent(ctx context.Context, orgId, clusterID string) (
	*agentpb.GetClusterAppLaunchesResponse, error) {
	currentTime := time.Now()
	lastFetchedTime := time.Unix(s.orgClusterIDCache[orgId+"-"+clusterID], 0)

	if currentTime.After(lastFetchedTime) {
		// cache expired re-trigger the cache
		agentClient, aErr := s.agentHandeler.GetAgent(orgId, clusterID)
		if aErr == nil {
			resp, err := agentClient.GetClient().GetClusterAppLaunches(ctx, &agentpb.GetClusterAppLaunchesRequest{})
			if err == nil {
				updateErr := s.serverStore.UpdateCacheAppLaunches(orgId, clusterID, resp.LaunchConfigList)
				if updateErr == nil {
					s.mutex.Lock()
					s.orgClusterIDCache[orgId+"-"+clusterID] = currentTime.Add(delayTimeinMin * time.Minute).Unix()
					s.mutex.Unlock()

					return resp, err
				}

			}
		}
	}

	// If any failure happens return from cache.
	return s.serverStore.GetCacheAppLaunches(orgId, clusterID)
}
