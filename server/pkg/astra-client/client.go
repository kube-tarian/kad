package astra

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	astraCredHostKey  = "host"
	astraCredTokenKey = "token"
)

type Config struct {
	EntityName           string `envconfig:"ASTRA_ENTITY_NAME" required:"true"`
	CredentailIdentifier string `envconfig:"ASTRA_CRED_IDENTIFIER" required:"true"`
}

type Client struct {
	session *client.StargateClient
}

func NewClient() (*Client, error) {
	conf := &Config{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, fmt.Errorf("astra config read faile, %v", err)
	}

	serviceCredential, err := credential.GetGenericCredential(context.Background(),
		conf.EntityName, conf.CredentailIdentifier)
	if err != nil {
		return nil, err
	}

	host := serviceCredential[astraCredHostKey]
	password := serviceCredential[astraCredTokenKey]
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(password),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to astra db, %w", err)
	}

	session, err := client.NewStargateClientWithConn(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating stargate client, %w", err)
	}
	return &Client{session: session}, nil
}

func (c *Client) Session() *client.StargateClient {
	return c.session
}

func (c *Client) Close() {
}
