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
	"scooter_micro/config"
	"scooter_micro/proto"
	"scooter_micro/repository"
	"scooter_micro/routing"
	"scooter_micro/routing/grpcserver"
	"scooter_micro/routing/httpserver"
	"scooter_micro/service"
)

var scooterIdMap = make(map[uint64]proto.ScooterService_RegisterServer)
var StructCh = make(chan *proto.ScooterClient)

func main() {
	log.Println("Starting scooter microservice")
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.PG_HOST,
		config.PG_PORT,
		config.POSTGRES_DB)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicf("%s: failed to open db connection - %v", "scooter_micro", err)
	}
	defer db.Close()

	scooterRepo := repository.NewScooterRepo(db)
	conn, err := grpc.DialContext(context.Background(), net.JoinHostPort("", config.ORDER_GRPC_PORT),
		grpc.WithInsecure())
	if err != nil {
		log.Printf("gRPC connection to %v port failed. With: %v\n", config.ORDER_GRPC_PORT, err)
	}

	log.Printf("gRPC connected port: %v.", config.ORDER_GRPC_PORT)

	orderClient := proto.NewOrderServiceClient(conn)
	scooterService := service.NewScooterService(scooterRepo, orderClient)
	scooterList, err := scooterService.GetAllScooters(context.Background(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
	}

	handler := routing.NewRouter(scooterService, StructCh)

	httpServer := httpserver.New(handler, StructCh, scooterService, httpserver.Port(config.HTTP_PORT))
	handler.HandleFunc("/scooter", httpServer.ScooterHandler)

	getIdFromStructInArray(scooterList, httpServer.ScooterIdMap)
	grpcServer := grpcserver.NewGrpcServer()
	proto.RegisterScooterServiceServer(grpcServer, httpServer)
	reflection.Register(grpcServer)

	http.ListenAndServe(":" + config.HTTP_PORT, handler)
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
