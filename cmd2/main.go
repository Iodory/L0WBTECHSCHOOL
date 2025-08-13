package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testex/internal/models"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	if err := Producer("127.0.0.1:9092", "test-topic"); err != nil {
		log.Fatal(err)
	}
}

func Producer(brokerAddr, topic string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddr},
		Topic:   topic,
	})

	defer writer.Close()

	for i := 0; i < 10; i++ {
		order := models.Order{
			OrderUID:    fmt.Sprintf("order-%d", i),
			TrackNumber: fmt.Sprintf("track-%d", i),
			Entry:       "web",
			Delivery: models.Delivery{
				OrderUID: fmt.Sprintf("order-%d", i),
				Name:     "John Doe",
				Phone:    "+1234567890",
				Zip:      "12345",
				City:     "CityName",
				Address:  "Some street, 1",
				Region:   "RegionName",
				Email:    "john@example.com",
			},
			Payment: models.Payment{
				OrderUID:    fmt.Sprintf("order-%d", i),
				Transaction: "txn12345",
				Amount:      1000,
				Currency:    "USD",
			},
			Items: []models.Item{
				{
					OrderUID:    fmt.Sprintf("order-%d", i),
					ChrtID:      int64(i),
					TrackNumber: fmt.Sprintf("track-%d", i),
					Price:       1000,
					Name:        "Item Name",
					TotalPrice:  1000,
				},
			},
			DateCreated: time.Now(),
		}

		data, err := json.Marshal(order)
		if err != nil {
			return err
		}

		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: data,
			},
		)
		if err != nil {
			return err
		}

		log.Printf("Отправлен заказ %s", order.OrderUID)
		time.Sleep(1 * time.Second)
	}

	return nil
}
