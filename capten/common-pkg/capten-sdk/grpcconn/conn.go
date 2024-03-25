package grpcconn

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/capten-sdk/captensdkpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn      *grpc.ClientConn
	conf      *GRPCConfig
	sdkClient captensdkpb.CaptenSdkClient
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

	client := &GRPCClient{
		conn:      conn,
		sdkClient: captensdkpb.NewCaptenSdkClient(conn),
		conf:      conf,
	}
	return client, nil
}

func (g *GRPCClient) Close() {
	g.conn.Close()
}

func (g *GRPCClient) SetupDatabase(ctx context.Context, dbOemName, dbName, serviceUserName string) (*captensdkpb.DBSetupResponse, error) {
	return g.sdkClient.SetupDatabase(ctx, &captensdkpb.DBSetupRequest{
		DbOemName:       dbOemName,
		DbName:          dbName,
		ServiceUserName: serviceUserName,
	})
}
