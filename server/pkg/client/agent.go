package client

import (
	"context"
	"fmt"
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Configuration struct {
	AgentAddress string `envconfig:"AGENT_ADDRESS" default:"localhost"`
	AgentPort    int    `envconfig:"AGENT_PORT" default:"9091"`
}

type Agent struct {
	cfg        *Configuration
	connection *grpc.ClientConn
	client     agentpb.AgentClient
	log        logging.Logger
}

// NewAgent returns agent object creates grpc connection for given address
func NewAgent(log logging.Logger) (*Agent, error) {
	cfg, err := fetchConfiguration()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.AgentAddress, cfg.AgentPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("failed to connect: %v", err)
		return nil, err
	}
	log.Infof("gRPC connection started to %s:%d", cfg.AgentAddress, cfg.AgentPort)

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
	a.log.Info("gRPC connection closed")
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	session, err := db.New()
	if err != nil {
		logrus.Error("failed to get db session", err)
		return nil, err
	}

	//todo make customerID dynamic : Ganesh
	endpoint, err := session.GetEndpoint("1")
	if err != nil {
		logrus.Error("failed to get db session", err)
		return nil, err
	}

	cfg.AgentAddress = endpoint
	return cfg, err
}
