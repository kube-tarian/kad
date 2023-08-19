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

func (c *Client) RegisterWithIam() error {
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
	err = credential.PutIamOauthCredential(context.Background(), res.ClientId, res.ClientSecret)
	if err != nil {
		return err
	}
	return nil
}

// at the line cm.WithIamYamlPath("provide the yaml location here"),
// the roles and actions should be added to ConfigMap
// the the location should be provided
func (c *Client) RegisterRolesActions() error {
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	// Create an instance of IamConn with desired options
	// the order of calling the options should be same as given in example
	iamConn := cm.NewIamConn(
		cm.WithGrpcDialOption(grpcOpts...),
		cm.WithIamAddress(c.cfg.IAMURL),
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
		"ory_url", c.oryClient.GetURL(),
		"ory_pat", c.oryClient.GetPAT(),
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
