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
	"scooter_micro/proto"
	"scooter_micro/repository"
	"scooter_micro/routing"
	"scooter_micro/routing/grpcserver"
	"scooter_micro/routing/httpserver"
	"scooter_micro/service"
)

var scooterIdMap = make(map[uint64]proto.ScooterService_RegisterServer)
var Structure = make(chan *proto.ScooterClient)

func main() {
	log.Println("Starting scooter microservice")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

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
	scooterList, err := scooterService.GetAllScooters(context.Background(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
	}

	handler := routing.NewRouter(scooterService, Structure)

	httpServer := httpserver.New(handler, httpserver.Port("8085"))
	handler.HandleFunc("/scooter", httpServer.ScooterHandler)

	getIdFromStructInArray(scooterList, httpServer.ScooterIdMap)
	grpcServer := grpcserver.NewGrpcServer()
	proto.RegisterScooterServiceServer(grpcServer, httpServer)
	reflection.Register(grpcServer)

	http.ListenAndServe(":8085", handler)
}

func getIdFromStructInArray(from *proto.ScooterList,
	to map[uint64]proto.ScooterService_RegisterServer) map[uint64]proto.ScooterService_RegisterServer {
	for _, v := range from.Scooters {
		for i := 0; i < len(from.Scooters); i++ {
			to[v.Id] = nil
		}
	}
	return to
}
