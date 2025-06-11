package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sistem-manajemen-armada/internal/config"
	"sistem-manajemen-armada/internal/models"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQService(cfg *config.AppConfig) (*RabbitMQService, error) {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitMQUser, cfg.RabbitMQPass, cfg.RabbitMQHost, cfg.RabbitMQPort)
	var conn *amqp.Connection
	var err error

	for i := range 10 {
		conn, err = amqp.Dial(connString)
		if err == nil {
			log.Println("Successfully connected to RabbitMQ!")
			break
		}
		log.Printf("Waiting for RabbitMQ to be ready... (%d/10) %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after multiple attempts: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare exchange and queue
	err = channel.ExchangeDeclare(
		"armada.events", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange %w", err)
	}

	_, err = channel.QueueDeclare(
		"geofence_alerts", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = channel.QueueBind(
		"geofence_alerts", // queue name
		"geofence_entry",  // routing key (could be specific event types)
		"armada.events",   // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQService{conn: conn, channel: channel}, nil
}

func (s *RabbitMQService) PublishGeofenceEvent(event models.GeofenceEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal geofence event: %w", err)
	}

	err = s.channel.Publish(
		"armada.events",  // exchange
		"geofence_entry", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	log.Printf("Published geofence event for vehicle %s", event.VehicleID)
	return nil
}

func (s *RabbitMQService) StartGeofenceWorker() {
	messages, err := s.channel.Consume(
		"geofence_alerts", // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for d := range messages {
			log.Printf("Received a geofence alert message: %s", d.Body)
		}
	}()
}

func (s *RabbitMQService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	log.Println("RabbitMQ connection closed.")
}
