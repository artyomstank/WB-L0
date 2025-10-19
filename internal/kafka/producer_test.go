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

func TestProducer_SendOrders(t *testing.T) {
	cfg := config.Config{
		Kafka: config.Kafka{
			Host:  "localhost",
			Port:  9092,
			Topic: "test-topic",
		},
	}

	producer := NewProducer(cfg)
	defer producer.Close()

	testOrder := &models.Order{
		OrderUID:    "test-123",
		TrackNumber: "WBIL123456789",
		DateCreated: time.Now(),
	}

	ctx := context.Background()
	err := producer.SendOrders(ctx, testOrder)
	assert.NoError(t, err)

	// Test nil writer
	p := &Producer{writer: nil}
	err = p.SendOrders(ctx, testOrder)
	assert.Error(t, err)
}
