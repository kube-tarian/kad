package grpcconn

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/capten-sdk/captensdkpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	conf   *GRPCConfig
	Client captensdkpb.CaptenSdkClient
}

type GRPCConfig struct {
	CaptenServerAddress string `envconfig:"CAPTEN_SERVER_ADDRESS" required:"true"`
}

func NewGRPCConfig() (*GRPCConfig, error) {
	conf := &GRPCConfig{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, fmt.Errorf("reading Capten server address failed, %v", err)
	}
	return conf, nil
}

func NewGRPCClient() (*GRPCClient, error) {
	conf, err := NewGRPCConfig()
	if err != nil {
		return nil, err
	}
	return NewGRPCClientWithConfig(conf)
}

func NewGRPCClientWithConfig(conf *GRPCConfig) (*GRPCClient, error) {
	conn, err := grpc.Dial(conf.CaptenServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("gRPC client did not connect to server: %v", err)
	}

	serverpb.NewServerClient(conn)
	client := &GRPCClient{
		conn:   conn,
		Client: captensdkpb.NewCaptenSdkClient(conn),
		conf:   conf,
	}
	return client, nil
}

func (g *GRPCClient) Close() {
	g.conn.Close()
}
