package api

import (
	"context"
	"encoding/base64"

	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/config"
	iamclient "github.com/kube-tarian/kad/server/pkg/iam-client"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"github.com/kube-tarian/kad/server/pkg/store"
	"google.golang.org/grpc/metadata"
)

const (
	organizationIDAttribute = "organizationid"
	clusterIDAttribute      = "clusterid"
)

type Server struct {
	serverpb.UnimplementedServerServer
	serverStore   store.ServerStore
	agentHandeler *agent.AgentHandler
	log           logging.Logger
	oryClient     oryclient.OryClient
	iam           iamclient.IAMRegister
	cfg           config.ServiceConfig
}

func NewServer(log logging.Logger, cfg config.ServiceConfig, serverStore store.ServerStore,
	oryClient oryclient.OryClient, iam iamclient.IAMRegister) (*Server, error) {
	return &Server{
		serverStore:   serverStore,
		agentHandeler: agent.NewAgentHandler(log, cfg, serverStore, oryClient),
		log:           log,
		oryClient:     oryClient,
		iam:           iam,
		cfg:           cfg,
	}, nil
}

func metadataContextToMap(ctx context.Context) map[string]string {
	metadataMap := make(map[string]string)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadataMap
	}

	for key, values := range md {
		if len(values) > 0 {
			metadataMap[key] = values[0]
		}
	}
	return metadataMap
}

func encodeBase64BytesToString(val []byte) string {
	if len(val) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(val)
}

func decodeBase64StringToBytes(val string) []byte {
	if len(val) == 0 {
		return nil
	}
	dval, _ := base64.StdEncoding.DecodeString(val)
	return dval
}
