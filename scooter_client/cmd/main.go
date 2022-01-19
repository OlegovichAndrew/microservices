package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"scooter_client/model"
	"scooter_client/proto"
	"scooter_client/service"
	"time"
)

var destination model.Location

func main() {
	conn, err := grpc.DialContext(context.Background(), net.JoinHostPort("", os.Getenv("GRPC_PORT")),
		grpc.WithInsecure())
	if err != nil {
		log.Printf("gRPC connection to %v port failed. With: %v\n", os.Getenv("GRPC_PORT"), err)
	}

	log.Printf("gRPC connected port: %v.", os.Getenv("GRPC_PORT"))

	client := proto.NewScooterServiceClient(conn)
	stream, err := client.Register(context.Background())
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	ctx := stream.Context()
	done := make(chan bool)

	scooterClient := service.NewScooterClient(0, 0.0, 0.0, 0.0, stream)

	fmt.Printf("This is a scooter client: %v\n", scooterClient)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				fmt.Println(err)
				return
			}

			fmt.Printf("!Received from server: %v\n", resp)

			if err != nil {
				log.Fatalf("can not receive %v", err)
			}

			destination.Latitude = resp.DestLatitude
			destination.Longitude = resp.DestLongitude

			scooterClient.Longitude = resp.Longitude
			scooterClient.Latitude = resp.Latitude
			scooterClient.BatteryRemain = resp.BatteryRemain
			scooterClient.ID = resp.Id

			fmt.Printf("Scooter client is:%v\n", scooterClient)
			fmt.Printf("Destination is:%v\n", scooterClient)

		}
	}()

	go func() {
		for {
			if scooterClient.ID != 0 {
				err = scooterClient.Run(destination)
				if err != nil {
					fmt.Println(err)
				}
			}

			msg := &proto.ClientMessage{
				Id:        scooterClient.ID,
				Latitude:  scooterClient.Latitude,
				Longitude: scooterClient.Longitude,
			}

			fmt.Printf("Send to server this message: %v\n", msg)
			err := scooterClient.Stream.Send(msg)
			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(time.Second * 3)

		}
	}()

	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
		close(done)
	}()

	<-done

}
