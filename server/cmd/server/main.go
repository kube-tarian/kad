package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/kube-tarian/kad/server/pkg/api"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/server"
)

var log = logging.NewLogger()

func main() {
	cfg, err := config.FetchConfiguration()
	if err != nil {
		log.Fatalf("Fetching application configuration failed, %v", err)
	}

	s, err := server.NewServer(log)
	if err != nil {
		log.Fatalf("Server initialization failed, %v", err)
	}

	r := gin.Default()
	api.Setup(r, s)

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
