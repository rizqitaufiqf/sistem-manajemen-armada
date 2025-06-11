package config

import (
	"log"
	"os"
)

type AppConfig struct {
	// Database
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string

	// MQTT
	MQTTHost string

	// RabbitMQ
	RabbitMQHost string
	RabbitMQPort string
	RabbitMQUser string
	RabbitMQPass string

	// API
	APIPort string
}

func LoadConfig() *AppConfig {
	cfg := &AppConfig{
		PostgresHost: getEnv("POSTGRES_HOST", "db"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresUser: getEnv("POSTGRES_USER", "user"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "password"),
		PostgresDB:   getEnv("POSTGRES_DB", "vehicle_db"),
		MQTTHost:     getEnv("MQTT_BROKER", "tcp://mqtt_broker:1883"),
		RabbitMQHost: getEnv("RABBITMQ_HOST", "rabbitmq"),
		RabbitMQPort: getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUser: getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPass: getEnv("RABBITMQ_PASS", "guest"),
		APIPort:      getEnv("API_PORT", "8080"),
	}

	log.Println("Configuration loaded.")
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
