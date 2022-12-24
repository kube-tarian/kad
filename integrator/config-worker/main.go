package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/config-worker/pkg/application"
)

func main() {
	logger := logging.NewLogger()
	logger.Infof("Started config worker\n")
	app := application.New(logger)
	go app.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	app.Close()
	logger.Infof("Exiting config worker\n")
}
