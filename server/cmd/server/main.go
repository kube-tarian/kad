package main

import (
	//
	//"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	//"github.com/kube-tarian/kad/server/pkg/handler"

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

	s, err := handler.NewAPIHandler(log)
	if err != nil {
		log.Fatalf("APIHandler initialization failed, %v", err)
	}

	r := gin.Default()
	r.Use(middleware.OapiRequestValidator(swagger))
	r = api.RegisterHandlers(r, s)

	//this part from 43-49 is used to invoke the /deploy endpoint to run and test the handle deploy locally
	handler := &handler.APIHanlder{}
	router := gin.Default()
	router.POST("/deploy", func(c *gin.Context) {

		handler.PostDeploy(c)
	})

	go func() {
		if err := r.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
			log.Fatalf("failed to start server : %s", err.Error())
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Infof("Interrupt received, exiting")
}
