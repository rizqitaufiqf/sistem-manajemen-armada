package services

import (
	"encoding/json"
	"log"
	"sistem-manajemen-armada/internal/config"
	"sistem-manajemen-armada/internal/models"
	"sistem-manajemen-armada/internal/repository"
	"sistem-manajemen-armada/internal/utils"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTService struct {
	client          mqtt.Client
	repo            *repository.PostgreSQLRepository
	rabbitmqService *RabbitMQService
	cfg             *config.AppConfig
}

func NewMQTTService(cfg *config.AppConfig, repo *repository.PostgreSQLRepository, rabbitMQService *RabbitMQService) *MQTTService {
	options := mqtt.NewClientOptions().AddBroker(cfg.MQTTHost).SetClientID("backend_subscriber")
	options.SetKeepAlive(60 * time.Second)
	options.SetPingTimeout(10 * time.Second)
	options.SetCleanSession(true)
	options.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Successfully Connected to MQTT broker!")
	})
	options.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	})

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	return &MQTTService{client: client, repo: repo, rabbitmqService: rabbitMQService, cfg: cfg}
}

func (s *MQTTService) Disconnect() {
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect(250)
		log.Println("MQTT client disconnected.")
	}
}

func (s *MQTTService) SubscribeToVehicleLocationTopic() {
	topic := "/fleet/vehicle/+/location"
	token := s.client.Subscribe(topic, 1, s.mqttHandler)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}
	log.Printf("Subscribed to MQTT topic: %s", topic)
}

func (s *MQTTService) mqttHandler(client mqtt.Client, message mqtt.Message) {
	log.Printf("Received MQTT message on topic: %s", message.Topic())

	var vehicleLocation models.VehicleLocation
	if err := json.Unmarshal(message.Payload(), &vehicleLocation); err != nil {
		log.Printf("Error unmarshalling MQTT payload: %v", err)
		return
	}

	//? delete this
	// log.Printf("MQTT Data => Vehicle ID: %s, Lat: %.6f, Lon: %.6f, Timestamp: %d",
	// 	vehicleLocation.VehicleID, vehicleLocation.Latitude, vehicleLocation.Longitude, vehicleLocation.Timestamp)

	validationErrors := vehicleLocation.ValidateVehicleLocationData()
	if len(validationErrors) > 0 {
		log.Printf("Error saving location to DB: %s", strings.Join(validationErrors, ", "))
		return
	}

	// store data
	if err := s.repo.InsertVehicleLocation(vehicleLocation); err != nil {
		log.Printf("Error saving location to DB: %v", err)
	} else {
		log.Printf("Location for vehicle %s saved to DB.", vehicleLocation.VehicleID)
	}

	s.geofenceCheck(vehicleLocation)
}

func (s *MQTTService) geofenceCheck(vehicleLocation models.VehicleLocation) {
	// Kantor Pusat TJ Latitude : -6.2526, Longitude: 106.8736
	geofenceCenterLat := s.cfg.GeofenceLat
	geofenceCenterLon := s.cfg.GeofenceLon
	radiusMeters := s.cfg.GeofenceRadius

	if utils.IsInGeofence(vehicleLocation.Latitude, vehicleLocation.Longitude, geofenceCenterLat, geofenceCenterLon, radiusMeters) {
		log.Printf("Vehicle %s entered geofence!", vehicleLocation.VehicleID)

		eventPayload :=
			models.GeofenceEvent{
				VehicleID: vehicleLocation.VehicleID,
				Event:     "geofence_entry",
				Location: models.Location{
					Latitude:  vehicleLocation.Latitude,
					Longitude: vehicleLocation.Longitude,
				},
				Timestamp: vehicleLocation.Timestamp,
			}

		if s.rabbitmqService != nil {
			if err := s.rabbitmqService.PublishGeofenceEvent(eventPayload); err != nil {
				log.Printf("Error publishing geofence event to RabbitMQ: %v", err)
			}
		} else {
			log.Println("RabbitMQService not initialized, cannot publish geofence event.")
		}
	}
}
