package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/opentelemetry"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) GetClusterApp(ctx context.Context, request *serverpb.GetClusterAppRequest) (
	*serverpb.GetClusterAppResponse, error) {

	_, span := opentelemetry.GetTracer(request.ClusterID).
		Start(opentelemetry.BuildContext(ctx), "CaptenServer")
	defer span.End()

	span.SetAttributes(attribute.String("Cluster ID", request.ClusterID))
	span.SetAttributes(attribute.String("AppReleaseName", request.AppReleaseName))
	return &serverpb.GetClusterAppResponse{}, fmt.Errorf("not implemented")
}
