//go:build unit
// +build unit

package repo

import (
	"L0-wb/internal/models"
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (*PostgresRepo, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return &PostgresRepo{DB: db}, mock, func() { db.Close() }
}

func TestCreateDeliveryTx(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()
	ctx := context.Background()

	mock.ExpectBegin()
	tx, err := repo.DB.Begin()
	require.NoError(t, err)

	// validation error
	_, err = repo.CreateDeliveryTx(ctx, tx, models.Delivery{Name: "", Phone: "123"})
	require.Error(t, err)
	_, err = repo.CreateDeliveryTx(ctx, tx, models.Delivery{Name: "Ivan", Phone: ""})
	require.Error(t, err)

	// db error
	mock.ExpectQuery("INSERT INTO delivery").
		WithArgs("Ivan", "123", "", "", "", "", "").
		WillReturnError(errors.New("db error"))
	_, err = repo.CreateDeliveryTx(ctx, tx, models.Delivery{Name: "Ivan", Phone: "123"})
	require.Error(t, err)

	// success
	rows := sqlmock.NewRows([]string{"id"}).AddRow(42)
	mock.ExpectQuery("INSERT INTO delivery").
		WithArgs("Ivan", "123", "", "", "", "", "").
		WillReturnRows(rows)
	id, err := repo.CreateDeliveryTx(ctx, tx, models.Delivery{Name: "Ivan", Phone: "123"})
	require.NoError(t, err)
	require.Equal(t, 42, id)
}

func TestCreatePaymentTx(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()
	ctx := context.Background()

	mock.ExpectBegin()
	tx, err := repo.DB.Begin()
	require.NoError(t, err)

	// validation errors
	_, err = repo.CreatePaymentTx(ctx, tx, models.Payment{Transaction: "", Provider: "p", Amount: 1})
	require.Error(t, err)
	_, err = repo.CreatePaymentTx(ctx, tx, models.Payment{Transaction: "t", Provider: "", Amount: 1})
	require.Error(t, err)
	_, err = repo.CreatePaymentTx(ctx, tx, models.Payment{Transaction: "t", Provider: "p", Amount: 0})
	require.Error(t, err)

	// db error
	mock.ExpectQuery("INSERT INTO payment").
		WithArgs("t", "", "", "p", 1, 0, "", 0, 0, 0).
		WillReturnError(errors.New("db error"))
	_, err = repo.CreatePaymentTx(ctx, tx, models.Payment{Transaction: "t", Provider: "p", Amount: 1})
	require.Error(t, err)

	// success
	rows := sqlmock.NewRows([]string{"id"}).AddRow(7)
	mock.ExpectQuery("INSERT INTO payment").
		WithArgs("t", "", "", "p", 1, 0, "", 0, 0, 0).
		WillReturnRows(rows)
	id, err := repo.CreatePaymentTx(ctx, tx, models.Payment{Transaction: "t", Provider: "p", Amount: 1})
	require.NoError(t, err)
	require.Equal(t, 7, id)
}

func TestCreateItemTx(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()
	ctx := context.Background()

	mock.ExpectBegin()
	tx, err := repo.DB.Begin()
	require.NoError(t, err)

	// validation errors
	_, err = repo.CreateItemTx(ctx, tx, models.Item{Name: "", Price: 1}, "uid")
	require.Error(t, err)
	_, err = repo.CreateItemTx(ctx, tx, models.Item{Name: "n", Price: 0}, "uid")
	require.Error(t, err)
	_, err = repo.CreateItemTx(ctx, tx, models.Item{Name: "n", Price: 1}, "")
	require.Error(t, err)

	// db error
	mock.ExpectQuery("INSERT INTO item").
		WithArgs(0, "", 1, "", "n", 0, "", 0, 0, "", 0, "uid").
		WillReturnError(errors.New("db error"))
	_, err = repo.CreateItemTx(ctx, tx, models.Item{Name: "n", Price: 1}, "uid")
	require.Error(t, err)

	// success
	rows := sqlmock.NewRows([]string{"id"}).AddRow(5)
	mock.ExpectQuery("INSERT INTO item").
		WithArgs(0, "", 1, "", "n", 0, "", 0, 0, "", 0, "uid").
		WillReturnRows(rows)
	id, err := repo.CreateItemTx(ctx, tx, models.Item{Name: "n", Price: 1}, "uid")
	require.NoError(t, err)
	require.Equal(t, 5, id)
}

func TestCreateOrder(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()
	ctx := context.Background()

	order := models.Order{
		OrderUID:          "uid",
		TrackNumber:       "track",
		Entry:             "entry",
		Delivery:          models.Delivery{Name: "Ivan", Phone: "123"},
		Payment:           models.Payment{Transaction: "t", Provider: "p", Amount: 1},
		Locale:            "ru",
		InternalSignature: "sig",
		CustomerID:        "cid",
		DeliveryService:   "svc",
		Shardkey:          "shard",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "oof",
		Items: []models.Item{
			{Name: "item1", Price: 10},
			{Name: "item2", Price: 20},
		},
	}

	// validation error
	err := repo.CreateOrder(ctx, models.Order{})
	require.Error(t, err)

	// mock transaction begin
	mock.ExpectBegin()

	// delivery insert
	rows1 := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO delivery (name, phone, zip, city, address, region, email) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id")).
		WithArgs("Ivan", "123", "", "", "", "", "").
		WillReturnRows(rows1)

	// payment insert
	rows2 := sqlmock.NewRows([]string{"id"}).AddRow(2)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id")).
		WithArgs("t", "", "", "p", 1, 0, "", 0, 0, 0).
		WillReturnRows(rows2)

	// order insert (fix: используем простую строку без спецсимволов)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO orders (")).
		WithArgs(order.OrderUID, order.TrackNumber, order.Entry, 1, 2, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// items insert
	rows3 := sqlmock.NewRows([]string{"id"}).AddRow(11)
	mock.ExpectQuery("INSERT INTO item").
		WithArgs(0, "", 10, "", "item1", 0, "", 0, 0, "", 0, "uid").
		WillReturnRows(rows3)
	rows4 := sqlmock.NewRows([]string{"id"}).AddRow(12)
	mock.ExpectQuery("INSERT INTO item").
		WithArgs(0, "", 20, "", "item2", 0, "", 0, 0, "", 0, "uid").
		WillReturnRows(rows4)

	// commit
	mock.ExpectCommit()

	err = repo.CreateOrder(ctx, order)
	require.NoError(t, err)
}

func TestGetOrder(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()
	ctx := context.Background()

	orderUID := "uid"
	deliveryID := 1
	paymentID := 2
	dateCreated := time.Now()

	// validation error
	_, err := repo.GetOrder(ctx, "")
	require.Error(t, err)

	// order row
	orderRows := sqlmock.NewRows([]string{
		"order_uid", "track_number", "entry", "delivery_id", "payment_id", "locale", "internal_signature",
		"customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard",
	}).AddRow(
		orderUID, "track", "entry", deliveryID, paymentID, "ru", "sig",
		"cid", "svc", "shard", 1, dateCreated, "oof",
	)
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1",
	)).WithArgs(orderUID).WillReturnRows(orderRows)

	// delivery row
	deliveryRows := sqlmock.NewRows([]string{
		"name", "phone", "zip", "city", "address", "region", "email",
	}).AddRow("Ivan", "123", "", "", "", "", "")
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT name, phone, zip, city, address, region, email FROM delivery WHERE id = $1",
	)).WithArgs(deliveryID).WillReturnRows(deliveryRows)

	// payment row
	paymentRows := sqlmock.NewRows([]string{
		"transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee",
	}).AddRow("t", "", "", "p", 1, 0, "", 0, 0, 0)
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE id = $1",
	)).WithArgs(paymentID).WillReturnRows(paymentRows)

	// items rows
	itemsRows := sqlmock.NewRows([]string{
		"chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status",
	}).AddRow(0, "", 10, "", "item1", 0, "", 0, 0, "", 0).
		AddRow(0, "", 20, "", "item2", 0, "", 0, 0, "", 0)
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM item WHERE order_uid = $1",
	)).WithArgs(orderUID).WillReturnRows(itemsRows)

	order, err := repo.GetOrder(ctx, orderUID)
	require.NoError(t, err)
	require.Equal(t, orderUID, order.OrderUID)
	require.Equal(t, "Ivan", order.Delivery.Name)
	require.Equal(t, "t", order.Payment.Transaction)
	require.Len(t, order.Items, 2)
	require.Equal(t, "item1", order.Items[0].Name)
	require.Equal(t, "item2", order.Items[1].Name)
}

func TestGetDelivery(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()
	ctx := context.Background()

	deliveryID := 1

	// db error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE id = $1")).
		WithArgs(deliveryID).
		WillReturnError(errors.New("db error"))
	_, err := repo.GetDelivery(ctx, deliveryID)
	require.Error(t, err)

	// success
	rows := sqlmock.NewRows([]string{
		"name", "phone", "zip", "city", "address", "region", "email",
	}).AddRow("Ivan", "123", "", "", "", "", "")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE id = $1")).
		WithArgs(deliveryID).
		WillReturnRows(rows)
	delivery, err := repo.GetDelivery(ctx, deliveryID)
	require.NoError(t, err)
	require.Equal(t, "Ivan", delivery.Name)
}
func TestGetPayment(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()
	ctx := context.Background()

	paymentID := 1

	// db error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE id = $1")).
		WithArgs(paymentID).
		WillReturnError(errors.New("db error"))
	_, err := repo.GetPayment(ctx, paymentID)
	require.Error(t, err)

	// success
	rows := sqlmock.NewRows([]string{
		"transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee",
	}).AddRow("t", "", "", "p", 1, 0, "", 0, 0, 0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE id = $1")).
		WithArgs(paymentID).
		WillReturnRows(rows)
	payment, err := repo.GetPayment(ctx, paymentID)
	require.NoError(t, err)
	require.Equal(t, "t", payment.Transaction)
}
