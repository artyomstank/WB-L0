package kafka

import (
	"L0-wb/config"
	"L0-wb/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Service описывает поведение сервисного слоя,
// куда мы будем сохранять заказ.
type Service interface {
	SaveOrder(ctx context.Context, order *models.Order) error
}

type Consumer struct {
	reader  *kafka.Reader
	topic   string
	timeout time.Duration
	service Service
}

// NewConsumer создаёт Kafka consumer
func NewConsumer(cfg config.Config, service Service) (*Consumer, error) {
	brokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddr},
		Topic:          cfg.Kafka.Topic,
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // авто-коммит каждую секунду
	})
	logrus.WithFields(logrus.Fields{
		"brokers": []string{brokerAddr},
		"topic":   cfg.Kafka.Topic,
	}).Info("Kafka consumer initialized")

	to := 5 * time.Second

	return &Consumer{
		reader:  reader,
		topic:   cfg.Kafka.Topic,
		timeout: to,
		service: service,
	}, nil
}

// Close закрывает consumer
func (c *Consumer) Close() error {
	if c.reader != nil {
		err := c.reader.Close()
		c.reader = nil
		if err != nil {
			logrus.Errorf("ошибка при закрытии consumer: %v", err)
			return err
		}
		logrus.Infof("consumer для топика %s успешно закрыт", c.topic)
	}
	return nil
}

// ConsumeMessages слушает Kafka и обрабатывает сообщения
func (c *Consumer) ConsumeMessages(ctx context.Context) {
	logrus.Infof("Старт чтения сообщений из Kafka (topic: %s)", c.topic)
	for {
		if ctx.Err() != nil {
			logrus.Info("consumer stopped by context")
			return
		}
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				logrus.Info("consumer stopped by context")
				return
			}
			logrus.WithError(err).Error("read message error")
			time.Sleep(time.Second)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			logrus.WithError(err).Errorf("unmarshal order error, raw message: %s", string(m.Value))
			continue
		}

		if order.OrderUID == "" {
			logrus.Errorf("invalid order: missing order_uid. raw message: %s", string(m.Value))
			continue
		}

		// Отправляем заказ в сервисный слой
		if err := c.service.SaveOrder(ctx, &order); err != nil {
			logrus.WithError(err).Errorf("failed to save order %s", order.OrderUID)
		} else {
			logrus.Infof("order %s saved successfully", order.OrderUID)
		}

		logrus.WithFields(logrus.Fields{
			"partition": m.Partition,
			"offset":    m.Offset,
			"order_uid": order.OrderUID,
		}).Info("message processed")
	}
}
