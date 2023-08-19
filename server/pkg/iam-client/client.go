package iamclient

import (
	"context"

	cm "github.com/intelops/go-common/iam"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/credential"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	iampb "github.com/kube-tarian/kad/server/pkg/pb/iampb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	IamURL string `envconfig:"IAM_URL" required:"true"`
}
type Client struct {
	oryClient oryclient.OryClient
	log       logging.Logger
	oryURL    string
	oryPAT    string
}

func NewClient(ory oryclient.OryClient, log logging.Logger) (*Client, error) {
	cfg, err := ory.GetOryEnv()
	if err != nil {
		return nil, err
	}
	serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
		cfg.OryEntityName, cfg.CredentialIdentifier)
	if err != nil {
		return nil, err
	}
	oryPAT := serviceCredential.AdditionalData["ORY_PAT"]
	oryURL := serviceCredential.AdditionalData["ORY_URL"]
	return &Client{
		oryClient: ory,
		log:       log,
		oryURL:    oryURL,
		oryPAT:    oryPAT,
	}, nil
}
func (c *Client) RegisterWithIam() error {
	cfg, err := getIamEnv()
	if err != nil {
		return err
	}

	iamURL := cfg.IamURL
	conn, err := grpc.Dial(iamURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	iamclient := iampb.NewOauthServiceClient(conn)
	c.log.Info("Registering capten as client in ory through...")
	oauthClientReq := &iampb.CreateClientCredentialsClientRequest{
		ClientName: "CaptenServer",
	}
	res, err := iamclient.CreateClientCredentialsClient(context.Background(), oauthClientReq)
	if err != nil {
		return err
	}
	err = credential.PutIamOauthCredential(context.Background(), res.ClientId, res.ClientSecret)
	if err != nil {
		return err
	}
	return nil
}

func getIamEnv() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// at the line cm.WithIamYamlPath("provide the yaml location here"),
// the roles and actions should be added to ConfigMap
// the the location should be provided
func (c *Client) RegisterRolesActions() error {
	cfg, err := getIamEnv()
	if err != nil {
		return err
	}

	iamURL := cfg.IamURL
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	// Create an instance of IamConn with desired options
	// the order of calling the options should be same as given in example
	iamConn := cm.NewIamConn(
		cm.WithGrpcDialOption(grpcOpts...),
		cm.WithIamAddress(iamURL),
		// TODO: here need to add the roles and actions yaml location
		cm.WithIamYamlPath("provide the yaml location here"),
	)
	ctx := context.Background()
	tkn, err := c.oryClient.GetCaptenServiceRegOauthToken()
	if err != nil {
		err = errors.WithMessage(err, "error getting capten service reg oauth token")
		return err
	}
	if tkn == nil {
		return errors.New("capten service reg oauth token is nil")
	}
	md := metadata.Pairs(
		"oauth_token", *tkn,
		"ory_url", c.oryURL,
		"ory_pat", c.oryPAT,
	)
	newCtx := metadata.NewOutgoingContext(ctx, md)
	// Update action roles
	err = iamConn.UpdateActionRoles(newCtx)
	if err != nil {
		c.log.Errorf("Failed to update action roles: %v", err)
		return err
	}
	return nil
}
