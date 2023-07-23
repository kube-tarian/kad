package api

import (
	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"github.com/kube-tarian/kad/server/pkg/store"
)

type Server struct {
	serverpb.UnimplementedServerServer
	serverStore   store.ServerStore
	agentHandeler *agent.AgentHandler
	log           logging.Logger
}

func NewServer(log logging.Logger, serverStore store.ServerStore) (*Server, error) {
	return &Server{
		serverStore:   serverStore,
		agentHandeler: agent.NewAgentHandler(serverStore),
		log:           log,
	}, nil
}
