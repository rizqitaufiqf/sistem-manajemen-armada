package main

import (
	"log"
	"sistem-manajemen-armada/internal/config"
	"sistem-manajemen-armada/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := config.LoadConfig()

	rabbitMQService, err := services.NewRabbitMQService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ service: %v", err)
	}
	defer rabbitMQService.Close()

	log.Println("RabbitMQ Worker started. Waiting for messages...")
	// Keep the worker running indefinitely
	select {}
}
