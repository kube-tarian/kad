package main

import (
	"crypto/tls"
	"log"

	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	host := "<DBID>-<region>.apps.astra.datastax.com:443"
	token := "AstraCS:"
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	log.Printf("connecting to astra %s", host)
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(token),
		),
	)
	if err != nil {
		log.Fatalf("failed to connect to astra db, %v", err)
	}

	log.Printf("connected to %v", conn)
	session, err := client.NewStargateClientWithConn(conn)
	if err != nil {
		log.Fatalf("error creating stargate client, %v", err)
	}

	log.Printf("exec query")
	res, err := session.ExecuteQuery(&proto.Query{Cql: "DESCRIBE 'capten'"})
	if err != nil {
		log.Fatalf("error exec query, %v", err)
	}

	result := res.GetResult()
	log.Printf("result, %v", result)
}
