package service

import (
	"context"
	"order_micro/proto"
	"order_micro/repository"
)

type OrderInterface interface {
	CreateOrder(ctx context.Context, info *proto.TripInfo) (*proto.Order, error)
}

type OrderService struct {
	Repo *repository.OrderRepo
	*proto.UnimplementedOrderServiceServer
}

func NewOrderService(repo *repository.OrderRepo) *OrderService {
	return &OrderService{Repo: repo}
}

func (os *OrderService) CreateOrder(ctx context.Context, info *proto.TripInfo) (*proto.Order, error) {
	return os.Repo.CreateOrder(ctx, info)
}