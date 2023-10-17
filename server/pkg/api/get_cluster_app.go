package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetClusterApp(ctx context.Context, request *serverpb.GetClusterAppRequest) (
	*serverpb.GetClusterAppResponse, error) {
	return &serverpb.GetClusterAppResponse{}, fmt.Errorf("not implemented")
}
