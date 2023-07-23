package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
)

func (a *Agent) SyncApp(ctx context.Context, request *agentpb.SyncAppRequest) (*agentpb.SyncAppResponse, error) {
	// var appConfig types.AppConfig
	// if err := yaml.Unmarshal(request.Payload, &appConfig); err != nil {
	// 	a.log.Errorf("could not unmarshal appConfig yaml: %v", err)
	// 	return nil, err
	// }

	// if err := a.as.AddAppConfig(appConfig); err != nil {
	// 	a.log.Errorf("could not insert, err: %v", err)
	// 	return &agentpb.SyncAppResponse{
	// 		Status:        agentpb.StatusCode(1),
	// 		StatusMessage: "FAILED",
	// 	}, err
	// }

	return &agentpb.SyncAppResponse{
		Status:        agentpb.StatusCode(0),
		StatusMessage: "SUCCESS",
	}, nil
}
