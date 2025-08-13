package kafffka

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

func StartConsumer(ctx context.Context, brokerAddress []string, topic, groupID string, orderCache *cache.OrderCache) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokerAddress,
		Topic:     topic,
		GroupID:   groupID,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	defer func() {
		if err := r.Close(); err != nil {
			log.Println("Ошибка закрытия consumer:", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Контекст отменён, consumer останавливается")
			return
		default:
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Ошибка чтения сообщения: %v", err)
				time.Sleep(time.Second)
				continue
			}

			order := models.Order{}
			log.Printf("получено сообщение %s", string(m.Value))

			err = json.Unmarshal(m.Value, &order)
			if err != nil {
				log.Printf("Ошибка парсинга: %v", err)
				continue
			}
			if err := order.Validate(); err != nil {
				log.Printf("Ошибка валидации: %v", err)
				continue
			}

			if err := service.InsertOrder(db.DB, order); err != nil {
				log.Printf("Ошибка вставки в БД: %v", err)
				continue
			}

			if err := orderCache.Set(order); err != nil {
				log.Printf("Ошибка записи кеша: %v", err)
				continue
			}
		}
	}
}
