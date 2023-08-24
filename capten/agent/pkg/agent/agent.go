package agent

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogo/status"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"
	"github.com/kube-tarian/kad/capten/agent/pkg/config"
	"github.com/kube-tarian/kad/capten/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	ory "github.com/ory/client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"go.temporal.io/sdk/client"
)

var _ agentpb.AgentServer = &Agent{}

type Agent struct {
	agentpb.UnimplementedAgentServer
	tc  *temporalclient.Client
	as  *captenstore.Store
	log logging.Logger
}

func NewAgent(log logging.Logger, cfg *config.SericeConfig) (*Agent, error) {
	var tc *temporalclient.Client
	var err error

	tc, err = temporalclient.NewClient(log)
	if err != nil {
		return nil, err
	}

	as, err := captenstore.NewStore(log)
	if err != nil {
		// ignoring store failure until DB user creation working
		// return nil, err
		log.Errorf("failed to initialize store, %v", err)
	}

	agent := &Agent{
		tc:  tc,
		as:  as,
		log: log,
	}
	return agent, nil
}

func (a *Agent) Ping(ctx context.Context, request *agentpb.PingRequest) (*agentpb.PingResponse, error) {
	a.log.Infof("Ping request received")
	return &agentpb.PingResponse{Status: agentpb.StatusCode_OK}, nil
}

func (a *Agent) SubmitJob(ctx context.Context, request *agentpb.JobRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved event %+v", request)
	worker, err := a.getWorker(request.Operation)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	run, err := worker.SendEvent(ctx, request.Payload.GetValue())
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}

func (a *Agent) getWorker(operatoin string) (workers.Worker, error) {
	switch operatoin {
	default:
		return nil, fmt.Errorf("unsupported operation %s", operatoin)
	}
}

func prepareJobResponse(run client.WorkflowRun, name string) *agentpb.JobResponse {
	if run != nil {
		return &agentpb.JobResponse{Id: run.GetID(), RunID: run.GetRunID(), WorkflowName: name}
	}
	return &agentpb.JobResponse{}
}

func (a *Agent) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	tk, oryUrl, oryPat, err := a.extractDetailsFromContext(ctx)
	if err != nil {
		a.log.Errorf("error occured while extracting oauth token, oryurl, and ory pat token error: %v", err.Error())
		return nil, status.Error(codes.Unauthenticated, "invalid or missing token")
	}
	oryApiClient := NewOrySdk(a.log, oryUrl)
	isValid, err := verifyToken(a.log, oryPat, tk, oryApiClient)
	if err != nil || !isValid {
		return nil, status.Error(codes.Unauthenticated, "invalid or missing token")
	}

	return handler(ctx, req)
}

// NewOrySdk creates a oryAPIClient using the oryURL
// and returns it
func NewOrySdk(log logging.Logger, oryURL string) *ory.APIClient {
	log.Info("creating a ory client")
	config := ory.NewConfiguration()
	config.Servers = ory.ServerConfigurations{{
		URL: oryURL,
	}}

	return ory.NewAPIClient(config)
}

func (a *Agent) extractDetailsFromContext(ctx context.Context) (oauthToken, oryURL, oryPAT string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		errMsg := "failed to extract metadata from context"
		a.log.Errorf(errMsg)
		return "", "", "", errors.New(errMsg)
	}

	if values, ok := md["oauth_token"]; ok && len(values) > 0 {
		oauthToken = values[0]
	} else {
		errMsg := "missing oauth_token in metadata"
		a.log.Errorf(errMsg)
		return "", "", "", errors.New(errMsg)
	}

	if values, ok := md["ory_url"]; ok && len(values) > 0 {
		oryURL = values[0]
	} else {
		errMsg := "missing ory_url in metadata"
		a.log.Errorf(errMsg)
		return "", "", "", errors.New(errMsg)
	}

	if values, ok := md["ory_pat"]; ok && len(values) > 0 {
		oryPAT = values[0]
	} else {
		errMsg := "missing ory_pat in metadata"
		a.log.Errorf(errMsg)
		return "", "", "", errors.New(errMsg)
	}

	return oauthToken, oryURL, oryPAT, nil
}

func verifyToken(log logging.Logger, oryPAT, token string, oryApiClient *ory.APIClient) (bool, error) {
	oryAuthedContext := context.WithValue(context.Background(), ory.ContextAccessToken, oryPAT)
	introspect, _, err := oryApiClient.OAuth2Api.IntrospectOAuth2Token(oryAuthedContext).Token(token).Scope("").Execute()
	if err != nil {
		log.Errorf("Failed to introspect token: %v", err)
		return false, err
	}
	if !introspect.Active {
		log.Error("Token is not active")
	}
	return introspect.Active, nil
}
