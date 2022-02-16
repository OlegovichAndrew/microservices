package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	proto "scooter_micro/proto"
	"scooter_micro/repository/mock"
	"scooter_micro/service"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
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
	var (
		scooterService *service.ScooterService
		mockCtrl       *gomock.Controller
		repoScooter    *mock.MockScooterRepository
		order          proto.OrderServiceClient
		expectedError  error
		ctx            context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repoScooter = mock.NewMockScooterRepository(mockCtrl)
		expectedError = errors.New("expectedError")
		ctx = context.Background()
		conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		Expect(err).ToNot(HaveOccurred())
		order = proto.NewOrderServiceClient(conn)
		scooterService = &service.ScooterService{Repo: repoScooter, Order: order}
	})

	Describe(".GetAllStations()", func() {
		var stationList *proto.StationList

		Context("mocked in a case of success", func() {
			JustBeforeEach(func() {
				stationList = &proto.StationList{Stations: []*proto.Station{{Id: 1, Name: "station1"},
					{Id: 2, Name: "station2"}}}
				repoScooter.EXPECT().GetAllStations(ctx, &proto.Request{}).Return(stationList, nil)
			})

			It("station list is not nil", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
			It("stations have names", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list.Stations[0].Name).ShouldNot(BeEmpty())
				Expect(list.Stations[1].Name).ShouldNot(BeEmpty())
				Expect(err).Should(BeNil())
			})
			It("stations names are strings", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list.Stations[0].Name).Should(BeAssignableToTypeOf("string"))
				Expect(err).Should(BeNil())
			})
			It("list has stations", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(len(list.Stations)).ShouldNot(BeZero())
				Expect(err).Should(BeNil())
			})
			It("station list should have length = 2", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list.Stations).Should(HaveLen(2))
				Expect(err).Should(BeNil())
			})
			It("station list should have ID higher or equal than 1", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list.Stations[0].Id).Should(HaveValue(BeNumerically(">=", 1)))
				Expect(err).Should(BeNil())
			})
			It("station  .IsActive should be false", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list.Stations[1].IsActive).Should(BeFalse())
				Expect(err).Should(BeNil())
			})
			It("station .latitude should be float", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list.Stations[0].Latitude).Should(BeAssignableToTypeOf(1.5))
				Expect(err).Should(BeNil())
			})
			It("station .latitude should be lesser than 5", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list.Stations[0].Latitude).Should(BeNumerically("<", 5))
				Expect(err).Should(BeNil())
			})
			It("the first station HAS field .latitude which is lesser or equal to zero", func() {
				list, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(list.Stations[0]).Should(HaveField("Latitude", BeNumerically("<=", 0)))
				Expect(err).Should(BeNil())
			})

		})

		Context("mocked with an expected error", func() {
			JustBeforeEach(func() {
				stationList = &proto.StationList{Stations: []*proto.Station{{Id: 1, Name: "station1"},
					{Id: 2, Name: "station2"}}}

				repoScooter.EXPECT().GetAllStations(ctx, &proto.Request{}).Return(stationList, expectedError)
			})

			It("error is not nil", func() {
				_, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(err).ShouldNot(BeNil())
			})
			It("error matches with expectedError", func() {
				_, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(err).Should(MatchError(expectedError))
			})
			It("should have occurred", func() {
				_, err := scooterService.GetAllStations(ctx, &proto.Request{})
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe(".GetAllScooters()", func() {
		var scooterList *proto.ScooterList

		Context("mocked in a case of success", func() {
			JustBeforeEach(func() {
				scooterList = &proto.ScooterList{Scooters: []*proto.Scooter{{Id: 1, ScooterModel: "model1"},
					{Id: 2, ScooterModel: "model2"}}}

				repoScooter.EXPECT().GetAllScooters(ctx, &proto.Request{}).Return(scooterList, nil).AnyTimes()
			})

			It("list is equal to mocked one", func() {
				Expect(scooterService.GetAllScooters(ctx, &proto.Request{})).To(Equal(scooterList))
				Expect(scooterService.GetAllScooters(ctx, &proto.Request{})).Error().Should(BeNil())
			})
			It("list length longer than 1", func() {
				list, err := scooterService.GetAllScooters(ctx, &proto.Request{})
				Expect(len(list.Scooters)).To(BeNumerically(">=", 1))
				Expect(err).To(BeNil())
			})
			It("ID's higher than 0", func() {
				list, err := scooterService.GetAllScooters(ctx, &proto.Request{})
				Expect(list.Scooters[1].Id).To(BeNumerically(">=", 1))
				Expect(list.Scooters[0].Id).To(BeNumerically(">=", 1))
				Expect(err).To(BeNil())
			})
		})

		Context("mocked with an error", func() {
			JustBeforeEach(func() {
				repoScooter.EXPECT().GetAllScooters(ctx, &proto.Request{}).Return(&proto.ScooterList{}, expectedError).AnyTimes()
			})
			It("returns an empty list and an error", func() {
				result, err := scooterService.GetAllScooters(ctx, &proto.Request{})
				Expect(result).To(Equal(&proto.ScooterList{}))
				Expect(err).NotTo(BeNil())
				Expect(err).Should(MatchError(expectedError))
			})
		})
	})

	Describe(".GetScooterById()", func() {
		Context("mocked in a case of success", func() {
			JustBeforeEach(func() {
				repoScooter.EXPECT().GetScooterById(ctx, &proto.ScooterID{Id: 1}).Return(&proto.Scooter{Id: 1, MaxWeight: 120.0,
					ScooterModel: "Xiaomi", BatteryRemain: 50, CanBeRent: true, StationID: 3}, nil)
			})

			It("scooter ID = 1", func() {
				result, err := scooterService.GetScooterById(ctx, &proto.ScooterID{Id: 1})
				Expect(result.Id).Should(Equal(uint64(1)))
				Expect(err).Should(BeNil())
			})
			It("scooter .CanBeRent should be true", func() {
				result, err := scooterService.GetScooterById(ctx, &proto.ScooterID{Id: 1})
				Expect(result.CanBeRent).Should(BeTrue())
				Expect(err).Should(BeNil())
			})
			It("scooter .BatteryRemain should be > 20", func() {
				result, err := scooterService.GetScooterById(ctx, &proto.ScooterID{Id: 1})
				Expect(result.BatteryRemain).Should(BeNumerically(">=", 50))
				Expect(err).Should(BeNil())
			})
			It("returned value should be similar to type .Scooter", func() {
				result, err := scooterService.GetScooterById(ctx, &proto.ScooterID{Id: 1})
				Expect(result).Should(BeAssignableToTypeOf(&proto.Scooter{}))
				Expect(err).Should(BeNil())
			})
			It("scooter model should be 'Xiaomi' ", func() {
				result, err := scooterService.GetScooterById(ctx, &proto.ScooterID{Id: 1})
				Expect(result.ScooterModel).Should(BeEquivalentTo("Xiaomi"))
				Expect(err).Should(BeNil())
			})
		})

		Context("mocked with an error", func() {
			JustBeforeEach(func() {
				repoScooter.EXPECT().GetScooterById(ctx, &proto.ScooterID{Id: 1}).Return(&proto.Scooter{Id: 1, MaxWeight: 120.0,
					ScooterModel: "Xiaomi", BatteryRemain: 50, CanBeRent: true, StationID: 3}, expectedError)
			})

			It("error not nil", func() {
				_, err := scooterService.GetScooterById(ctx, &proto.ScooterID{Id: 1})
				Expect(err).ShouldNot(BeNil())
			})
			It("error matches with expectedError", func() {
				_, err := scooterService.GetScooterById(ctx, &proto.ScooterID{Id: 1})
				Expect(err).Should(MatchError(expectedError))
			})
			It("should have occurred", func() {
				_, err := scooterService.GetScooterById(ctx, &proto.ScooterID{Id: 1})
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe(".SendCurrentStatus()", func() {
		var sendStatus *proto.SendStatus

		Context("mocked in a case of success", func() {
			JustBeforeEach(func() {
				sendStatus = &proto.SendStatus{StationID: 1, ScooterID: 2, Latitude: 54.78, Longitude: 47.66, BatteryRemain: 40.0}
				repoScooter.EXPECT().SendCurrentStatus(ctx, sendStatus).Return(nil).AnyTimes()
			})

			It("with a nil error", func() {
				Expect(scooterService.SendCurrentStatus(ctx, sendStatus)).ShouldNot(HaveOccurred())
			})
		})

		Context("mocked with an error", func() {
			JustBeforeEach(func() {
				sendStatus = &proto.SendStatus{StationID: 1, ScooterID: 2, Latitude: 54.78, Longitude: 47.66, BatteryRemain: 40.0}
				repoScooter.EXPECT().SendCurrentStatus(ctx, sendStatus).Return(expectedError).AnyTimes()
			})

			It("error occurred", func() {
				Expect(scooterService.SendCurrentStatus(ctx, sendStatus)).Should(HaveOccurred())
			})
			It("error matches with expectedError", func() {
				Expect(scooterService.SendCurrentStatus(ctx, sendStatus)).Error().Should(MatchError(expectedError))
			})
		})
	})

	Describe(".GetScooterStatus", func() {
		var scooterStatus *proto.ScooterStatus

		Context("mocked in a case of success", func() {
			JustBeforeEach(func() {
				scooterStatus = &proto.ScooterStatus{Latitude: 55.5, Longitude: 48.8, StationID: &proto.StationID{Id: 432}, BatteryRemain: 34.2}
				repoScooter.EXPECT().GetScooterStatus(ctx, &proto.ScooterID{Id: 3}).Return(scooterStatus, nil).AnyTimes()
			})

			It("returns a status which is equal to expected", func() {
				status, err := scooterService.GetScooterStatus(ctx, &proto.ScooterID{Id: 3})
				Expect(status).Should(BeEquivalentTo(scooterStatus))
				Expect(err).Should(BeNil())
			})
			It("status should have a StationID field witth ID:432", func() {
				status, err := scooterService.GetScooterStatus(ctx, &proto.ScooterID{Id: 3})
				Expect(status).Should(HaveField("StationID", &proto.StationID{Id: 432}))
				Expect(err).Should(BeNil())
			})
			It("should have battery remain > 20", func() {
				status, err := scooterService.GetScooterStatus(ctx, &proto.ScooterID{Id: 3})
				Expect(status).Should(HaveField("BatteryRemain", BeNumerically(">", 20)))
				Expect(err).Should(BeNil())
			})
			It("should have a field like in base scooterStatus", func() {
				status, err := scooterService.GetScooterStatus(ctx, &proto.ScooterID{Id: 3})
				Expect(status).Should(HaveField("BatteryRemain", BeElementOf(scooterStatus.BatteryRemain)))
				Expect(err).Should(BeNil())
			})
		})

		Context("mocked with an error", func() {
			JustBeforeEach(func() {
				scooterStatus = &proto.ScooterStatus{Latitude: 55.5, Longitude: 48.8, StationID: &proto.StationID{Id: 432}, BatteryRemain: 34.2}
				repoScooter.EXPECT().GetScooterStatus(ctx, &proto.ScooterID{Id: 3}).Return(scooterStatus, expectedError).AnyTimes()
			})

			It("error not nil", func() {
				_, err := scooterService.GetScooterStatus(ctx, &proto.ScooterID{Id: 3})
				Expect(err).ShouldNot(BeNil())
			})
			It("error matches with expectedError", func() {
				_, err := scooterService.GetScooterStatus(ctx, &proto.ScooterID{Id: 3})
				Expect(err).Should(MatchError(expectedError))
			})
			It("should have occurred", func() {
				_, err := scooterService.GetScooterStatus(ctx, &proto.ScooterID{Id: 3})
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe(".CreateScooterStatusInRent()", func() {
		var statusInRent *proto.ScooterStatusInRent

		Context("mocked in a case of success", func() {
			JustBeforeEach(func() {
				statusInRent = &proto.ScooterStatusInRent{ScooterID: 33, Id: 2, StationID: 43, DateTime: &timestamppb.Timestamp{Nanos: int32(432432454)}}
				repoScooter.EXPECT().CreateScooterStatusInRent(ctx, &proto.ScooterID{Id: 33}).Return(statusInRent, nil).AnyTimes()
			})

			It("has been created with nil error", func() {
				rentStatus, err := scooterService.CreateScooterStatusInRent(ctx, &proto.ScooterID{Id: 33})
				testFunc := func(id1 uint64) uint64 { return id1 + 3 }
				Expect(rentStatus.ScooterID).To(WithTransform(testFunc, Equal(uint64(36))))
				Expect(err).Should(BeNil())
			})
			It("DateTime should be of type timestambpb", func() {
				rentStatus, err := scooterService.CreateScooterStatusInRent(ctx, &proto.ScooterID{Id: 33})
				Expect(rentStatus.DateTime).Should(BeAssignableToTypeOf(timestamppb.Now()))
				Expect(err).Should(BeNil())
			})
			It("should math with JSON", func() {
				rentStatus, err := scooterService.CreateScooterStatusInRent(ctx, &proto.ScooterID{Id: 33})
				statusToJson, _ := json.Marshal(statusInRent)
				finalStatusToJson, _ := json.Marshal(rentStatus)

				Expect(statusToJson).Should(MatchJSON(finalStatusToJson))
				Expect(err).Should(BeNil())
			})
		})

		Context("mocked with an error", func() {
			JustBeforeEach(func() {
				statusInRent = &proto.ScooterStatusInRent{ScooterID: 33, Id: 2, StationID: 43, DateTime: &timestamppb.Timestamp{Nanos: int32(432432454)}}
				repoScooter.EXPECT().CreateScooterStatusInRent(ctx, &proto.ScooterID{Id: 33}).Return(statusInRent, expectedError).AnyTimes()
			})

			It("error not nil", func() {
				_, err := scooterService.CreateScooterStatusInRent(ctx, &proto.ScooterID{Id: 33})
				Expect(err).ShouldNot(BeNil())
			})
			It("error matches with expectedError", func() {
				_, err := scooterService.CreateScooterStatusInRent(ctx, &proto.ScooterID{Id: 33})
				Expect(err).Should(MatchError(expectedError))
			})
			It("should have occurred", func() {
				_, err := scooterService.CreateScooterStatusInRent(ctx, &proto.ScooterID{Id: 33})
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
