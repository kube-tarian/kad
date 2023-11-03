package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	captenmodel "github.com/kube-tarian/kad/capten/model"
)

const (
	crossplaneConfigUseCase string = "crossplane"
)

func (a *Agent) RegisterCrossplaneProject(ctx context.Context, request *captenpluginspb.RegisterCrossplaneProjectRequest) (
	*captenpluginspb.RegisterCrossplaneProjectResponse, error) {

	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Register Crossplane Git project %s request recieved", request.Id)

	crossplaneProject, err := a.as.GetCrossplaneProjectForID(request.Id)
	if err != nil {
		a.log.Infof("failed to get git project %s, %v", request.Id, err)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	if crossplaneProject.Status != string(model.CrossplaneProjectConfigurationFailed) &&
		crossplaneProject.Status != string(model.CrossplaneProjectAvailable) {
		a.log.Infof("currently the Crossplane project configuration on-going %s, %v", request.Id, crossplaneProject.Status)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_OK,
			StatusMessage: "Crossplane configuration on-going",
		}, nil
	}

	crossplaneProject.Status = string(model.CrossplaneProjectConfigurationOngoing)
	if err := a.as.UpsertCrossplaneProject(crossplaneProject); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "inserting data to Crossplane db got failed",
		}, err
	}

	if ok, err := a.isProjectRegisteredWithArgoCD(ctx, crossplaneProject.GitProjectUrl); !ok && err == nil {
		accessToken, userID, err := a.getGitProjectCredential(ctx, crossplaneProject.GitProjectId)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.RegisterCrossplaneProjectResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "Error occured while fetching Crossplane git project AccessToken and User Id",
			}, nil
		}

		if err := a.addProjectToArgoCD(ctx, crossplaneProject.GitProjectUrl, userID, accessToken); err != nil {
			a.log.Errorf("failed to add Repository to ArgoCD : %v ", err)
			return &captenpluginspb.RegisterCrossplaneProjectResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "Error occured while adding Crossplane Repository to ArgoCD",
			}, err
		}
	} else if err != nil {
		a.log.Errorf("failed to add Repository to ArgoCD : %v ", err)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Failed to check weather Crossplane Repository is added to ArgoCD or not",
		}, err
	}

	providers, err := a.GetCrossplaneProviders()
	if err != nil {
		a.log.Errorf("failed to fetch provider info, %v", err)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch provider info",
		}, nil
	}
	// start the config-worker routine
	go a.configureCrossplaneGitRepo(crossplaneProject, providers)

	a.log.Infof("Crossplane Git project %s registration triggerred", request.Id)
	return &captenpluginspb.RegisterCrossplaneProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully registered crossplane project",
	}, nil
}

func (a *Agent) GetCrossplaneProject(ctx context.Context, request *captenpluginspb.GetCrossplaneProjectsRequest) (
	*captenpluginspb.GetCrossplaneProjectsResponse, error) {
	a.log.Infof("Get Crossplane Git projects request recieved")

	project, err := a.as.GetCrossplaneProject()
	if err != nil {
		a.log.Errorf("failed to get Crossplane Project, %v", err)
		return &captenpluginspb.GetCrossplaneProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get Crossplane Project",
		}, err
	}

	crossplaneProject := &captenpluginspb.CrossplaneProject{
		Id:            project.Id,
		GitProjectUrl: project.GitProjectUrl,
		Status:        project.Status,
	}

	a.log.Infof("Fetched Crossplane Git project, id: %v", project.Id)
	return &captenpluginspb.GetCrossplaneProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the Crossplane projects",
		Project:       crossplaneProject,
	}, nil
}

func (a *Agent) UnRegisterCrossplaneProject(ctx context.Context, request *captenpluginspb.UnRegisterCrossplaneProjectRequest) (
	*captenpluginspb.UnRegisterCrossplaneProjectResponse, error) {

	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UnRegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("UnRegister Crossplane Git project %s request recieved", request.Id)

	crossplaneProject, err := a.as.GetCrossplaneProjectForID(request.Id)
	if err != nil {
		a.log.Infof("failed to get git project %s, %v", request.Id, err)
		return &captenpluginspb.UnRegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	crossplaneProject.Status = string(model.CrossplaneProjectAvailable)
	if err := a.as.UpsertCrossplaneProject(crossplaneProject); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &captenpluginspb.UnRegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "inserting data to Crossplane db got failed",
		}, err
	}

	a.log.Infof("UnRegister Crossplane Git project %s request processed", request.Id)
	return &captenpluginspb.UnRegisterCrossplaneProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully delete the Crossplane project",
	}, nil
}

func (a *Agent) configureCrossplaneGitRepo(req *model.CrossplaneProject, providers []model.CrossplaneProvider) {
	ci := captenmodel.CrossplaneUseCase{Type: crossplaneConfigUseCase, RepoURL: req.GitProjectUrl,
		VaultCredIdentifier: req.GitProjectId, PushToDefaultBranch: !a.createPr, CrossplaneProviders: providers}
	wd := workers.NewConfig(a.tc, a.log)

	wkfId, err := wd.SendAsyncEvent(context.TODO(), &captenmodel.ConfigureParameters{Resource: crossplaneConfigUseCase}, ci)
	if err != nil {
		req.Status = string(model.CrossplaneProjectConfigurationFailed)
		req.WorkflowId = "NA"
		if err := a.as.UpsertCrossplaneProject(req); err != nil {
			a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
			return
		}
		a.log.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectUrl, err)
		return
	}

	a.log.Infof("Crossplane Git project %s config workflow event %s created", wkfId)

	req.Status = string(model.CrossplaneProjectConfigured)
	req.WorkflowId = wkfId
	req.WorkflowStatus = string(model.WorkFlowStatusStarted)
	if err := a.as.UpsertCrossplaneProject(req); err != nil {
		a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
		return
	}

	go a.monitorCrossplaneWorkflow(req, wkfId)
	a.log.Infof("Crossplane Git project %s registration completed", req.Id)
}

func (a *Agent) monitorCrossplaneWorkflow(req *model.CrossplaneProject, wkfId string) {
	// during system reboot start monitoring, add it in map or somewhere.
	wd := workers.NewConfig(a.tc, a.log)
	wkfResp, err := wd.GetWorkflowInformation(context.TODO(), wkfId)
	if err != nil {
		req.Status = string(model.CrossplaneProjectConfigurationFailed)
		req.WorkflowStatus = string(model.WorkFlowStatusFailed)
		if err := a.as.UpsertCrossplaneProject(req); err != nil {
			a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
			return
		}
		a.log.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectUrl, err)
		return
	}

	a.log.Infof("Monitoring Crossplane Git project %s config workflow event %s created", wkfId)

	req.Status = string(model.CrossplaneProjectConfigured)
	req.WorkflowStatus = wkfResp.Status
	if err := a.as.UpsertCrossplaneProject(req); err != nil {
		a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
		return
	}

	a.log.Infof("Crossplane Git project %s monitoring completed", req.Id)
}

func (a *Agent) GetCrossplaneProviders() (providers []model.CrossplaneProvider, err error) {
	crossplaneProviders, err := a.as.GetCrossplaneProviders()
	if err != nil {
		return nil, err
	}
	for _, crossplaneProvider := range crossplaneProviders {
		crossplaneProviderObj := model.CrossplaneProvider{
			Id:              crossplaneProvider.GetId(),
			CloudType:       crossplaneProvider.GetCloudType(),
			ProviderName:    crossplaneProvider.GetProviderName(),
			CloudProviderId: crossplaneProvider.GetCloudProviderId(),
			Status:          crossplaneProvider.GetStatus(),
		}
		providers = append(providers, crossplaneProviderObj)
	}
	return
}
