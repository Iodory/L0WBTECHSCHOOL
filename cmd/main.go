package main

import (
	"log"
	"net/http"
	"testex/internal/cache"
	"testex/internal/db"
	"testex/internal/hanlders"
	"testex/internal/kafka"
)

func main() {
	order := cache.NewOrderCache()

	dsn := "postgres://postgres:G0Lang_DB!7x@9zQw@localhost:5432/LOWB?sslmode=disable"
	if err := db.DBConn(dsn); err != nil {
		log.Fatal(err)
	}

	broker := "localhost:9092"
	topic := "test-topic"
	groupID := "test-group"
	go kafka.StartConsumer(broker, topic, groupID, order)

	http.HandleFunc("/order", hanlders.GetOrder(order))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
