package astra

import (
	"crypto/tls"
	"fmt"
	"sync"

	vaultClient "github.com/kube-tarian/kad/server/pkg/Client"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	once     sync.Once
	astraObj *astra
)

const (
	keyspace      = "capten"
	astraHost     = "DB_HOST"
	astraPassword = "DB_PASSWORD"
)

type astra struct {
	stargateClient *client.StargateClient
}

func New() (*astra, error) {
	var err error
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	once.Do(func() {
		astraObj = &astra{}
		astraObj.stargateClient, err = connect()
		if err != nil {
			logger.Error("failed to connect to astra db", zap.Error(err))
		}
	})

	return astraObj, err
}

func connect() (*client.StargateClient, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Astra DB configuration
	//const astraUri = "0d175de3-c592-43f7-adf5-1bdda2761385-us-east1.apps.astra.datastax.com:443"
	//const bearerToken = "AstraCS:kYZPvIeLpthElpvKXQZUWHZF:32613fec5fe0be7f3cff755c2a09c5a411f0b0516d5521fc1fe8f3cbb3bf74ef"
	a, err := New()
	if err != nil {
		return nil, fmt.Errorf("error creating new astra object %w", err)
	}
	// TODO: pass customerId
	creds, err := a.FetchCreds("1", "vault")
	if err != nil {
		return nil, fmt.Errorf("error fetchingcreds from secret server %w", err)
	}
	host := creds.Host
	bearerToken := creds.Password
	// Create connection with authentication
	// For Astra DB:
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(bearerToken),
		),
	)

	stargateClient, err := client.NewStargateClientWithConn(conn)
	if err != nil {
		logger.Error("error creating stargate client", zap.Error(err))
		return nil, err
	}

	return stargateClient, nil
}

func (a *astra) GetAgentInfo(customerID string) (*types.AgentInfo, error) {
	agentInfo := types.AgentInfo{}
	selectQuery := &pb.Query{
		Cql: fmt.Sprintf("SELECT endpoint FROM %s.endpoints WHERE customer_id = %s", keyspace, customerID),
	}

	response, err := a.stargateClient.ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed fetch endpoint details %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	if len(result.Rows[0].Values) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	agentInfo.Endpoint, _ = client.ToString(result.Rows[0].Values[0])
	return &agentInfo, nil
}

func (a *astra) RegisterEndpoint(customerID, endpoint string, fileContentMap map[string]string) error {
	_, err := a.stargateClient.ExecuteQuery(&pb.Query{
		Cql: fmt.Sprintf("INSERT INTO %s.endpoints (customer_id, endpoint) VALUES ('%s', '%s');",
			keyspace, customerID, endpoint),
	})
	return err
}

func (a *astra) FetchCreds(customerID, secretPlugin string) (*types.DbCreds, error) {
	// TODO: based on secretPlugin fetch creds, currently vault is the secretPlugin
	vaultSession, err := vaultClient.NewVault()
	if err != nil {
		return nil, fmt.Errorf("failed to get vault session: %w", err)
	}

	// TODO: use exposed vault method when available
	credsMap, err := vaultSession.GetCreds("secret", customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get certs")
	}

	return &types.DbCreds{
		Host:     credsMap[astraHost],
		Password: credsMap[astraPassword],
	}, nil
}
