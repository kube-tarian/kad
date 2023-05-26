package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/types"

	"go.uber.org/zap"
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
func NewAgent(cfg *types.AgentConfiguration) (*Agent, error) {

	logger := log.GetLogger()
	defer logger.Sync()

	conn, err := getConnection(cfg)
	if err != nil {
		logger.Error("failed to connect", zap.Error(err))
		return nil, err
	}

	logger.Info("gRPC connection started",
		zap.String("address", cfg.Address),
		zap.Int("port", cfg.Port))
	agentClient := agentpb.NewAgentClient(conn)
	return &Agent{
		cfg:        cfg,
		connection: conn,
		client:     agentClient,
	}, nil
}

func getConnection(cfg *types.AgentConfiguration) (*grpc.ClientConn, error) {
	if !cfg.TlsEnabled {
		return grpc.Dial(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	tlsCreds, err := loadTLSCredentials(cfg)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		grpc.WithTransportCredentials(tlsCreds))

}

func (a *Agent) GetClient() agentpb.AgentClient {
	return a.client
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
