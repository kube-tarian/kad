package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kube-tarian/kad/agent/pkg/logging"
	rpcapi "github.com/kube-tarian/kad/server/pkg/api"
	"github.com/kube-tarian/kad/server/pkg/config"
	iamclient "github.com/kube-tarian/kad/server/pkg/iam-client"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"github.com/kube-tarian/kad/server/pkg/store"
)

func main() {
	log := logging.NewLogger()
	log.Infof("Staring Server")

	cfg, err := config.GetServiceConfig()
	if err != nil {
		log.Fatal("failed to load service congfig", err)
	}

	err = iamclient.RegisterService(log)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = iamclient.RegisterCerbosPolicy(log)
	if err != nil {
		log.Fatalf("%v", err)
	}

	serverStore, err := store.NewStore(log, cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to %s database", cfg.Database, err)
	}

	if cfg.CleanupDatabase {
		log.Infof("Cleaning database")
		err = serverStore.CleanupDatabase()
		if err != nil {
			log.Fatal("failed to initialize %s db, %w", cfg.Database, err)
		}
	}

	err = serverStore.InitializeDatabase()
	if err != nil {
		log.Fatal("failed to initialize %s db, %w", cfg.Database, err)
	}

	oryclient, err := oryclient.NewOryClient(log)
	if err != nil {
		log.Fatal("OryClient initialization failed", err)
	}

	iamCfg, err := iamclient.NewConfig()
	if err != nil {
		log.Fatal("faield to get the iam config", err)
	}

	iamClient, err := iamclient.NewClient(log, oryclient, iamCfg)
	if err != nil {
		log.Fatal("faield to initialize the iam client", err)
	}
	rpcServer, err := rpcapi.NewServer(log, cfg, serverStore, oryclient, iamClient)
	if err != nil {
		log.Fatal("grpc server initialization failed", err)
	}

	target := fmt.Sprintf("%s:%d", cfg.ServerGRPCHost, cfg.ServerGRPCPort)
	listener, err := net.Listen("tcp", target)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	var grpcServer *grpc.Server
	if cfg.AuthEnabled {
		log.Info("Server Authentication enabled")
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(rpcServer.AuthInterceptor))
	} else {
		log.Info("Server Authentication disabled")
		grpcServer = grpc.NewServer()
	}

	serverpb.RegisterServerServer(grpcServer, rpcServer)
	captenpluginspb.RegisterCaptenPluginsServer(grpcServer, rpcServer)
	pluginstorepb.RegisterPluginStoreServer(grpcServer, rpcServer)
	log.Info("Server listening at ", listener.Addr())
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("failed to start grpc server: ", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	log.Info("interrupt received, exiting")
}
