package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/application"
	"github.com/kube-tarian/kad/integrator/pkg/logging"
)

func main() {
	logger := logging.NewLogger()
	logger.Infof("Started deployment worker\n")
	app := application.New(logger)
	go app.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	app.Close()
	logger.Infof("Exiting deployment worker\n")
}
