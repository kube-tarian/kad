package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/kube-tarian/kad/server/pkg/db"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Configuration struct {
	AgentAddress string `envconfig:"AGENT_ADDRESS" default:"localhost:9091"`
	IsSSLEnabled bool   `envconfig:"IS_SSL_ENABLED" default:"false"`
}

type Agent struct {
	cfg        *Configuration
	connection *grpc.ClientConn
	client     agentpb.AgentClient
	log        logging.Logger
}

// NewAgent returns agent object creates grpc connection for given address
func NewAgent(log logging.Logger) (*Agent, error) {
	cfg, err := fetchConfiguration(log)
	if err != nil {
		return nil, err
	}

	var conn *grpc.ClientConn
	if cfg.IsSSLEnabled {
		// TODO: loadTLSCredential to be implemented when mtls is introduced
		return nil, fmt.Errorf("SSL is not supported currently to agent")

		// tlsCredentials, lErr := loadTLSCredentials()
		// if lErr != nil {
		// 	log.Errorf("cannot load TLS credentials: ", lErr)
		// 	return nil, lErr
		// }
		// conn, err = grpc.Dial(cfg.AgentAddress, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		conn, err = grpc.Dial(cfg.AgentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if err != nil {
		log.Errorf("failed to connect: %v", err)
		return nil, err
	}
	log.Infof("gRPC connection started to %s", cfg.AgentAddress)

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

func fetchConfiguration(log logging.Logger) (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Errorf("env configuration fetch failed, %v", err)
		return nil, err
	}

	session, err := db.New()
	if err != nil {
		log.Error("failed to get db session", err)
		return nil, err
	}

	//todo make customerID dynamic : Ganesh
	endpoint, err := session.GetEndpoint("1")
	if err != nil {
		log.Error("failed to get db session", err)
		return nil, err
	}

	cfg.AgentAddress = endpoint
	return cfg, err
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	certificate, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		panic("Load client certification failed: " + err.Error())
	}

	pemServerCA, err := ioutil.ReadFile("certs/dev.optimizor.app.crt")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}
