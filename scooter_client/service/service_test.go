package service_test

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"scooter_client/model"
	"scooter_client/proto"
	"scooter_client/service"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var serv proto.ScooterServiceServer

func init() {
	serv = proto.UnimplementedScooterServiceServer{}
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	proto.RegisterScooterServiceServer(s, serv)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

var _ = Describe("Service", func() {

	var protoScooterClient proto.ScooterServiceClient
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		Expect(err).To(BeNil())

		protoScooterClient = proto.NewScooterServiceClient(conn)
	})

	Describe(".gRPC scooter message()", func() {
		var scooterCl *service.ScooterClient
		var stream proto.ScooterService_RegisterClient
		var err error
		var expecterError error

		Context("correct sending by stream with correct data", func() {
			BeforeEach(func() {
				stream, err = protoScooterClient.Register(ctx)
				Expect(err).Should(BeNil())
				scooterCl = &service.ScooterClient{ID: 7, Latitude: 52.18, Longitude: 47.17, BatteryRemain: 52.7, Stream: stream}
			})
			It("return nil error", func() {
				Expect(scooterCl.GrpcScooterMessage()).Should(Succeed())
				Expect(scooterCl.GrpcScooterMessage()).Error().Should(BeNil())
			})
		})

		Context("with an error", func() {
			BeforeEach(func() {
				expecterError = errors.New("ID shouldn't be zero or below")
				stream, err = protoScooterClient.Register(ctx)
				Expect(err).Should(BeNil())
				scooterCl = &service.ScooterClient{Latitude: 52.18, Longitude: 47.17, BatteryRemain: 52.7, Stream: stream}
			})
			It("when scooter ID = 0", func() {
				Expect(scooterCl.GrpcScooterMessage()).Should(HaveOccurred())
				Expect(scooterCl.GrpcScooterMessage()).Error().Should(MatchError(expecterError))
			})
		})
	})

	Describe(".Run()", func() {
		var scooterCl *service.ScooterClient
		var stream proto.ScooterService_RegisterClient
		var station model.Location
		var scooterStartState *service.ScooterClient
		var err error

		Context("in a case of success", func() {
			BeforeEach(func() {
				station = model.Location{Longitude: 59.0, Latitude: 49.0}
				stream, err = protoScooterClient.Register(ctx)
				Expect(err).Should(BeNil())
				scooterCl = &service.ScooterClient{ID: 7, Latitude: 52.18, Longitude: 47.17, BatteryRemain: 52.7, Stream: stream}
				scooterStartState = scooterCl
			})

			It("with all parameters", func() {
				result, err := scooterCl.Run(station)
				Expect(result).Should(Not(BeNil()))
				Expect(err).Should(BeNil())
			})
			It("should return a scooter status which coordinates are close to station", func() {
				result, err := scooterCl.Run(station)
				Expect(result.Latitude).Should(BeNumerically(">=", station.Latitude))
				Expect(result.Longitude).Should(BeNumerically("<=", station.Longitude))
				Expect(err).Should(BeNil())
			})

			It("battery remain should be lesser than at start", func() {
				result, err := scooterCl.Run(station)
				Expect(result.BatteryRemain).Should(BeNumerically("<=", scooterStartState.BatteryRemain))
				Expect(err).Should(BeNil())
			})
			It("ID after finish should be 0", func() {
				_, err := scooterCl.Run(station)
				Expect(scooterCl.ID).Should(BeZero())
				Expect(err).Should(BeNil())
			})
		})
	})
})
