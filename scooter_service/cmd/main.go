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
	"net/http"
	"os"
	"scooter_micro/config"
	"scooter_micro/proto"
	"scooter_micro/repository"
	"scooter_micro/routing"
	"scooter_micro/routing/grpcserver"
	"scooter_micro/routing/httpserver"
	"scooter_micro/service"
)

func main() {
	log.Println("Starting scooter microservice")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PG_HOST,
		config.PG_PORT,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicf("%s: failed to open db connection - %v", "scooter_micro", err)
	}
	defer db.Close()

	scooterRepo := repository.NewScooterRepo(db)
	conn, err := grpc.DialContext(context.Background(), net.JoinHostPort("", os.Getenv("ORDER_GRPC_PORT")),
		grpc.WithInsecure())
	if err != nil {
		log.Printf("gRPC connection to %v port failed. With: %v\n", os.Getenv("ORDER_GRPC_PORT"), err)
	}

	log.Printf("gRPC connected port: %v.", os.Getenv("ORDER_GRPC_PORT"))
	orderClient := proto.NewOrderServiceClient(conn)
	scooterService := service.NewScooterService(scooterRepo, orderClient)

	handler := routing.NewRouter(scooterService)
	httpServer := httpserver.New(handler, httpserver.Port("8085"))
	handler.HandleFunc("/scooter", httpServer.ScooterHandler)
	grpcServer := grpcserver.NewGrpcServer()
	proto.RegisterScooterServiceServer(grpcServer, httpServer)
	reflection.Register(grpcServer)
	http.ListenAndServe(":8085", handler)
}
