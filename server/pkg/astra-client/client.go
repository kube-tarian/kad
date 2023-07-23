package astraclient

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
	astraCredDBIDKey     = "ASTRA_DB_ID"
	astraCredDBRegionKey = "ASTRA_DB_REGION"
	astraCredTokenKey    = "TOKEN"
)

type Config struct {
	DbHost               string `envconfig:"ASTRA_DB_HOST" default:"apps.astra.datastax.com"`
	EntityName           string `envconfig:"ASTRA_ENTITY_NAME" required:"true"`
	CredentailIdentifier string `envconfig:"ASTRA_CRED_IDENTIFIER" required:"true"`
	Keyspace             string `envconfig:"ASTRA_DB_NAME" default:"capten"`
}

type Client struct {
	session *client.StargateClient
	conf    *Config
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

	dbAddress, err := prepareDBAddress(conf, serviceCredential)
	if err != nil {
		return nil, err
	}
	token := serviceCredential[astraCredTokenKey]
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.Dial(dbAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(token),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to astra db, %w", err)
	}

	session, err := client.NewStargateClientWithConn(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating stargate client, %w", err)
	}
	return &Client{conf: conf, session: session}, nil
}

func (c *Client) Keyspace() string {
	return c.conf.Keyspace
}

func (c *Client) Session() *client.StargateClient {
	return c.session
}

func (c *Client) Close() {
}

func prepareDBAddress(conf *Config, serviceCredential map[string]string) (string, error) {
	dbID := serviceCredential[astraCredDBIDKey]
	region := serviceCredential[astraCredDBRegionKey]

	if len(dbID) == 0 || len(region) == 0 || len(serviceCredential[astraCredTokenKey]) == 0 {
		return "", fmt.Errorf("invalid credential")
	}
	return fmt.Sprintf("%s-%s.%s:443", serviceCredential[astraCredDBIDKey],
		serviceCredential[astraCredDBRegionKey], conf.DbHost), nil
}
