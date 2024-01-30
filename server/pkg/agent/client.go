package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
	"github.com/pkg/errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	ClusterName string
	Address     string
	CaCert      string
	Cert        string
	Key         string
	ServicName  string
	AuthEnabled bool
}

type Agent struct {
	cfg                 *Config
	connection          *grpc.ClientConn
	agentClient         agentpb.AgentClient
	captenPluginsClient captenpluginspb.CaptenPluginsClient
	log                 logging.Logger
}

func NewAgent(log logging.Logger, cfg *Config, oryClient oryclient.OryClient) (*Agent, error) {
	log.Infof("connecting to agent %s", cfg.Address)
	conn, err := getConnection(cfg, oryClient)
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

	captenPluginsClient := captenpluginspb.NewCaptenPluginsClient(conn)
	if captenPluginsClient == nil {
		return nil, errors.WithMessage(err, "failed to get agent capten plugins client")
	}

	return &Agent{
		log:                 log,
		cfg:                 cfg,
		connection:          conn,
		agentClient:         agentClient,
		captenPluginsClient: captenPluginsClient,
	}, nil
}

func getConnection(cfg *Config, oryClient oryclient.OryClient) (*grpc.ClientConn, error) {
	address, port, tls, err := parseAgentConnectionConfig(cfg.Address)
	if err != nil {
		return nil, err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithUnaryInterceptor(timeout.UnaryClientInterceptor(60 * time.Second)),
	}

	if cfg.AuthEnabled {
		dialOptions = append(dialOptions, grpc.WithUnaryInterceptor(authInterceptor(oryClient, cfg.ServicName)))
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

	return grpc.DialContext(context.Background(), fmt.Sprintf("%s:%s", address, port), dialOptions...)
}

func (a *Agent) GetClient() agentpb.AgentClient {
	return a.agentClient
}

func (a *Agent) GetCaptenPluginsClient() captenpluginspb.CaptenPluginsClient {
	return a.captenPluginsClient
}

func (a *Agent) Close() {
	a.connection.Close()
}

func authInterceptor(oryClient oryclient.OryClient, serviceName string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		oauthCred, err := oryClient.GetServiceOauthCredential(ctx, serviceName)
		if err != nil {
			return err
		}

		md := metadata.Pairs(
			"oauth_token", oauthCred.AccessToken,
			"ory_url", oauthCred.OryURL,
			"ory_pat", oauthCred.OryPAT,
		)
		newCtx := metadata.NewOutgoingContext(ctx, md)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
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
