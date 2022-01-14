package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"scooter_client/config"
	"scooter_client/proto"
)

func main() {
	log.Println("Starting scooter client")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PG_HOST,
		config.PG_PORT,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicf("%s: failed to open db connection - %v", "order_micro", err)
	}
	defer db.Close()

	conn, err := grpc.DialContext(context.Background(), net.JoinHostPort("", os.Getenv("ORDER_GRPC_PORT")),
		grpc.WithInsecure())
	if err != nil {
		log.Printf("gRPC connection to %v port failed. With: %v\n", os.Getenv("ORDER_GRPC_PORT"), err)
	}

	scooterClient := proto.NewScooterServiceClient(conn)

}
