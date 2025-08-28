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

type Producer struct {
	writer  *kafka.Writer
	topic   string
	timeout time.Duration
}

func NewProducer(cfg config.Config) *Producer {
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

	// Таймаут можно вынести в конфиг, если нужно
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
			// Генерируем тестовый заказ
			order := p.GenerateTestOrder() // Нужно реализовать эту функцию
			err := p.SendOrders(ctx, order)
			if err != nil {
				logrus.WithError(err).Error("failed to send order")
			}

			// Пауза между отправками
			time.Sleep(1 * time.Second)
		}
	}
}
func (p *Producer) GenerateTestOrder() *models.Order {
	now := time.Now()

	return &models.Order{
		OrderUID:          "test-" + fmt.Sprintf("%d", now.UnixNano()),
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test_customer",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       now,
		OofShard:          "1",

		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},

		Payment: models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},

		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
			{
				ChrtID:      9934931,
				TrackNumber: "WBILMTESTTRACK2",
				Price:       200,
				Rid:         "cd4219087a764ae0btest",
				Name:        "Lipstick",
				Sale:        10,
				Size:        "1",
				TotalPrice:  180,
				NmID:        2389213,
				Brand:       "Maybelline",
				Status:      202,
			},
		},
	}
}
