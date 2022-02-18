package service

//go:generate mockgen -destination=./mock/mock_scooterRepository.go -package=mock -source=./service.go

import (
	"context"
	"scooter_micro/proto"
	"scooter_micro/repository"
)

const (
	step          = 0.0001
	dischargeStep = 0.1
	interval      = 450
)

//ScooterService is a service which responsible for gRPC scooter.
type ScooterService struct {
	Repo  repository.ScooterRepository
	Order proto.OrderServiceClient
	*proto.UnimplementedScooterServiceServer
}

//NewScooterService creates a new GrpcScooterService.
func NewScooterService(repoScooter repository.ScooterRepository, order proto.OrderServiceClient) *ScooterService {
	return &ScooterService{
		Repo:  repoScooter,
		Order: order,
	}
}

//GetAllScooters gives the access to the ScooterRepo.GetAllScooters function.
func (gss *ScooterService) GetAllScooters(ctx context.Context, request *proto.Request) (*proto.ScooterList, error) {
	return gss.Repo.GetAllScooters(ctx, request)
}

func (gss *ScooterService) GetAllScootersByStationID(ctx context.Context, id *proto.StationID) (*proto.ScooterList,
	error) {
	return gss.Repo.GetAllScootersByStationID(ctx, id)
}

func (gss *ScooterService) GetAllStations(ctx context.Context, request *proto.Request) (*proto.StationList,
	error) {
	return gss.Repo.GetAllStations(ctx, request)
}

//GetScooterById gives the access to the ScooterRepo.GetScooterById function.
func (gss *ScooterService) GetScooterById(ctx context.Context, id *proto.ScooterID) (*proto.Scooter, error) {
	return gss.Repo.GetScooterById(ctx, id)
}

func (gss *ScooterService) GetStationById(ctx context.Context, id *proto.StationID) (*proto.Station, error) {
	return gss.Repo.GetStationById(ctx, id)
}

//GetScooterStatus gives the access to the ScooterRepo.GetScooterStatus function.
func (gss *ScooterService) GetScooterStatus(ctx context.Context, status *proto.ScooterID) (*proto.ScooterStatus, error) {
	return gss.Repo.GetScooterStatus(ctx, status)
}

//SendCurrentStatus gives the access to the ScooterRepo.SendCurrentStatus function.
func (gss *ScooterService) SendCurrentStatus(ctx context.Context, status *proto.SendStatus) error {
	return gss.Repo.SendCurrentStatus(ctx, status)
}

//CreateScooterStatusInRent gives the access to the ScooterRepo.CreateScooterStatusInRent function.
func (gss *ScooterService) CreateScooterStatusInRent(ctx context.Context, id *proto.ScooterID) (*proto.ScooterStatusInRent,
	error) {
	return gss.Repo.CreateScooterStatusInRent(ctx, id)
}
