package kafka

import (
	"L0-wb/internal/models"
	"context"
)

type ProducerInterface interface {
	SendOrders(ctx context.Context, order *models.Order) error
	RunProducer(ctx context.Context) error
	Close() error
	GenerateTestOrder() *models.Order
}

type ConsumerInterface interface {
	ConsumeMessages(ctx context.Context) error
	Close() error
}

type MessageProcessor interface {
	SaveOrder(ctx context.Context, order *models.Order) error
}
