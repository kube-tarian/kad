package oryclient

import (
	"context"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/credential"
	ory "github.com/ory/client-go"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Config struct {
	OryEntityName        string `envconfig:"ORY_ENTITY_NAME" required:"true"`
	CredentialIdentifier string `envconfig:"ORY_CRED_IDENTIFIER" required:"true"`
}

type OauthAccessCredential struct {
	OryPAT      string
	OryURL      string
	AccessToken string
}

type Client struct {
	oryPAT string
	conn   *ory.APIClient
	log    logging.Logger
	oryURL string
}

type OryClient interface {
	GetSessionTokenFromContext(ctx context.Context) (string, error)
	Authorize(ctx context.Context, accessToken string) (context.Context, error)
	GetServiceOauthCredential(ctx context.Context, serviceName string) (*OauthAccessCredential, error)
	GetURL() string
	GetPAT() string
}

func NewOryClient(log logging.Logger) (OryClient, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
		cfg.OryEntityName, cfg.CredentialIdentifier)
	if err != nil {
		return nil, err
	}
	oryPAT := serviceCredential.AdditionalData["ORY_PAT"]
	oryURL := serviceCredential.AdditionalData["ORY_URL"]
	conn := newOryAPIClient(log, oryURL)
	return &Client{
		oryPAT: oryPAT,
		conn:   conn,
		log:    log,
		oryURL: oryURL,
	}, nil
}

func newOryAPIClient(log logging.Logger, oryURL string) *ory.APIClient {
	config := ory.NewConfiguration()
	config.Servers = ory.ServerConfigurations{{
		URL: oryURL,
	}}
	return ory.NewAPIClient(config)
}

func (c *Client) GetURL() string {
	return c.oryURL
}

func (c *Client) GetPAT() string {
	return c.oryPAT
}

func (c *Client) GetOryTokenUrl() string {
	tokenUrl := c.oryURL + "/oauth2/token"
	return tokenUrl
}

func (c *Client) GetSessionTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Failed to get metadata from context")
	}

	bearerToken := md.Get("authorization")
	if len(bearerToken) == 0 || len(strings.Split(bearerToken[0], " ")) != 2 {
		return "", status.Error(codes.Unauthenticated, "No access token provided")
	}

	accessToken := bearerToken[0]
	if len(accessToken) < 8 || accessToken[:7] != "Bearer " {
		return "", status.Error(codes.Unauthenticated, "Invalid access token")
	}
	return accessToken[7:], nil
}

func (c *Client) Authorize(ctx context.Context, accessToken string) (context.Context, error) {
	ctx = context.WithValue(ctx, ory.ContextAccessToken, c.oryPAT)
	sessionInfo, _, err := c.conn.IdentityApi.GetSession(ctx, accessToken).Execute()
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "Failed to introspect session id, %v", err)
	}

	c.log.Infof("session id: %v", sessionInfo.Id)
	if !sessionInfo.GetActive() {
		return ctx, status.Error(codes.Unauthenticated, "session id is not active")
	}
	return ctx, nil
}

func (c *Client) GetServiceOauthCredential(ctx context.Context, serviceName string) (*OauthAccessCredential, error) {

	clientId, clientSecret, err := credential.GetServiceOauthCredential(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	conf := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid email offline"},
		TokenURL:     c.GetOryTokenUrl(),
	}

	oauthToken, err := conf.Token(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "error while fetching oauth token")
	}

	return &OauthAccessCredential{
		OryPAT:      c.oryPAT,
		OryURL:      c.oryURL,
		AccessToken: oauthToken.AccessToken,
	}, nil
}
