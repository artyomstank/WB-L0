package service

import (
	"L0-wb/internal/models"
	"context"
)

type Service interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error)
	GetOrderResponse(ctx context.Context, orderUID string) (*models.OrderResponse, error)
	CreateOrder(ctx context.Context, order *models.Order) error
	SaveOrder(ctx context.Context, order *models.Order) error
	RestoreCache(ctx context.Context) error
	Close() error
}
