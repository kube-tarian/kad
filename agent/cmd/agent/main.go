package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"intelops.io/agent/pkg/agentpb"
	"intelops.io/agent/pkg/config"
	"intelops.io/agent/pkg/server"
)

func main() {
	listener, err := net.Listen("tcp", config.DefaultAgentPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	agentpb.RegisterAgentServer(grpcServer, &server.Agent{})
	log.Println("server listening at ", listener.Addr())
	//Todo: handle the interrupts

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to start agent : %v", err)
	}
}
