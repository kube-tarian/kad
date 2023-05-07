package main

import (
	"fmt"
	"github.com/kube-tarian/kad/vaultserv/pkg/api"
	"github.com/kube-tarian/kad/vaultserv/pkg/config"
	"github.com/kube-tarian/kad/vaultserv/pkg/pb/vaultservpb"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Debug("staring vaultserv")
	vaultServer, err := api.NewVaultServ()
	if err != nil {
		log.Fatal("failed to start vaultserv", zap.Error(err))
	}

	cfg, err := config.FetchConfiguration()
	if err != nil {
		log.Fatal("Fetching application configuration failed", zap.Error(err))
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	vaultservpb.RegisterVaultServer(grpcServer, vaultServer)
	log.Info("Server listening at 9098")

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("failed to start vaultserv", zap.Error(err))
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	grpcServer.Stop()
	log.Debug("exiting vaultserv")
}
