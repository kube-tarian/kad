package iamclient

import (
	"context"

	cm "github.com/intelops/go-common/iam"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/server/pkg/credential"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	iampb "github.com/kube-tarian/kad/server/pkg/pb/iampb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type IAMRegister interface {
	RegisterAppClientSecrets(ctx context.Context, clientName, redirectURL, organisationid string) (string, string, error)
	GetOAuthURL() string
}

type Client struct {
	oryClient oryclient.OryClient
	log       logging.Logger
	cfg       Config
}

func NewClient(log logging.Logger, ory oryclient.OryClient, cfg Config) (*Client, error) {
	return &Client{
		oryClient: ory,
		log:       log,
		cfg:       cfg,
	}, nil
}

func (c *Client) RegisterService() error {
	conn, err := grpc.Dial(c.cfg.IAMURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	iamclient := iampb.NewOauthServiceClient(conn)
	oauthClientReq := &iampb.CreateClientCredentialsClientRequest{
		ClientName: c.cfg.ServiceName,
	}
	res, err := iamclient.CreateClientCredentialsClient(context.Background(), oauthClientReq)
	if err != nil {
		return err
	}

	err = credential.StoreServiceOauthCredential(context.Background(), c.cfg.ServiceName, res.ClientId, res.ClientSecret)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) RegisterRolesActions() error {

	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	iamConn := cm.NewIamConn(
		cm.WithGrpcDialOption(grpcOpts...),
		cm.WithIamAddress(c.cfg.IAMURL),
		cm.WithIamYamlPath(c.cfg.ServiceRolesConfigFilePath),
	)

	ctx := context.Background()
	oauthCred, err := c.oryClient.GetServiceOauthCredential(ctx, c.cfg.ServiceName)
	if err != nil {
		return errors.WithMessage(err, "error while getting service oauth token")
	}

	newCtx := metadata.AppendToOutgoingContext(context.Background(),
		"oauth_token", oauthCred.AccessToken, "ory_url", c.oryClient.GetURL(), "ory_pat", c.oryClient.GetPAT())

	err = iamConn.UpdateActionRoles(newCtx)
	if err != nil {
		return errors.WithMessage(err, "Failed to update action roles")
	}
	return nil
}

func (c *Client) GetOAuthURL() string {
	return c.oryClient.GetURL()
}

func (c *Client) RegisterAppClientSecrets(ctx context.Context, clientName, redirectURL, organisationid string) (string, string, error) {
	conn, err := grpc.Dial(c.cfg.IAMURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", "", err
	}
	defer conn.Close()

	iamclient := iampb.NewOauthServiceClient(conn)
	md := metadata.Pairs(
		"organisationid", organisationid,
	)
	newCtx := metadata.NewOutgoingContext(ctx, md)
	res, err := iamclient.CreateOauthClient(newCtx, &iampb.OauthClientRequest{
		ClientName: clientName, RedirectUris: []string{redirectURL},
	})
	if err != nil {
		return "", "", err
	}
	return res.ClientId, res.ClientSecret, nil
}

func (c *Client) RegisterCerbosPolicy() error {
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	iamConn := cm.NewIamConn(
		cm.WithGrpcDialOption(grpcOpts...),
		cm.WithIamAddress(c.cfg.IAMURL),
		cm.WithCerbosYamlPath(c.cfg.CerbosResourcePolicyFilePath),
	)

	ctx := context.Background()
	oauthCred, err := c.oryClient.GetServiceOauthCredential(ctx, c.cfg.ServiceName)
	if err != nil {
		return errors.WithMessage(err, "error while getting service oauth token")
	}

	newCtx := metadata.AppendToOutgoingContext(context.Background(),
		"oauth_token", oauthCred.AccessToken, "ory_url", c.oryClient.GetURL(), "ory_pat", c.oryClient.GetPAT())

	err = iamConn.RegisterCerbosResourcePolicies(newCtx)
	if err != nil {
		return errors.WithMessage(err, "Failed to update action roles")
	}
	return nil
}
