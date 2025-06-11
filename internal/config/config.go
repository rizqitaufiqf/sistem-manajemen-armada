package config

import (
	"log"
	"os"
	"strconv"
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

	GeofenceLat    float64
	GeofenceLon    float64
	GeofenceRadius float64
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

	geofenceLatStr := getEnv("GEOFENCE_LAT", "-6.2526")
	geofenceLonStr := getEnv("GEOFENCE_LON", "106.8736")
	geofenceRadiusStr := getEnv("GEOFENCE_RADIUS", "50.0")

	lat, err := strconv.ParseFloat(geofenceLatStr, 64)
	if err != nil {
		log.Fatalf("Error parsing GEOFENCE_LAT '%s' to float64: %v", geofenceLatStr, err)
	}
	cfg.GeofenceLat = lat

	lon, err := strconv.ParseFloat(geofenceLonStr, 64)
	if err != nil {
		log.Fatalf("Error parsing GEOFENCE_LON '%s' to float64: %v", geofenceLonStr, err)
	}
	cfg.GeofenceLon = lon

	radius, err := strconv.ParseFloat(geofenceRadiusStr, 64)
	if err != nil {
		log.Fatalf("Error parsing GEOFENCE_RADIUS '%s' to float64: %v", geofenceRadiusStr, err)
	}
	cfg.GeofenceRadius = radius

	log.Println("Configuration loaded.")
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
