package main

import (
	"log"
	"sistem-manajemen-armada/internal/config"
	"sistem-manajemen-armada/internal/repository"
	"sistem-manajemen-armada/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := config.LoadConfig()

	// Initialize Database Repository
	repo, err := repository.NewPostgreSQLRepository(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL repository: %v", err)
	}
	defer repo.DB.Close()

	//Init RabbitMQ Service
	rabbitMQService, err := services.NewRabbitMQService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ service: %v", err)
	}
	defer rabbitMQService.Close()

	// Init MQTT Service
	mqttService := services.NewMQTTService(cfg, repo, rabbitMQService)
	defer mqttService.Disconnect()

	log.Println("MQTT Subscriber started. Waiting for messages...")
	// Keep the subscriber running indefinitely
	select {}
}
