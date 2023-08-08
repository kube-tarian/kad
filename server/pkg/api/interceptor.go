package api

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// UnaryInterceptor is a gRPC server-side interceptor that handles authentication for unary RPCs.
// It first attempts to retrieve an access token from the context using the ORY client interface.
// If the token retrieval is successful, it then tries to authorize the token using the ORY client interface.
// If either step fails, the interceptor logs the error and returns it, halting the RPC.
// If both steps are successful, the interceptor invokes the provided handler with the updated context and request.
func (s Server) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	accessToken, err := s.oryClient.GetSessionTokenFromContext(ctx)
	if err != nil {
		s.log.Debugf("error occured while fetching the token from the context. Error - %s", err.Error())
		return nil, err
	}

	ctx, err = s.oryClient.Authorize(ctx, accessToken)
	if err != nil {
		s.log.Info(fmt.Sprintf("Error occurred while authorizing the session id from context. Session Id - %s\nError - %s", accessToken, err.Error()))
		return nil, err
	}

	return handler(ctx, req)
}
