package main

import (
	"fmt"
	"github.com/kube-tarian/kad/server/pkg/db"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/handler"
	"github.com/kube-tarian/kad/server/pkg/logging"
)

var log = logging.NewLogger()

func main() {
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Failed to get the swagger, %v", err)
	}

	cfg, err := config.FetchConfiguration()
	if err != nil {
		log.Fatalf("Fetching application configuration failed, %v", err)
	}

	server, err := handler.NewAPIHandler(log)
	if err != nil {
		log.Fatalf("APIHandler initialization failed, %v", err)
	}

	_, err = db.New()
	if err != nil {
		log.Fatalf("Failed to connect to cassandra database", err)
	}

	r := gin.Default()
	r.Use(middleware.OapiRequestValidator(swagger))
	r = api.RegisterHandlers(r, server)

	go func() {
		if err := r.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
			log.Fatalf("failed to start server : %s", err.Error())
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	server.CloseAll()
	log.Infof("Interrupt received, exiting")
}
