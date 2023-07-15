package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"github.com/kube-tarian/kad/server/pkg/types"

	"go.uber.org/zap"
)

type Api struct {
	serverpb.UnimplementedServerServer
	agents map[string]*client.Agent
	vault  *client.Vault
}

func New() (*Api, error) {
	vaultClient, err := client.NewVault()
	if err != nil {
		return nil, err
	}

	return &Api{
		agents: make(map[string]*client.Agent),
		vault:  vaultClient,
	}, nil
}

func (a *Api) NewClusterRegistration(ctx context.Context, request *serverpb.NewClusterRegistrationRequest) (
	*serverpb.NewClusterRegistrationResponse, error) {
	logger := log.GetLogger()
	defer logger.Sync()

	cfg := config.GetConfig()
	orgId, ok := ctx.Value("organizationID").(string)
	if !ok || orgId == "" {
		logger.Error("organizationID is missing")
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "missing orgID",
		}, fmt.Errorf("organizationID is missing")
	}

	agentConfig := &types.AgentConfiguration{
		Address:    request.AgentEndpoint,
		Port:       cfg.GetInt(types.AgentPortCfgKey),
		CaCert:     request.ClientCAChainData,
		Cert:       request.ClientKeyData,
		Key:        request.ClientCertData,
		TlsEnabled: cfg.GetBool(types.AgentTlsEnabledCfgKey),
	}

	if err := a.ConnectClient(orgId, request.ClusterName, agentConfig); err != nil {
		logger.Error("failed to connect agent", zap.Error(err),
			zap.String("orgId", orgId),
			zap.String("cluster-name", request.ClusterName))
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed connect cluster",
		}, err
	}

	session, err := db.New(cfg.GetString(types.ServerDbCfgKey))
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed register cluster",
		}, err
	}

	err = session.RegisterCluster(orgId, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed register cluster",
		}, err
	}

	err = a.vault.PutCert(ctx, orgId, request.ClusterName,
		request.ClientCAChainData, request.ClientKeyData, request.ClientCertData)
	if err != nil {
		logger.Error("failed to store cert in vault", zap.Error(err))
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed register cluster",
		}, err
	}

	return &serverpb.NewClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "registered cluster successfully",
	}, nil
}

func (a *Api) UpdateClusterRegistration(ctx context.Context, request *serverpb.UpdateClusterRegistrationRequest) (
	*serverpb.UpdateClusterRegistrationResponse, error) {

	logger := log.GetLogger()
	defer logger.Sync()

	cfg := config.GetConfig()
	orgId, ok := ctx.Value("organizationID").(string)
	if !ok || orgId == "" {
		logger.Error("organizationID is missing")
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "missing orgID",
		}, fmt.Errorf("organizationID is missing")
	}

	agentConfig := &types.AgentConfiguration{
		Address:    request.AgentEndpoint,
		Port:       cfg.GetInt(types.AgentPortCfgKey),
		CaCert:     request.ClientCAChainData,
		Cert:       request.ClientKeyData,
		Key:        request.ClientCertData,
		TlsEnabled: cfg.GetBool(types.AgentTlsEnabledCfgKey),
	}

	if err := a.ReConnect(orgId, request.ClusterName, agentConfig); err != nil {
		logger.Error("failed to connect agent", zap.Error(err),
			zap.String("orgId", orgId),
			zap.String("cluster-name", request.ClusterName))
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed connect cluster",
		}, err
	}

	session, err := db.New(cfg.GetString(types.ServerDbCfgKey))
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update cluster info",
		}, err
	}

	err = session.UpdateCluster(orgId, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update cluster info",
		}, err
	}

	err = a.vault.PutCert(ctx, orgId, request.ClusterName,
		request.ClientCAChainData, request.ClientKeyData, request.ClientCertData)
	if err != nil {
		logger.Error("failed to store cert in vault", zap.Error(err))
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update cluster info",
		}, err
	}

	return &serverpb.UpdateClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster info update success",
	}, nil
}

func (a *Api) DeleteClusterRegistration(ctx context.Context, request *serverpb.DeleteClusterRegistrationRequest) (
	*serverpb.DeleteClusterRegistrationResponse, error) {
	cfg := config.GetConfig()
	logger := log.GetLogger()
	defer logger.Sync()

	orgId, ok := ctx.Value("organizationID").(string)
	if !ok || orgId == "" {
		logger.Error("organizationID is missing")
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "missing orgID",
		}, fmt.Errorf("organizationID is missing")
	}

	a.Close(getClusterKey(orgId, request.ClusterName))
	session, err := db.New(cfg.GetString(types.ServerDbCfgKey))
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete cluster info",
		}, err
	}

	err = session.DeleteCluster(orgId, request.ClusterName)
	if err != nil {
		logger.Error("failed to delete cluster", zap.Error(err))
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to delete cluster info",
		}, err
	}
	return &serverpb.DeleteClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster deletion success",
	}, nil
}

func (a *Api) GetClusters(ctx context.Context, request *serverpb.GetClustersRequest) (
	*serverpb.GetClustersResponse, error) {
	cfg := config.GetConfig()
	logger := log.GetLogger()
	defer logger.Sync()

	orgId, ok := ctx.Value("organizationID").(string)
	if !ok || orgId == "" {
		logger.Error("organizationID is missing")
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "missing orgID",
		}, fmt.Errorf("organizationID is missing")
	}

	session, err := db.New(cfg.GetString(types.ServerDbCfgKey))
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	clusterDetails, err := session.GetClusters(orgId)
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	var data []*serverpb.ClusterInfo
	for _, cluster := range clusterDetails {
		data = append(data, &serverpb.ClusterInfo{
			ClusterName:   cluster.ClusterName,
			AgentEndpoint: cluster.Endpoint,
		})
	}

	return &serverpb.GetClustersResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "retrieved cluster details successfully",
		Data:          data,
	}, nil
}

func (a *Api) GetClusterApps(ctx context.Context, request *serverpb.GetClusterAppsRequest) (
	*serverpb.GetClusterAppsResponse, error) {

	return &serverpb.GetClusterAppsResponse{}, nil
}

func (a *Api) GetClusterAppConfig(ctx context.Context, request *serverpb.GetClusterAppConfigRequest) (
	*serverpb.GetClusterAppConfigResponse, error) {

	return &serverpb.GetClusterAppConfigResponse{}, nil
}
