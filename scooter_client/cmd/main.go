package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"io"
	"log"
	"scooter_client/config"
	"scooter_client/model"
	"scooter_client/proto"
	"scooter_client/service"
	"scooter_client/transport"
	"time"
)

var destination model.Location

const ClientID = "some_client"
const TopicName = "order"


func main() {
	conn, err := grpc.Dial(config.SERVER_CONN_GRPC_ADDRESS, grpc.WithInsecure() )
	if err != nil {
		log.Printf("gRPC connection to %v port failed. With: %v\n", config.GRPC_PORT, err)
	}


	log.Printf("gRPC connected port: %v.", config.GRPC_PORT)

	client := proto.NewScooterServiceClient(conn)
	stream, err := client.Register(context.Background())

	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	ctx := stream.Context()
	done := make(chan bool)
	var currentStationID uint64
	scooterClient := service.NewScooterClient(0, 0.0, 0.0, 0.0, stream)

	producer := transport.CreateProducer([]string{config.KAFKA_BROKER}, ClientID)

	err = transport.CreateTopic([]string{config.KAFKA_BROKER}, TopicName, 1, 1)
	if err != nil && err.(*sarama.TopicError).Err != sarama.ErrTopicAlreadyExists {
		log.Fatalln("Failed to create kafka topic:", err)
		return
	}


	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				fmt.Println(err)
				return
			}

			fmt.Printf("Received from server: %v\n", resp)

			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			currentStationID = uint64(resp.StationID)
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
			//If I got data from the server, I will start scooter moving.
			if scooterClient.ID != 0 {
				currentStatus, err := scooterClient.Run(destination)
				if err != nil {
					fmt.Println(err)
				}
				currentStatus.StationID = currentStationID

				fmt.Println(currentStatus)

				//Remote call for server's method.
				_ ,err = client.SendCurrentStatus(ctx, currentStatus )
				if err != nil {
					fmt.Println(err)
				}

				msg, err := json.Marshal(currentStatus)
				if err != nil {fmt.Println(err)}
				err = transport.SendMessage(producer, TopicName, string(msg))
				if err != nil {fmt.Println(err)}
			}
			// a mock message for keeping the stream.
			msg := &proto.ClientMessage{
				Id:        scooterClient.ID,
				Latitude:  scooterClient.Latitude,
				Longitude: scooterClient.Longitude,
			}

			fmt.Printf("Sent to server this message: %v\n", msg)
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

