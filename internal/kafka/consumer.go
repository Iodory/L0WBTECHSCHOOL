package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testex/internal/cache"
	"testex/internal/db"
	"testex/internal/models"
	"testex/internal/service"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokerAddress []string, topic, groupID string, orderCashe *cache.OrderCache) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokerAddress,
		Topic:     topic,
		GroupID:   groupID,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					errCh <- nil
					return
				}

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

			if err := orderCashe.Set(order); err != nil {
				log.Print("Ошибка записи кеша", err)
				continue
			}
		}
	}()

	select {
	case <-sig:
		log.Print("Получен сигнал стоп")
		cancel()
	case err := <-errCh:
		if err != nil {
			log.Println("Ошибка в consumer", err)
		}
	}

	if err := r.Close(); err != nil {
		log.Println("Ошибка закрытия", err)
	}
	log.Println("Consumer остановлен")
}
