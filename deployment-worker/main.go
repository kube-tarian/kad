package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kube-tarian/kad/deployment-worker/pkg/application"
)

func main() {
	log.Printf("Started deployment worker\n")
	app := application.New()
	go app.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	app.Close()
	log.Printf("Exiting deployment worker\n")
}
