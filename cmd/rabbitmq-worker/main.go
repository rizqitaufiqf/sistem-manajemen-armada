package main

import (
	"log"
	"sistem-manajemen-armada/internal/config"
	"sistem-manajemen-armada/internal/services"
)

func main() {
	// enable this if you want use from local
	// godotenv.Load()
	cfg := config.LoadConfig()

	rabbitMQService, err := services.NewRabbitMQService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ service: %v", err)
	}
	defer rabbitMQService.Close()

	// Start Geofence Worker
	rabbitMQService.StartGeofenceWorker()

	log.Println("RabbitMQ Worker started. Waiting for messages...")
	// Keep the worker running indefinitely
	select {}
}
