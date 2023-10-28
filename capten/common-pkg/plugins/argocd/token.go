package argocd

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	sessionpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/util/io"
)

const TokenPath = "api/v1/session"

type TokenResponse struct {
	Token string `json:"token" required:"true"`
}

func getNewAPIClient(cfg *Configuration) (apiclient.Client, error) {

	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: cfg.ServiceURL,
		Insecure:   !cfg.IsSSLEnabled,
	})
	if err != nil {
		return nil, err
	}

	sessConn, sessionClient, err := client.NewSessionClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(sessConn)

	sessionRequest := sessionpkg.SessionCreateRequest{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	createdSession, err := sessionClient.Create(context.Background(), &sessionRequest)
	if err != nil {
		return nil, err
	}

	return apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: cfg.ServiceURL,
		Insecure:   !cfg.IsSSLEnabled,
		AuthToken:  createdSession.Token,
	})
}
