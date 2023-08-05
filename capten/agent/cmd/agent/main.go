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
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"
	"github.com/kube-tarian/kad/capten/agent/pkg/config"
	"google.golang.org/grpc/reflection"
)

var log = logging.NewLogger()

func main() {
	log.Infof("Staring Agent")

	cfg, err := config.GetServiceConfig()
	if err != nil {
		log.Fatalf("service config reading failed, %v", err)
	}

	if err := runAllMigrations(log); err != nil {
		log.Fatalf("Error while running migrations: %v", err)
	}

	s, err := agent.NewAgent(log)
	if err != nil {
		log.Fatalf("Agent initialization failed, %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	agentpb.RegisterAgentServer(grpcServer, s)
	log.Infof("Agent listening at %v", listener.Addr())
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

func runAllMigrations(log logging.Logger) error {
	return captenstore.Migrate(log)
}
