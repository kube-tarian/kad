package server

import (
	"context"
	"intelops.io/agent/pkg/agentpb"
)

type Agent struct {
	agentpb.UnimplementedAgentServer
}

func (a *Agent) SubmitJob(ctx context.Context, request *agentpb.JobRequest) (*agentpb.JobResponse, error) {
	return &agentpb.JobResponse{}, nil
}
