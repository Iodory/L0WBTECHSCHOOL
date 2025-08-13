package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testex/internal/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func DBConn(dsn string) error {
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("Ошибка при ping %v", err)
	}

	fmt.Println("Подключение успешно")
	return nil
}

func LoadAllOrdersFromDB(db *sql.DB) ([]models.Order, error) {
	rows, err := db.Query("SELECT order_uid, track_number, entry FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry); err != nil {
			log.Printf("Ошибка сканирования: %v", err)
			continue
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
