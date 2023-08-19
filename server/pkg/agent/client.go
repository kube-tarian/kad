package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"strings"

	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/pkg/errors"

	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	ClusterName string
	Address     string
	CaCert      string
	Cert        string
	Key         string
}

type Agent struct {
	cfg        *Config
	connection *grpc.ClientConn
	client     agentpb.AgentClient
	log        logging.Logger
}

func NewAgent(log logging.Logger, cfg *Config, oryclient oryclient.OryClient) (*Agent, error) {
	log.Infof("connecting to agent %s", cfg.Address)
	conn, err := getConnection(cfg, oryclient)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to connect to agent")
	}

	agentClient := agentpb.NewAgentClient(conn)
	pingResp, err := agentClient.Ping(context.TODO(), &agentpb.PingRequest{})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to ping agent")
	}
	if pingResp.Status != agentpb.StatusCode_OK {
		return nil, errors.WithMessage(err, "ping failed")
	}

	return &Agent{
		log:        log,
		cfg:        cfg,
		connection: conn,
		client:     agentClient,
	}, nil
}

func getConnection(cfg *Config, client oryclient.OryClient) (*grpc.ClientConn, error) {
	address, port, tls, err := parseAgentConnectionConfig(cfg.Address)
	if err != nil {
		return nil, err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithUnaryInterceptor(client.UnaryInterceptor),
	}

	if !tls {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsCreds, err := loadTLSCredentials(cfg)
		if err != nil {
			return nil, err
		}
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(tlsCreds))
	}

	return grpc.Dial(fmt.Sprintf("%s:%s", address, port), dialOptions...)
}

func (a *Agent) GetClient() agentpb.AgentClient {
	return a.client
}

func (a *Agent) Close() {
	a.connection.Close()
}

func loadTLSCredentials(config *Config) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(config.CaCert)) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	agentCert, err := tls.X509KeyPair([]byte(config.Cert), []byte(config.Key))
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{agentCert},
	}
	return credentials.NewTLS(tlsConfig), nil
}

func parseAgentConnectionConfig(agentAddress string) (agentHost, agentPort string, tlsEnabled bool, err error) {
	var parsedURL *url.URL
	parsedURL, err = url.Parse(agentAddress)
	if err != nil {
		return
	}
	agentHost = parsedURL.Host
	agentPort = parsedURL.Port()
	if strings.EqualFold(parsedURL.Scheme, "https") {
		tlsEnabled = true
		if len(agentPort) == 0 {
			agentPort = "443"
		}
		return
	}
	if len(agentPort) == 0 {
		agentPort = "80"
	}
	return
}
