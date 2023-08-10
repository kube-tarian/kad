package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/climon/pkg/application"
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
