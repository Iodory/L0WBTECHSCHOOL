package kafka

import (
	"context"
	"encoding/json"
	"log"
	"testex/internal/cache"
	"testex/internal/db"
	"testex/internal/models"
	"testex/internal/service"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokerAddress, topic, groupID string, orderCashe *cache.OrderCache) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokerAddress},
		Topic:     topic,
		GroupID:   groupID,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	ctx := context.Background()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Ошибка чтения сообщения %v", err)
			time.Sleep(time.Second)
			continue
		}

		order := models.Order{}

		log.Printf("получено сообщение %s", string(m.Value))

		err = json.Unmarshal(m.Value, &order)
		if err != nil {
			log.Printf("Ошибка парсинга %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Print("Успешный парсинг")

		if err := service.InsertOrder(db.DB, order); err != nil {
			log.Printf("Ошибка вставки в БД: %v", err)
			continue
		}

		orderCashe.Set(order)
	}
}
