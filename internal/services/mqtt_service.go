package services

import (
	"log"
	"sistem-manajemen-armada/internal/config"
	"sistem-manajemen-armada/internal/repository"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTService struct {
	client          mqtt.Client
	repo            *repository.PostgreSQLRepository
	rabbitmqService *RabbitMQService
}

func NewMQTTService(cfg *config.AppConfig, repo *repository.PostgreSQLRepository, rabbitMQService *RabbitMQService) *MQTTService {
	options := mqtt.NewClientOptions().AddBroker(cfg.MQTTHost).SetClientID("backend_subscriber")
	options.SetKeepAlive(60 * time.Second)
	options.SetPingTimeout(10 * time.Second)
	options.SetCleanSession(true)
	options.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Connected to MQTT broker!")
	})
	options.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	})

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	return &MQTTService{client: client, repo: repo, rabbitmqService: rabbitMQService}
}

func (s *MQTTService) Disconnect() {
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect(250)
		log.Println("MQTT client disconnected.")
	}
}
