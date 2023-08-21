package api

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func (s *Server) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
