package main

import (
	"intelops.io/climon/pkg/temporal"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	temporalObj := temporal.New(os.Getenv("TEMPORAL_ADDRESS"))
	if err := temporalObj.StartWorkers(); err != nil {
		log.Fatalln("failed to start worker", err)
	}

	log.Println("Started workers")
	closeChan := make(chan bool)
	go handleSignal(closeChan)
	<-closeChan
	temporalObj.StopWorkers()
	log.Println("Stopped worker")
}

func handleSignal(closeChan chan bool) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	closeChan <- true
}
