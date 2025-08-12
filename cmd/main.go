package main

import (
	"fmt"
	"log"
	"net/http"
	"testex/internal/cache"
	"testex/internal/config"
	"testex/internal/db"
	"testex/internal/hanlders"
	"testex/internal/kafka"
	"testex/internal/migrations"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")

	if err != nil {
		log.Fatal("Ошибка загрузки конфига", err)
	}

	order, err := cache.NewOrderCache(5 * time.Minute)
	if err != nil {
		log.Fatal("Ошибка создания кеша", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	if err := db.DBConn(dsn); err != nil {
		log.Fatal(err)
	}

	broker := cfg.Kafka.Brokers
	topic := cfg.Kafka.Topic
	groupID := cfg.Kafka.GroupID

	migrations.RunMigration(dsn)

	go kafka.StartConsumer(broker, topic, groupID, order)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../index.html")
	})
	http.HandleFunc("/order", hanlders.GetOrder(order))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
