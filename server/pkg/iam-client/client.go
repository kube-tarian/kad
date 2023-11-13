package iamclient

import (
	"context"
	"log"

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

type IAMRegister interface {
	RegisterAppClientSecrets(ctx context.Context, clientName, redirectURL string) (string, string, error)
	GetOAuthURL() string
}

type CerbosEnv struct {
	CerbosUrl        string `envconfig:"CERBOS_URL" required:"true"`
	CerbosUsername   string `envconfig:"CERBOS_USERNAME" required:"true"`
	CerbosEntityName string `envconfig:"CERBOS_ENTITY_NAME" required:"true"`
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
func GetCerbosEnv() (*CerbosEnv, error) {
	cfg := &CerbosEnv{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Client) RegisterService() error {
	log.Println("Registering Service")
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
	log.Println("Registering Roles Actiosn")

	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	iamConn := cm.NewIamConn(
		cm.WithGrpcDialOption(grpcOpts...),
		cm.WithIamAddress(c.cfg.IAMURL),
		cm.WithIamYamlPath(c.cfg.ServiceRolesConfigFilePath),
		cm.WithOryCreds(c.oryClient.GetURL(), c.oryClient.GetPAT()),
	)
	if err := iamConn.InitializeOrySdk(); err != nil {
		return err
	}
	ctx := context.Background()
	oauthCred, err := c.oryClient.GetServiceOauthCredential(ctx, c.cfg.ServiceName)
	if err != nil {
		return errors.WithMessage(err, "error while getting service oauth token")
	}

	newCtx := metadata.AppendToOutgoingContext(context.Background(),
		"oauth_token", oauthCred.AccessToken)

	err = iamConn.UpdateActionRoles(newCtx)
	if err != nil {
		return errors.WithMessage(err, "Failed to update action roles")
	} else {
		log.Println("Updated Roles")
	}
	return nil
}

func (c *Client) GetOAuthURL() string {
	return c.oryClient.GetURL()
}

func (c *Client) RegisterAppClientSecrets(ctx context.Context, clientName, redirectURL string) (string, string, error) {
	conn, err := grpc.Dial(c.cfg.IAMURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", "", err
	}
	defer conn.Close()

	iamclient := iampb.NewOauthServiceClient(conn)
	res, err := iamclient.CreateOauthClient(context.Background(), &iampb.OauthClientRequest{
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
	ctx := context.Background()
	oauthCred, err := c.oryClient.GetServiceOauthCredential(ctx, c.cfg.ServiceName)
	if err != nil {
		return errors.WithMessage(err, "error while getting service oauth token")
	}

	iamConn := cm.NewIamConn(
		cm.WithGrpcDialOption(grpcOpts...),
		cm.WithIamAddress(c.cfg.IAMURL),
		cm.WithCerbosYamlPath(c.cfg.CerbosResourcePolicyFilePath),
		cm.WithOryCreds(c.oryClient.GetURL(), c.oryClient.GetPAT()),
	)
	if err := iamConn.InitializeOrySdk(); err != nil {
		return err
	}

	newCtx := metadata.AppendToOutgoingContext(context.Background(),
		"oauth_token", oauthCred.AccessToken)

	err = iamConn.RegisterCerbosResourcePolicies(newCtx)
	if err != nil {
		return errors.WithMessage(err, "Failed to update action roles")
	}
	return nil
}

func (c *Client) Interceptor() (*cm.ClientsAndConfigs, error) {
	log.Println("Interceptor Started")
	cfg, err := GetCerbosEnv()
	if err != nil {
		return nil, err
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
		cfg.CerbosEntityName, cfg.CerbosUsername)
	if err != nil {
		return nil, err
	}
	cerbosPassword := serviceCredential.Password
	iamConn := cm.NewIamConn(
		cm.WithGrpcDialOption(grpcOpts...),
		cm.WithIamAddress(c.cfg.IAMURL),
		cm.WithCerbosYamlPath(c.cfg.CerbosResourcePolicyFilePath),
		cm.WithInterceptorYamlPath(c.cfg.InterceptorYamlPath),
		cm.WithOryCreds(c.oryClient.GetURL(), c.oryClient.GetPAT()),
		cm.WithScope(c.cfg.ServiceName),
		cm.WithCerbosCreds(cfg.CerbosUrl, cfg.CerbosUsername, cerbosPassword),
	)
	if err := iamConn.InitializeOrySdk(); err != nil {
		return nil, err
	} else {
		log.Println("Initiaized Ory Sdk")
	}

	if err := iamConn.IntializeCerbosSdk(); err != nil {
		return nil, err
	}

	return iamConn, nil

}
