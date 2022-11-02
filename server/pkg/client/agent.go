package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"intelops.io/server/pkg/pb/agentpb"
	"log"
)

type Agent struct {
	connection *grpc.ClientConn
	client     agentpb.AgentClient
}

// NewAgent returns agent object creates grpc connection for given address
func NewAgent(address string) (*Agent, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("failed to connect:", err)
		return nil, err
	}

	agentClient := agentpb.NewAgentClient(conn)
	return &Agent{
		connection: conn,
		client:     agentClient,
	}, nil
}

func (a *Agent) GetClient() agentpb.AgentClient {
	return a.client
}

func (a *Agent) Close() {
	a.connection.Close()
}
