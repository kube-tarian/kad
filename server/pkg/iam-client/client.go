package iamclient

import (
	"context"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/credential"
	iampb "github.com/kube-tarian/kad/server/pkg/pb/iampb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FetchSecret interface {
	GetSecrets(ctx context.Context, clientName, redirectURL string) (string, string, error)
	GetURL() string
}

type IAM struct {
	URL                  string
	log                  logging.Logger
	IamEntityName        string `envconfig:"IAM_ENTITY_NAME" required:"true"`
	CredentialIdentifier string `envconfig:"IAM_CRED_IDENTITY" required:"true"`
}

func (iam *IAM) GetURL() string {
	return iam.URL
}

func (iam *IAM) GetSecrets(ctx context.Context, clientName, redirectURL string) (string, string, error) {
	conn, err := grpc.Dial(iam.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (iam *IAM) RegisterWithIam() error {
	// each time whenever server starts it creates the client
	_, _, err := credential.GetIamOauthCredential(context.Background())
	if err == nil {
		iam.log.Info("Registration successful, re-using older registration..")
		return nil
	}

	conn, err := grpc.Dial(iam.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()

	iamclient := iampb.NewOauthServiceClient(conn)

	iam.log.Info("Registering capten as client in ory through...")

	oauthClientReq := &iampb.OauthClientRequest{
		ClientName:   "CaptenServer",
		RedirectUris: []string{"www.dummyurl.com"},
	}

	res, err := iamclient.CreateOauthClient(context.Background(), oauthClientReq)
	if err != nil {
		return err
	}

	err = credential.PutIamOauthCredential(context.Background(), res.ClientId, res.ClientSecret)
	if err != nil {
		return err
	}

	return nil
}

func New(log logging.Logger) (*IAM, error) {
	iam := &IAM{}
	if err := envconfig.Process("", iam); err != nil {
		return nil, err
	}

	iam.log = log

	serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
		iam.IamEntityName, iam.CredentialIdentifier)
	if err != nil {
		return nil, err
	}

	iam.URL = serviceCredential.AdditionalData["IAM_URL"]
	return iam, nil
}
