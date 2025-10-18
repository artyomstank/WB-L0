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
	now := time.Now()
	uid := now.UnixNano() // Уникальные айдишники

	return &models.Order{
		OrderUID:          fmt.Sprintf("test-%d", uid),
		TrackNumber:       fmt.Sprintf("WBTRACK-%d", uid),
		Entry:             "WBIL",
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        fmt.Sprintf("test_customer_%d", uid),
		DeliveryService:   "DHL",
		Shardkey:          "7",
		SmID:              99,
		DateCreated:       now,
		OofShard:          "1",

		Delivery: models.Delivery{
			Name:    "Artem Mayorov",
			Phone:   "+49590445501",
			Zip:     "1235312",
			City:    "Rostov",
			Address: " Marksa 34a",
			Region:  "Tula",
			Email:   fmt.Sprintf("test+%d@gmail.com", uid), // Уникальный email
		},

		Payment: models.Payment{
			Transaction:  fmt.Sprintf("txn-%d", uid), // Уникальный transaction_id
			RequestID:    "",
			Currency:     "USD",
			Provider:     "gazprombank",
			Amount:       1657,
			Bank:         "sigma bank",
			DeliveryCost: 1300,
			GoodsTotal:   357,
			CustomFee:    0,
		},

		Items: []models.Item{
			{
				ChrtID:      4321990,
				TrackNumber: fmt.Sprintf("WBITEMTRACK-%d-1", uid),
				Price:       250,
				Rid:         fmt.Sprintf("rid-%d-1", uid),
				Name:        "Glasses",
				Sale:        25,
				Size:        "0",
				TotalPrice:  187,
				NmID:        9291192,
				Brand:       "Miu miu",
				Status:      202,
			},
			{
				ChrtID:      4321991,
				TrackNumber: fmt.Sprintf("WBITEMTRACK-%d-2", uid),
				Price:       200,
				Rid:         fmt.Sprintf("rid-%d-2", uid),
				Name:        "Cup",
				Sale:        15,
				Size:        "200ml",
				TotalPrice:  170,
				NmID:        1319901,
				Brand:       "Expensive Cups",
				Status:      202,
			},
		},
	}
}
