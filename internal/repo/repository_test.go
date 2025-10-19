package repo

import (
	"L0-wb/internal/models"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLastOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostgresRepo{DB: db}
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		orderRows := sqlmock.NewRows([]string{
			"order_uid", "track_number", "entry", "delivery_id", "payment_id",
			"locale", "internal_signature", "customer_id", "delivery_service",
			"shardkey", "sm_id", "date_created", "oof_shard",
		}).AddRow(
			"test-123", "track1", "WBIL", 1, 1,
			"en", "sig1", "customer1", "test",
			"1", 1, time.Now(), "1",
		)

		deliveryRows := sqlmock.NewRows([]string{
			"name", "phone", "zip", "city", "address", "region", "email",
		}).AddRow("Test User", "+7999999999", "123456", "City", "Address", "Region", "test@test.com")

		paymentRows := sqlmock.NewRows([]string{
			"transaction", "request_id", "currency", "provider", "amount",
			"payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee",
		}).AddRow(
			"tx-1", "req-1", "USD", "stripe", 100,
			time.Now().Unix(), "bank1", 10, 90, 0,
		)

		itemRows := sqlmock.NewRows([]string{
			"chrt_id", "track_number", "price", "rid", "name",
			"sale", "size", "total_price", "nm_id", "brand", "status",
		}).AddRow(
			1, "track1", 100, "rid1", "Item 1",
			0, "M", 100, 1, "Brand", 200,
		)

		mock.ExpectQuery("SELECT (.+) FROM orders").WillReturnRows(orderRows)
		mock.ExpectQuery("SELECT (.+) FROM delivery").WillReturnRows(deliveryRows)
		mock.ExpectQuery("SELECT (.+) FROM payment").WillReturnRows(paymentRows)
		mock.ExpectQuery("SELECT (.+) FROM item").WillReturnRows(itemRows)

		orders, err := repo.GetLastOrders(ctx, 10)
		assert.NoError(t, err)
		assert.Len(t, orders, 1)
		assert.Equal(t, "test-123", orders[0].OrderUID)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM orders").WillReturnError(sqlmock.ErrCancelled)

		orders, err := repo.GetLastOrders(ctx, 10)
		assert.Error(t, err)
		assert.Nil(t, orders)
	})
}

func TestCreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostgresRepo{DB: db}
	ctx := context.Background()

	order := models.Order{
		OrderUID:    "test-123",
		TrackNumber: "track1",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:  "Test User",
			Phone: "+7999999999",
		},
		Payment: models.Payment{
			Transaction: "tx-1",
			Provider:    "stripe",
			Amount:      100,
		},
		Items: []models.Item{
			{
				Name:       "Item 1",
				Price:      100,
				TotalPrice: 100,
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO delivery").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("INSERT INTO payment").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO orders").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery("INSERT INTO item").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.CreateOrder(ctx, order)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		err := repo.CreateOrder(ctx, models.Order{})
		assert.Error(t, err)
	})
}
