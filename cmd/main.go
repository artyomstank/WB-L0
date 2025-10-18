package main

import (
	"L0-wb/config"
	"L0-wb/internal/db"
	"L0-wb/internal/handler"
	"L0-wb/internal/kafka"
	"L0-wb/internal/repo"
	"L0-wb/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	sqlDB := db.NewDB(cfg)
	if sqlDB == nil {
		log.Fatal("failed to initialize database")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}()

	pgRepo := repo.NewRepo(sqlDB)
	svc, err := service.NewService(pgRepo)
	if err != nil {
		log.Fatalf("failed to initialize service: %v", err)
	}

	h := handler.NewHandler(svc)

	// Создаем HTTP сервер до запуска консьюмера
	srv := handler.NewServer(cfg, h)

	cons, err := kafka.NewConsumer(*cfg, svc)
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем консьюмер в горутине и обрабатываем ошибки
	go func() {
		if err := cons.ConsumeMessages(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("consumer error: %v", err)
		}
	}()

	// Запускаем HTTP сервер в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	}()

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown signal received")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := cons.Close(); err != nil {
		log.Printf("Error stopping consumer: %v", err)
	}

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	cancel() // Отменяем основной контекст
	log.Println("Shutdown complete")
}
