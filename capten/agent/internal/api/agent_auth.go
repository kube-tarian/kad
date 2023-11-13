package api

import (
	"context"

	"github.com/gogo/status"
	"github.com/intelops/go-common/logging"
	ory "github.com/ory/client-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func (a *Agent) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	tk, oryUrl, oryPat, err := a.extractDetailsFromContext(ctx)
	if err != nil {
		a.log.Errorf("error occured while extracting oauth token, oryurl, and ory pat token error: %v", err.Error())
		return nil, status.Error(codes.Unauthenticated, "invalid or missing token")
	}
	oryApiClient := newOryAPIClient(a.log, oryUrl)
	isValid, err := verifyToken(a.log, oryPat, tk, oryApiClient)
	if err != nil || !isValid {
		return nil, status.Error(codes.Unauthenticated, "invalid or missing token")
	}

	return handler(ctx, req)
}

func newOryAPIClient(log logging.Logger, oryURL string) *ory.APIClient {
	config := ory.NewConfiguration()
	config.Servers = ory.ServerConfigurations{{
		URL: oryURL,
	}}
	return ory.NewAPIClient(config)
}

func (a *Agent) extractDetailsFromContext(ctx context.Context) (oauthToken, oryURL, oryPAT string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		errMsg := "failed to extract metadata from context"
		a.log.Errorf(errMsg)
		return "", "", "", errors.New(errMsg)
	}

	if values, ok := md["oauth_token"]; ok && len(values) > 0 {
		oauthToken = values[0]
	} else {
		errMsg := "missing oauth_token in metadata"
		a.log.Errorf(errMsg)
		return "", "", "", errors.New(errMsg)
	}

	if values, ok := md["ory_url"]; ok && len(values) > 0 {
		oryURL = values[0]
	} else {
		errMsg := "missing ory_url in metadata"
		a.log.Errorf(errMsg)
		return "", "", "", errors.New(errMsg)
	}

	if values, ok := md["ory_pat"]; ok && len(values) > 0 {
		oryPAT = values[0]
	} else {
		errMsg := "missing ory_pat in metadata"
		a.log.Errorf(errMsg)
		return "", "", "", errors.New(errMsg)
	}

	return oauthToken, oryURL, oryPAT, nil
}

func verifyToken(log logging.Logger, oryPAT, token string, oryApiClient *ory.APIClient) (bool, error) {
	oryAuthedContext := context.WithValue(context.Background(), ory.ContextAccessToken, oryPAT)
	introspect, _, err := oryApiClient.OAuth2Api.IntrospectOAuth2Token(oryAuthedContext).Token(token).Scope("").Execute()
	if err != nil {
		log.Errorf("Failed to introspect token: %v", err)
		return false, err
	}
	if !introspect.Active {
		log.Error("Token is not active")
	}
	return introspect.Active, nil
}
