package api

import (
	"context"
	"encoding/base64"
	"sync"

	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/config"
	iamclient "github.com/kube-tarian/kad/server/pkg/iam-client"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	pluginstore "github.com/kube-tarian/kad/server/pkg/plugin-store"
	"github.com/kube-tarian/kad/server/pkg/store"
	"google.golang.org/grpc/metadata"
)

const (
	organizationIDAttribute  = "organizationid"
	clusterIDAttribute       = "clusterid"
	delayTimeinMin           = 15
	credentialAccessTokenKey = "accessToken"
)

type Server struct {
	serverpb.UnimplementedServerServer
	captenpluginspb.UnimplementedCaptenPluginsServer
	pluginstorepb.UnimplementedPluginStoreServer
	serverStore   store.ServerStore
	agentHandeler *agent.AgentHandler
	log           logging.Logger
	oryClient     oryclient.OryClient
	iam           iamclient.IAMRegister
	cfg           config.ServiceConfig
	pluginStore   *pluginstore.PluginStore
	mutex         *sync.Mutex
}

func NewServer(log logging.Logger, cfg config.ServiceConfig, serverStore store.ServerStore,
	oryClient oryclient.OryClient, iam iamclient.IAMRegister) (*Server, error) {
	agentHandeler := agent.NewAgentHandler(log, cfg, serverStore, oryClient)
	pluginStore, err := pluginstore.NewPluginStore(log, serverStore, agentHandeler, iam)
	if err != nil {
		return nil, err
	}
	return &Server{
		serverStore:   serverStore,
		agentHandeler: agentHandeler,
		log:           log,
		oryClient:     oryClient,
		iam:           iam,
		cfg:           cfg,
		pluginStore:   pluginStore,
		mutex:         &sync.Mutex{},
	}, nil
}

func metadataContextToMap(ctx context.Context) map[string]string {
	metadataMap := make(map[string]string)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadataMap
	}

	for key, values := range md {
		if len(values) > 0 {
			metadataMap[key] = values[0]
		}
	}
	return metadataMap
}

func encodeBase64BytesToString(val []byte) string {
	if len(val) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(val)
}

func decodeBase64StringToBytes(val string) []byte {
	if len(val) == 0 {
		return nil
	}
	dval, _ := base64.StdEncoding.DecodeString(val)
	return dval
}

func getBase64DecodedString(encodedString string) (string, error) {
	decodedByte, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "", err
	}
	return string(decodedByte), nil
}
