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

	"go.uber.org/zap"

	"github.com/kube-tarian/kad/server/api"
	rpcapi "github.com/kube-tarian/kad/server/pkg/api"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/kube-tarian/kad/server/pkg/handler"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println("failed to load configuration", err)
		return
	}

	if err := log.New(cfg.GetString("log.level")); err != nil {
		fmt.Println("failed to configure logger", err)
		return
	}

	logger := log.GetLogger()
	defer logger.Sync()

	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Fatal("Failed to get the swagger", zap.Error(err))
	}

	server, err := handler.NewAPIHandler()
	if err != nil {
		logger.Fatal("APIHandler initialization failed", zap.Error(err))
	}

	_, err = db.New(cfg.GetString("server.db"))
	if err != nil {
		logger.Fatal("Failed to connect to cassandra database", zap.Error(err))
	}

	rpcServer, err := rpcapi.New()
	if err != nil {
		logger.Fatal("grpc server initialization failed", zap.Error(err))
	}

	target := fmt.Sprintf("%s:%d", cfg.GetString("server.host"), cfg.GetInt("server.tcpPort"))
	listener, err := net.Listen("tcp", target)
	if err != nil {
		logger.Fatal("Failed to listen: ", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	serverpb.RegisterServerServer(grpcServer, rpcServer)
	logger.Info("Server listening at", zap.Any("address", listener.Addr()))
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("failed to start grpc server: ", zap.Error(err))
		}
	}()

	r := gin.Default()
	r.Use(middleware.OapiRequestValidator(swagger))
	r = api.RegisterHandlers(r, server)

	go func() {
		serverAddress := fmt.Sprintf("%s:%d", cfg.GetString("server.host"), cfg.GetInt("server.port"))
		if err := r.Run(serverAddress); err != nil {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	server.CloseAll()
	logger.Info("interrupt received, exiting")
}
