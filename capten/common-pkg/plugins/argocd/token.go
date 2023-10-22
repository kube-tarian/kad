package argocd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	sessionpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/util/io"
)

const TokenPath = "api/v1/session"

type TokenResponse struct {
	Token string `json:"token" required:"true"`
}

func getNewAPIClient(cfg *Configuration) (apiclient.Client, error) {
	v, _ := json.Marshal(cfg)
	fmt.Println("Config valaues => " + string(v))
	fmt.Println("getNewAPIClient start")
	fmt.Println("cfg.ServiceURL => " + cfg.ServiceURL)
	fmt.Println("!cfg.IsSSLEnabled => " + fmt.Sprintf("%v", !cfg.IsSSLEnabled))
	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: cfg.ServiceURL,
		// Insecure:   !cfg.IsSSLEnabled,
	})
	if err != nil {
		fmt.Println("getNewAPIClient - NewClient - Error =>" + err.Error())
		return nil, err
	}

	fmt.Println("NewClient completed")

	sessConn, sessionClient, err := client.NewSessionClient()
	if err != nil {
		fmt.Println("sessionClient Error => " + err.Error())
		return nil, err
	}
	defer io.Close(sessConn)

	fmt.Println("sessConn sessionClient completed")

	sessionRequest := sessionpkg.SessionCreateRequest{
		Username: "admin",
		// Password: cfg.Password,
		Password: "Jgb4pU7gbY57PAnc",
	}
	createdSession, err := sessionClient.Create(context.Background(), &sessionRequest)
	if err != nil {
		fmt.Println("createdSession Error => " + err.Error())
		return nil, err
	}

	fmt.Println("createdSession completed")
	return apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: cfg.ServiceURL,
		Insecure:   !cfg.IsSSLEnabled,
		AuthToken:  createdSession.Token,
	})
}
