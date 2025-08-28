package main

import (
	"L0-wb/config"
	"L0-wb/internal/kafka"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.LoadConfig()
	log.Println("config initialized")

	prdcr := kafka.NewProducer(*cfg)

	defer prdcr.Close()
	log.Println("producer initialized")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем продюсера в отдельной горутине
	go func() {
		err := prdcr.RunProducer(ctx)
		if err != nil {
			log.Printf("producer stopped: %v", err)
		}
	}()

	// Канал чтения сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ждём  команду сигнала или пока контекст не умрет
	select {
	case <-quit:
	case <-ctx.Done():
	}

	log.Println("shutting down gracefully...")
}
