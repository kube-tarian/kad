package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agent"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/config"
	"github.com/kube-tarian/kad/capten/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	"google.golang.org/grpc/reflection"
)

var log = logging.NewLogger()

func main() {
	log.Debugf("Staring Agent")

	cfg, err := config.FetchConfiguration()
	if err != nil {
		log.Fatalf("Fetching application configuration failed, %v", err)
	}

	//Todo: read from env var
	store := cassandra.NewCassandraStore(log, nil)
	err = store.Connect([]string{"localhost:9042"}, "cassandra", "cassandra")
	if err != nil {
		log.Fatalf("could not connect to cassandra, err: %v", err)
	}

	temporalClient, err := temporalclient.NewClient(log)
	if err != nil {
		log.Fatalf("Agent initialization failed, %v", err)
	}

	s, err := agent.NewAgent(log, temporalClient, store)
	if err != nil {
		log.Fatalf("Agent initialization failed, %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	agentpb.RegisterAgentServer(grpcServer, s)
	log.Infof("Server listening at %v", listener.Addr())
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("Failed to start agent : %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	grpcServer.Stop()
	log.Debugf("Exiting Agent")
}
