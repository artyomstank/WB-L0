package main

import (
	"L0-wb/config"
	"L0-wb/internal/db"
	"L0-wb/internal/handler"
	"L0-wb/internal/kafka"
	"L0-wb/internal/repo"
	"L0-wb/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	sqlDB := db.NewDB(cfg) // создаём *sql.DB
	defer sqlDB.Close()

	pgRepo := repo.NewRepo(sqlDB)
	svc, err := service.NewService(pgRepo)
	if err != nil {
		log.Fatalf("service init error: %v", err)
	}
	h := handler.NewHandler(svc)

	cons, err := kafka.NewConsumer(*cfg, svc)
	if err != nil {
		log.Fatalf("kafka.NewConsumer error: %v", err)
	}

	go cons.ConsumeMessages(context.Background())

	srv := handler.NewServer(cfg, h)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != context.Canceled {
			log.Printf("server run error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown signal received")

	if err := cons.Close(); err != nil {
		log.Printf("Error stopping consumer: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	log.Println("Shutdown complete")
}
