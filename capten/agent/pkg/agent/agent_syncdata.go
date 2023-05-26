package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
)

func (a *Agent) Sync(ctx context.Context, request *agentpb.SyncRequest) (*agentpb.SyncResponse, error) {
	fmt.Printf("%+v", request)
	var syncDataRequest workers.SyncDataRequest
	if err := json.Unmarshal([]byte(request.Data), &syncDataRequest); err != nil {
		return nil, err
	}

	syncDataRequest.Type = request.Type
	fmt.Printf("%+v", syncDataRequest)
	syncWorker := workers.NewSync(a.client)
	_, err := syncWorker.SendEvent(context.Background(), syncDataRequest)
	if err != nil {
		return nil, err
	}

	return &agentpb.SyncResponse{
		Status:  "Success",
		Message: "success",
	}, nil
}
