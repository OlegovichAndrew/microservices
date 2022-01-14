package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"scooter_client/config"
)

//NewGrpcServer creates a new gRPC server on port 8080.
func NewGrpcServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", net.JoinHostPort("", config.GRPC_PORT))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		fmt.Printf("grpc server started on port: %v\n", config.GRPC_PORT)
		log.Fatal(grpcServer.Serve(listener))
	}()
	return grpcServer
}
