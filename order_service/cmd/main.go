package main

import (
	"context"
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
	"order_micro/transport"
)

const TopicName = "order"
const ClientID = "some_client"
const GroupConsumer = "some_group"


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

	group := transport.CreateConsumerGroup([]string{config.KAFKA_BROKER}, ClientID, GroupConsumer)

	go func() {
		for  {
			transport.ConsumeMessages(context.Background(), group, TopicName)
		}
	}()

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