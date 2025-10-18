package repo

import (
	"L0-wb/internal/models"
	"context"
	"database/sql"
)

type Repository interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrder(ctx context.Context, orderUID string) (models.Order, error)
	GetLastOrders(ctx context.Context, lim int) ([]models.Order, error)
	CreateDeliveryTx(ctx context.Context, tx *sql.Tx, del models.Delivery) (int, error)
	CreatePaymentTx(ctx context.Context, tx *sql.Tx, pay models.Payment) (int, error)
	CreateItemTx(ctx context.Context, tx *sql.Tx, item models.Item, orderUID string) (int, error)
	GetDelivery(ctx context.Context, deliveryID int) (models.Delivery, error)
	GetPayment(ctx context.Context, paymentID int) (models.Payment, error)
	GetItemsByOrderUID(ctx context.Context, orderUID string) (models.Items, error)
	Close() error
}
