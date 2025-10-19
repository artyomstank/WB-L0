package kafka

import (
	"L0-wb/config"
	"L0-wb/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"L0-wb/internal/generator"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	writer  *kafka.Writer
	topic   string
	timeout time.Duration
}

func NewProducer(cfg config.Config) ProducerInterface {
	// Собираем адрес брокера из Host + Port
	brokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddr},
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.Hash{},
		Async:    false,
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			logrus.Errorf(msg, args...)
		}),
		RequiredAcks: int(kafka.RequireAll),
	})

	to := 5 * time.Second

	return &Producer{
		writer:  writer,
		topic:   cfg.Kafka.Topic,
		timeout: to,
	}
}

func (p *Producer) Close() error {
	if p.writer != nil {
		err := p.writer.Close()
		p.writer = nil
		return err
	}
	return nil
}

func (p *Producer) SendOrders(ctx context.Context, order *models.Order) error {
	if p.writer == nil {
		return fmt.Errorf("writer is nil")
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}

	value, err := json.Marshal(order)
	if err != nil {
		logrus.WithError(err).Error("marshal order error")
		return err
	}

	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: value,
		Time:  time.Now(),
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if err := p.writer.WriteMessages(ctxTimeout, msg); err != nil {
		logrus.WithError(err).Error("send message error")
		return err
	}

	logrus.WithField("order_uid", order.OrderUID).Info("message sent to kafka")
	return nil
}
func (p *Producer) RunProducer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			order := p.GenerateTestOrder()
			err := p.SendOrders(ctx, order)
			if err != nil {
				logrus.WithError(err).Error("failed to send order")
			}

			// Пауза между отправками
			time.Sleep(2 * time.Second)
		}
	}
}

func (p *Producer) GenerateTestOrder() *models.Order {
	return generator.GenerateOrder()
}
