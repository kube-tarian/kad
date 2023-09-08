package api

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"text/template"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func (s *Server) AddStoreApp(ctx context.Context, request *serverpb.AddStoreAppRequest) (
	*serverpb.AddStoreAppResponse, error) {
	if request.AppConfig.AppName == "" || request.AppConfig.Version == "" {
		s.log.Infof("AppName or version is missing for add store app request")
		return &serverpb.AddStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "AppName or version is missing in the request",
		}, nil
	}

	config := &types.StoreAppConfig{
		ReleaseName:         request.AppConfig.ReleaseName,
		AppName:             request.AppConfig.AppName,
		Version:             request.AppConfig.Version,
		Category:            request.AppConfig.Category,
		Description:         request.AppConfig.Description,
		ChartName:           request.AppConfig.ChartName,
		RepoName:            request.AppConfig.RepoName,
		RepoURL:             request.AppConfig.RepoURL,
		Namespace:           request.AppConfig.Namespace,
		CreateNamespace:     request.AppConfig.CreateNamespace,
		PrivilegedNamespace: request.AppConfig.PrivilegedNamespace,
		Icon:                hex.EncodeToString(request.AppConfig.Icon),
		LaunchURL:           request.AppConfig.LaunchURL,
		LaunchUIDescription: request.AppConfig.LaunchUIDescription,
		OverrideValues:      encodeBase64BytesToString(request.AppValues.OverrideValues),
		LaunchUIValues:      encodeBase64BytesToString(request.AppValues.LaunchUIValues),
		TemplateValues:      encodeBase64BytesToString(request.AppValues.TemplateValues),
	}

	if err := s.serverStore.AddOrUpdateStoreApp(config); err != nil {
		s.log.Errorf("failed to add app config to store, %v", err)
		return &serverpb.AddStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed add app config to store",
		}, nil
	}

	return &serverpb.AddStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly added to store",
	}, nil
}

func (s *Server) UpdateStoreApp(ctx context.Context, request *serverpb.UpdateStoreAppRequest) (
	*serverpb.UpdateStoreAppRsponse, error) {
	if request.AppConfig.AppName == "" || request.AppConfig.Version == "" {
		s.log.Infof("AppName or version is missing for update store app request")
		return &serverpb.UpdateStoreAppRsponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "AppName or version is missing in the request",
		}, nil
	}

	config := &types.StoreAppConfig{
		ReleaseName:         request.AppConfig.ReleaseName,
		AppName:             request.AppConfig.AppName,
		Version:             request.AppConfig.Version,
		Category:            request.AppConfig.Category,
		Description:         request.AppConfig.Description,
		ChartName:           request.AppConfig.ChartName,
		RepoName:            request.AppConfig.RepoName,
		RepoURL:             request.AppConfig.RepoURL,
		Namespace:           request.AppConfig.Namespace,
		CreateNamespace:     request.AppConfig.CreateNamespace,
		PrivilegedNamespace: request.AppConfig.PrivilegedNamespace,
		Icon:                hex.EncodeToString(request.AppConfig.Icon),
		LaunchURL:           request.AppConfig.LaunchURL,
		LaunchUIDescription: request.AppConfig.LaunchUIDescription,
		OverrideValues:      encodeBase64BytesToString(request.AppValues.OverrideValues),
		LaunchUIValues:      encodeBase64BytesToString(request.AppValues.LaunchUIValues),
		TemplateValues:      encodeBase64BytesToString(request.AppValues.TemplateValues),
	}

	if err := s.serverStore.AddOrUpdateStoreApp(config); err != nil {
		s.log.Errorf("failed to update app config in store, %v", err)
		return &serverpb.UpdateStoreAppRsponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update app config in store",
		}, nil
	}

	return &serverpb.UpdateStoreAppRsponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly updated",
	}, nil
}

func (s *Server) DeleteStoreApp(ctx context.Context, request *serverpb.DeleteStoreAppRequest) (
	*serverpb.DeleteStoreAppResponse, error) {
	if request.AppName == "" || request.Version == "" {
		s.log.Infof("AppName or version is missing for delete store app request")
		return &serverpb.DeleteStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "AppName or version is missing in the request",
		}, nil
	}

	if err := s.serverStore.DeleteAppInStore(request.AppName, request.Version); err != nil {
		s.log.Errorf("failed to delete app config from store, %v", err)
		return &serverpb.DeleteStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to delete app config from store",
		}, nil
	}

	return &serverpb.DeleteStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly deleted",
	}, nil

}

func (s *Server) GetStoreApp(ctx context.Context, request *serverpb.GetStoreAppRequest) (
	*serverpb.GetStoreAppResponse, error) {
	if request.AppName == "" || request.Version == "" {
		s.log.Infof("AppName or version is missing for get store app request")
		return &serverpb.GetStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "AppName or version is missing in the request",
		}, nil
	}

	config, err := s.serverStore.GetAppFromStore(request.AppName, request.Version)
	if err != nil {
		s.log.Errorf("failed to get app config from store, %v", err)
		return &serverpb.GetStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get app config from store",
		}, nil
	}

	decodedIconBytes, _ := hex.DecodeString(config.Icon)
	appConfig := &serverpb.StoreAppConfig{
		AppName:             config.Name,
		Version:             config.Version,
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
		ReleaseName:         config.ReleaseName,
	}

	return &serverpb.GetStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly fetched from store",
		AppConfig:     appConfig,
	}, nil

}

func (s *Server) GetStoreApps(ctx context.Context, request *serverpb.GetStoreAppsRequest) (
	*serverpb.GetStoreAppsResponse, error) {

	configs, err := s.serverStore.GetAppsFromStore()
	if err != nil {
		s.log.Errorf("failed to get app config's from store, %v", err)
		return &serverpb.GetStoreAppsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get app config's from store",
		}, nil
	}

	appsData := []*serverpb.StoreAppsData{}
	for _, config := range *configs {
		decodedIconBytes, _ := hex.DecodeString(config.Icon)
		appsData = append(appsData, &serverpb.StoreAppsData{
			AppConfigs: &serverpb.StoreAppConfig{
				AppName:             config.Name,
				Version:             config.Version,
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
				ReleaseName:         config.ReleaseName,
			},
		})
	}

	return &serverpb.GetStoreAppsResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config's are sucessfuly fetched from store",
		Data:          appsData,
	}, nil
}

func (s *Server) GetStoreAppValues(ctx context.Context, request *serverpb.GetStoreAppValuesRequest) (
	*serverpb.GetStoreAppValuesResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}
	s.log.Infof("[org: %s] Get store app [%s:%s] values request for cluster %s recieved", orgId,
		request.AppName, request.Version, request.ClusterID)

	config, err := s.serverStore.GetAppFromStore(request.AppName, request.Version)
	if err != nil {
		s.log.Errorf("failed to get store app values, %v", err)
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get store app values",
		}, nil
	}

	marshaledOverride, err := yaml.Marshal(config.OverrideValues)
	if err != nil {
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to marshal values",
		}, nil
	}

	overrideValues, err := s.replaceGlobalValues(orgId, request.ClusterID, decodeBase64StringToBytes(string(marshaledOverride)))
	if err != nil {
		s.log.Errorf("failed to update overrided store app values, %v", err)
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update overrided store app values",
		}, nil
	}

	return &serverpb.GetStoreAppValuesResponse{
		Status:         serverpb.StatusCode_OK,
		StatusMessage:  "store app values sucessfuly fetched",
		OverrideValues: overrideValues,
	}, nil
}

func (s *Server) replaceGlobalValues(orgId, clusterID string, overridedValues []byte) ([]byte, error) {
	agent, err := s.agentHandeler.GetAgent(orgId, clusterID)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize agent for cluster %s", clusterID)
	}
	resp, err := agent.GetClient().GetClusterGlobalValues(context.TODO(), &agentpb.GetClusterGlobalValuesRequest{})
	if err != nil {
		return nil, err
	}
	if resp.Status != agentpb.StatusCode_OK {
		return nil, fmt.Errorf("failed to get global values for cluster %s", clusterID)
	}

	var globalValues map[string]interface{}
	err = yaml.Unmarshal(resp.GlobalValues, &globalValues)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal cluster values")
	}

	var overrideValues map[string]interface{}
	err = yaml.Unmarshal(overridedValues, &overrideValues)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal override values")
	}

	return replaceOverrideGlobalValues(overrideValues, globalValues)

}

func replaceOverrideGlobalValues(overrideValues map[string]interface{},
	globlaValues map[string]interface{}) (transformedData []byte, err error) {
	yamlData, err := yaml.Marshal(overrideValues)
	if err != nil {
		return
	}

	tmpl, err := template.New("templateVal").Parse(string(yamlData))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, globlaValues)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}

func (s *Server) DeployStoreApp(ctx context.Context, request *serverpb.DeployStoreAppRequest) (
	*serverpb.DeployStoreAppResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.DeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}
	s.log.Infof("[org: %s] Deploy store app [%s:%s] request for cluster %s recieved", orgId,
		request.AppName, request.Version, request.ClusterID)

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

	s.log.Infof("[org: %s] Store app [%s:%s] request request triggered for cluster %s", orgId,
		request.AppName, request.Version, request.ClusterID)

	return &serverpb.DeployStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app is successfully deployed",
	}, nil
}
