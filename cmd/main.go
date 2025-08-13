package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testex/internal/cache"
	"testex/internal/config"
	"testex/internal/db"
	"testex/internal/hanlders"
	kafffka "testex/internal/kafka"
	"testex/internal/migrations"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal("Ошибка загрузки конфига", err)
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

	orderCache := cache.NewOrderCache(16, 5*time.Minute)

	migrations.RunMigration(dsn)

	if err := cache.WarmUpCache(db.DB, orderCache); err != nil {
		log.Fatal("Ошибка прогрева", err)
	}

	broker := cfg.Kafka.Brokers
	topic := cfg.Kafka.Topic
	groupID := cfg.Kafka.GroupID

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafffka.StartConsumer(ctx, broker, topic, groupID, orderCache)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("/order", hanlders.GetOrder(orderCache))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP сервер упал: %v", err)
		}
	}()

	stroper := make(chan os.Signal, 1)
	signal.Notify(stroper, os.Interrupt, syscall.SIGINT)

	<-stroper
	log.Println("Сигнал стоп")
	cancel()

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()

	if err := srv.Shutdown(shutCtx); err != nil {
		log.Fatal(err)
	}

	log.Println("Приложение остановлено")
}
