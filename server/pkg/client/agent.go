package client

import (
	"context"
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Configuration struct {
	Port int `envconfig:"PORT" default:"9091"`
}

type Agent struct {
	cfg        *Configuration
	connection *grpc.ClientConn
	client     agentpb.AgentClient
}

// NewAgent returns agent object creates grpc connection for given address
func NewAgent() (*Agent, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("failed to connect:", err)
		return nil, err
	}

	agentClient := agentpb.NewAgentClient(conn)
	return &Agent{
		cfg:        cfg,
		connection: conn,
		client:     agentClient,
	}, nil
}

func (a *Agent) GetClient() agentpb.AgentClient {
	return a.client
}

func (a *Agent) SubmitJob(ctx context.Context, req *agentpb.JobRequest) (*agentpb.JobResponse, error) {
	return a.client.SubmitJob(ctx, req)
}

func (a *Agent) Close() {
	a.connection.Close()
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
