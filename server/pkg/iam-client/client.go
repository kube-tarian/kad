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

type Config struct {
	IamURL string `envconfig:"IAM_URL" required:"true"`
}

func RegisterWithIam(log logging.Logger) error {
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
	log.Info("Registering capten as client in ory through...")
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
