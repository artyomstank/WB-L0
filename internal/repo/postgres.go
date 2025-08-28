package repo

import (
	"L0-wb/internal/models"
	"context"
	"database/sql"
	"fmt"
)

func (pgs *PostgresRepo) CreateOrder(ctx context.Context, order models.Order) error {
	if order.OrderUID == "" {
		return fmt.Errorf(" order_uid cannot be empty")
	}

	// открываем транзакцию
	tx, err := pgs.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// откат при ошибке
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// zap.S().Errorf("rollback error: %v\n", rbErr) --добавить логи
			}
		}
	}()

	// создаём delivery запись (должна быть функция, которая делает INSERT и возвращает ID)
	deliveryID, err := pgs.CreateDeliveryTx(ctx, tx, order.Delivery)
	if err != nil {
		return fmt.Errorf("delivery creation error: %w", err)
	}

	// создаём payment запись
	paymentID, err := pgs.CreatePaymentTx(ctx, tx, order.Payment)
	if err != nil {
		return fmt.Errorf("payment creation error: %w", err)
	}

	// вставляем сам order
	query := `INSERT INTO orders (
		order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
	)`

	_, err = tx.ExecContext(ctx, query,
		order.OrderUID, order.TrackNumber, order.Entry, deliveryID, paymentID,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("order creation error: %w", err)
	}

	// вставляем все items
	for i, item := range order.Items {
		_, err := pgs.CreateItemTx(ctx, tx, item, order.OrderUID)
		if err != nil {
			return fmt.Errorf("item %d creation error: %w", i+1, err)
		}
	}

	// коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// create запросы для транзакции CreateOrder
func (pgs *PostgresRepo) CreateDeliveryTx(ctx context.Context, tx *sql.Tx, del models.Delivery) (int, error) {
	if del.Name == "" {
		return 0, fmt.Errorf(" recipient name is required")
	}
	if del.Phone == "" {
		return 0, fmt.Errorf(" recipient phone is required")
	}

	query := `INSERT INTO delivery (name, phone, zip, city, address, region, email) 
	          VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`
	var id int
	err := tx.QueryRowContext(ctx, query, del.Name, del.Phone, del.Zip, del.City, del.Address,
		del.Region, del.Email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

func (pgs *PostgresRepo) CreatePaymentTx(ctx context.Context, tx *sql.Tx, pay models.Payment) (int, error) {
	if pay.Transaction == "" {
		return 0, fmt.Errorf(" transaction number is required")
	}
	if pay.Provider == "" {
		return 0, fmt.Errorf(" payment provider is required")
	}
	if pay.Amount <= 0 {
		return 0, fmt.Errorf(" payment amount must be greater than zero")
	}
	query := `INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`

	var id int
	err := tx.QueryRowContext(ctx, query, pay.Transaction, pay.RequestID, pay.Currency, pay.Provider, pay.Amount,
		pay.PaymentDt, pay.Bank, pay.DeliveryCost, pay.GoodsTotal, pay.CustomFee).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (pgs *PostgresRepo) CreateItemTx(ctx context.Context, tx *sql.Tx, item models.Item, orderUID string) (int, error) {
	if item.Name == "" {
		return 0, fmt.Errorf(" item name is required")
	}
	if item.Price <= 0 {
		return 0, fmt.Errorf(" item price must be greater than zero")
	}
	if orderUID == "" {
		return 0, fmt.Errorf(" item order_uid is required")
	}

	query := `INSERT INTO item (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid) 
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`
	var id int
	err := tx.QueryRowContext(ctx, query, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
		item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status, orderUID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (pgs *PostgresRepo) GetOrder(ctx context.Context, orderUID string) (models.Order, error) {
	var order models.Order
	var deliveryID, paymentID int

	if orderUID == "" {
		return order, fmt.Errorf("order_uid cannot be empty")
	}

	query := `SELECT order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1`
	row := pgs.DB.QueryRowContext(ctx, query, orderUID)
	err := row.Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &deliveryID, &paymentID,
		&order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		return order, fmt.Errorf("order retrieval error: %w", err)
	}

	order.Delivery, err = pgs.GetDelivery(ctx, deliveryID)
	if err != nil {
		return order, fmt.Errorf("delivery data retrieval error: %w", err)
	}

	order.Payment, err = pgs.GetPayment(ctx, paymentID)
	if err != nil {
		return order, fmt.Errorf("payment data retrieval error: %w", err)
	}

	order.Items, err = pgs.GetItemsByOrderUID(ctx, orderUID)
	if err != nil {
		return order, fmt.Errorf("order items retrieval error: %w", err)
	}

	return order, nil
}

func (pgs *PostgresRepo) GetDelivery(ctx context.Context, deliveryID int) (models.Delivery, error) {
	var del models.Delivery

	query := `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE id = $1`
	row := pgs.DB.QueryRowContext(ctx, query, deliveryID)
	err := row.Scan(&del.Name, &del.Phone, &del.Zip, &del.City, &del.Address, &del.Region, &del.Email)
	if err != nil {
		return del, err
	}
	return del, nil
}
func (pgs *PostgresRepo) GetPayment(ctx context.Context, paymentID int) (models.Payment, error) {
	var pay models.Payment

	query := `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE id = $1`
	row := pgs.DB.QueryRowContext(ctx, query, paymentID)
	err := row.Scan(&pay.Transaction, &pay.RequestID, &pay.Currency, &pay.Provider, &pay.Amount,
		&pay.PaymentDt, &pay.Bank, &pay.DeliveryCost, &pay.GoodsTotal, &pay.CustomFee)
	if err != nil {
		return pay, err
	}
	return pay, nil
}
func (pgs *PostgresRepo) GetItemsByOrderUID(ctx context.Context, orderUID string) (models.Items, error) {
	var items models.Items

	query := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM item WHERE order_uid = $1`
	rows, err := pgs.DB.QueryContext(ctx, query, orderUID)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return items, err
	}
	return items, nil
}
func (pgs *PostgresRepo) GetLastOrders(ctx context.Context, lim int) ([]models.Order, error) {
	query := `SELECT order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders`
	rows, err := pgs.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("order query error: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var deliveryID, paymentID int
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &deliveryID, &paymentID, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, fmt.Errorf("order scanning error: %w", err)
		}

		order.Delivery, err = pgs.GetDelivery(ctx, deliveryID)
		if err != nil {
			return nil, fmt.Errorf("delivery data retrieval error for order %s: %w", order.OrderUID, err)
		}

		order.Payment, err = pgs.GetPayment(ctx, paymentID)
		if err != nil {
			return nil, fmt.Errorf("payment data retrieval error for order %s: %w", order.OrderUID, err)
		}

		order.Items, err = pgs.GetItemsByOrderUID(ctx, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("items retrieval error for order %s %w: ", order.OrderUID, err)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("order iteration error")
	}

	return orders, nil
}
