package captensdk

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/integrator/capten-sdk/agentpb"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClimonRequestPayload struct {
	PluginName string            `json:"plugin_name" required:"true"`
	Action     string            `json:"action" required:"true"`
	Data       ClimonRequestData `json:"data" required:"true"`
}

type ClimonRequestData struct {
	RepoName    string `json:"repo_name" required:"true"`
	RepoURL     string `json:"repo_url" required:"true"`
	ChartName   string `json:"chart_name" required:"true"`
	Namespace   string `json:"namespace" required:"true"`
	ReleaseName string `json:"release_name" required:"true"`
	Timeout     int    `json:"timeout" default:"5"`
}

type ClimonClient struct {
	log  logging.Logger
	conf *CaptenAgentConfiguration
	opts *TransportSSLOptions
}

func (c *Client) NewClimonClient(opts *TransportSSLOptions) (*ClimonClient, error) {
	return &ClimonClient{log: c.log, conf: c.conf, opts: opts}, nil
}

func (a *ClimonClient) createAgentConnection() (agentpb.AgentClient, *grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	if a.opts.IsSSLEnabled {
		tlsCredentials, lErr := loadTLSCredentials()
		if lErr != nil {
			a.log.Errorf("cannot load TLS credentials: ", lErr)
			return nil, nil, lErr
		}
		conn, err = grpc.Dial(fmt.Sprintf("%s:%d", a.conf.AgentAddress, a.conf.AgentPort), grpc.WithTransportCredentials(tlsCredentials))
	} else {
		conn, err = grpc.Dial(fmt.Sprintf("%s:%d", a.conf.AgentAddress, a.conf.AgentPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if err != nil {
		a.log.Errorf("failed to connect: %v", err)
		return nil, nil, err
	}
	a.log.Infof("gRPC connection started to %s:%d", a.conf.AgentAddress, a.conf.AgentPort)

	return agentpb.NewAgentClient(conn), conn, nil
}

func (a *ClimonClient) Create(req *agentpb.ClimonInstallRequest) (*agentpb.JobResponse, error) {
	agentConn, conn, err := a.createAgentConnection()
	if err != nil {
		a.log.Errorf("agent client connection creation failed, %v", err)
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	return agentConn.ClimonAppInstall(context.Background(), req)
}

func (a *ClimonClient) Update(req *agentpb.ClimonInstallRequest) (*agentpb.JobResponse, error) {
	return a.Create(req)
}

// Delete... TODO: For delete all parameters not required.
// It has to be enhanced with separate delete payload request
func (a *ClimonClient) Delete(req *agentpb.ClimonDeleteRequest) (*agentpb.JobResponse, error) {
	agentConn, conn, err := a.createAgentConnection()
	if err != nil {
		a.log.Errorf("agent client connection creation failed, %v", err)
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	return agentConn.ClimonAppDelete(context.Background(), req)
}
