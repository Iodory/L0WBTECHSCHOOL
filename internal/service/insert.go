package service

import (
	"database/sql"
	"fmt"
	"log"
	"testex/internal/models"
)

func InsertOrder(db *sql.DB, order models.Order) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.Exec(`
        INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        INSERT INTO delivery (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        INSERT INTO payment (
            transaction, order_uid, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.Payment.Transaction, order.OrderUID, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
            INSERT INTO items (
                order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
            ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	log.Printf("Заказ %s успешно вставлен", order.OrderUID)
	return nil
}

func LoadAllOrders(db *sql.DB) ([]models.Order, error) {
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			return nil, err
		}

		order.Delivery, _ = getDeliveryByOrderUID(db, order.OrderUID)
		order.Payment, _ = getPaymentByOrderUID(db, order.OrderUID)
		order.Items, _ = getItemsByOrderUID(db, order.OrderUID)

		orders = append(orders, order)
	}

	return orders, nil
}

func getDeliveryByOrderUID(db *sql.DB, orderUID string) (models.Delivery, error) {
	var d models.Delivery
	err := db.QueryRow("SELECT name, phone, zip, city, address, region, email FROM deliveries WHERE order_uid=$1", orderUID).
		Scan(&d.Name, &d.Phone, &d.Zip, &d.City, &d.Address, &d.Region, &d.Email)
	d.OrderUID = orderUID
	return d, err
}

func getPaymentByOrderUID(db *sql.DB, orderUID string) (models.Payment, error) {
	var p models.Payment
	err := db.QueryRow(`
        SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payments WHERE order_uid=$1`, orderUID).
		Scan(&p.Transaction, &p.RequestID, &p.Currency, &p.Provider, &p.Amount, &p.PaymentDt,
			&p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee)
	p.OrderUID = orderUID
	return p, err
}

func getItemsByOrderUID(db *sql.DB, orderUID string) ([]models.Item, error) {
	rows, err := db.Query(`
        SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
        FROM items WHERE order_uid=$1`, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var i models.Item
		err := rows.Scan(&i.ChrtID, &i.TrackNumber, &i.Price, &i.Rid, &i.Name, &i.Sale, &i.Size, &i.TotalPrice, &i.NmID, &i.Brand, &i.Status)
		if err != nil {
			return nil, err
		}
		i.OrderUID = orderUID
		items = append(items, i)
	}

	return items, nil
}

func GetOrderByID(db *sql.DB, orderUID string) (models.Order, error) {
	var order models.Order

	queryOrder := `
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, 
		       delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`
	err := db.QueryRow(queryOrder, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, fmt.Errorf("order not found")
		}
		return order, err
	}

	queryDelivery := `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries
		WHERE order_uid = $1
	`
	err = db.QueryRow(queryDelivery, orderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email)
	if err != nil {
		return order, err
	}

	queryPayment := `
		SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments
		WHERE order_uid = $1
	`
	err = db.QueryRow(queryPayment, orderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee)
	if err != nil {
		return order, err
	}

	queryItems := `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`
	rows, err := db.Query(queryItems, orderUID)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return order, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}
