package captensdk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/integrator/capten-sdk/agentpb"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
)

type DeploymentRequestPayload struct {
	PluginName string                `json:"plugin_name" required:"true"`
	Action     string                `json:"action" required:"true"`
	Data       DeploymentRequestData `json:"data" required:"true"`
}

type DeploymentRequestData struct {
	RepoName    string `json:"repo_name" required:"true"`
	RepoURL     string `json:"repo_url" required:"true"`
	ChartName   string `json:"chart_name" required:"true"`
	Namespace   string `json:"namespace" required:"true"`
	ReleaseName string `json:"release_name" required:"true"`
	Timeout     int    `json:"timeout" default:"5"`
}

type ApplicationClient struct {
	log  logging.Logger
	conf *CaptenAgentConfiguration
	opts *TransportSSLOptions
}

func (c *Client) NewApplicationClient(opts *TransportSSLOptions) (*ApplicationClient, error) {
	return &ApplicationClient{log: c.log, conf: c.conf, opts: opts}, nil
}

func (a *ApplicationClient) Create(req *DeploymentRequestPayload) (*agentpb.JobResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		a.log.Errorf("Deploy request prepration failed", err)
		return nil, err
	}

	agentConn, conn, err := a.createAgentConnection()
	if err != nil {
		a.log.Errorf("agent client connection creation failed, %v", err)
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	return a.send(payload, agentConn)
}

func (*ApplicationClient) send(payload []byte, agentConn agentpb.AgentClient) (*agentpb.JobResponse, error) {
	jobReq := &agentpb.JobRequest{
		Operation: "deployment",
		Payload:   &anypb.Any{Value: payload},
	}
	return agentConn.SubmitJob(context.Background(), jobReq)
}

func (a *ApplicationClient) createAgentConnection() (agentpb.AgentClient, *grpc.ClientConn, error) {
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

// Delete... TODO: For delete all parameters not required.
// It has to be enhanced with separate delete payload request
func (a *ApplicationClient) Delete(req *DeploymentRequestPayload) (*agentpb.JobResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		a.log.Errorf("Deploy request prepration failed", err)
		return nil, err
	}

	agentConn, conn, err := a.createAgentConnection()
	if err != nil {
		a.log.Errorf("agent client connection creation failed, %v", err)
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	return a.send(payload, agentConn)
}

func (a *ApplicationClient) Update(req *DeploymentRequestPayload) (*agentpb.JobResponse, error) {
	return a.Create(req)
}
