package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"order_micro/config"
	"order_micro/proto"
	"order_micro/repository"
	"order_micro/service"
)

func main() {
	log.Println("Starting order microservice")
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.PG_HOST,
		config.PG_PORT,
		config.POSTGRES_DB)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicf("%s: failed to open db connection - %v", "order_micro", err)
	}
	defer db.Close()

	orderRepo := repository.NewOrderRepo(db)
	service := service.NewOrderService(orderRepo)

	listener, err := net.Listen("tcp", net.JoinHostPort("", config.ORDER_GRPC_PORT))
	if err != nil {
		log.Panicf("%s: failed to listen on port - %v","order_micro", err)
	}

	server := grpc.NewServer()
	proto.RegisterOrderServiceServer(server, service)
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		log.Panicf("%s: failed to start grpc - %v", "order_micro", err)
	}
}
