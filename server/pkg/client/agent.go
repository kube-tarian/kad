package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Agent struct {
	cfg        *types.AgentConfiguration
	connection *grpc.ClientConn
	client     agentpb.AgentClient
	log        logging.Logger
}

// NewAgent returns agent object creates grpc connection for given address
func NewAgent(log logging.Logger, cfg *types.AgentConfiguration) (*Agent, error) {
	tlsCreds, err := loadTLSCredentials(cfg)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		log.Errorf("failed to connect: %v", err)
		return nil, err
	}

	log.Infof("gRPC connection started to %s:%d", cfg.Address, cfg.Port)
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

func loadTLSCredentials(cfg *types.AgentConfiguration) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(cfg.CaCert)) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	agentCert, err := tls.X509KeyPair([]byte(cfg.Cert), []byte(cfg.Key))
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{agentCert},
	}

	return credentials.NewTLS(config), nil
}
