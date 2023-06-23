package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/climon/pkg/application"
	"github.com/kube-tarian/kad/capten/climon/pkg/db/cassandra"
	"github.com/kube-tarian/kad/capten/climon/pkg/temporal"
)

func main() {
	logger := logging.NewLogger()
	logger.Infof("Started deployment worker\n")

	db, err := cassandra.Create(logger)
	if err != nil {
		logger.Fatalf("failed to create db connection", err)
	}

	temporalObj := temporal.New(os.Getenv("TEMPORAL_ADDRESS"))
	if err := temporalObj.StartWorkers(); err != nil {
		log.Fatalln("failed to start worker", err)
	}

	app := application.New(logger, nil)
	go app.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	app.Close()
	db.Close()
	temporalObj.StopWorkers()
	logger.Infof("Exiting deployment worker\n")
}
