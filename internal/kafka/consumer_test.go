//go:build integration

package kafka

import (
	"L0-wb/config"
	"L0-wb/internal/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockService struct {
	savedOrders []*models.Order
}

func (m *mockService) SaveOrder(_ context.Context, order *models.Order) error {
	m.savedOrders = append(m.savedOrders, order)
	return nil
}

func TestConsumer_ConsumeMessages(t *testing.T) {
	cfg := config.Config{
		Kafka: config.Kafka{
			Host:  "localhost",
			Port:  9092,
			Topic: "test-topic",
		},
	}

	mockSvc := &mockService{}
	consumer, err := NewConsumer(cfg, mockSvc)
	assert.NoError(t, err)
	defer consumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = consumer.ConsumeMessages(ctx)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
