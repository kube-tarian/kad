package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-gonic/gin"

	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/server/api"
	rpcapi "github.com/kube-tarian/kad/server/pkg/api"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/handler"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
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

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal("Failed to get the swagger", err)
	}

	serverStore, err := store.NewStore(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to %s database", cfg.Database, err)
	}

	err = serverStore.InitializeDb()
	if err != nil {
		log.Fatal("failed to initialize %s db, %w", cfg.Database, err)
	}

	server, err := handler.NewAPIHandler(log, serverStore)
	if err != nil {
		log.Fatal("APIHandler initialization failed", err)
	}
	oryclient, err := oryclient.NewOryClient(log)
	if err != nil {
		log.Fatal("OryClient initialization failed", err)
	}
	rpcServer, err := rpcapi.NewServer(log, serverStore, oryclient)
	if err != nil {
		log.Fatal("grpc server initialization failed", err)
	}

	target := fmt.Sprintf("%s:%d", cfg.ServerGRPCHost, cfg.ServerGRPCPort)
	listener, err := net.Listen("tcp", target)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(rpcServer.UnaryInterceptor))
	serverpb.RegisterServerServer(grpcServer, rpcServer)
	log.Info("Server listening at ", listener.Addr())
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("failed to start grpc server: ", err)
		}
	}()

	r := gin.Default()
	r.Use(middleware.OapiRequestValidator(swagger))
	r = api.RegisterHandlers(r, server)

	go func() {
		serverAddress := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
		if err := r.Run(serverAddress); err != nil {
			log.Fatal("failed to start server", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	log.Info("interrupt received, exiting")
}
