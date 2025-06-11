package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sistem-manajemen-armada/internal/config"
	"sistem-manajemen-armada/internal/models"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := config.LoadConfig()

	vehicleID := os.Getenv("VEHICLE_ID")
	if vehicleID == "" {
		vehicleID = "B1234XYZ" // Default vehicle ID
	}

	initialLat, err := strconv.ParseFloat(os.Getenv("INITIAL_LAT"), 64)
	if err != nil {
		initialLat = -6.2520
	}
	initialLon, err := strconv.ParseFloat(os.Getenv("INITIAL_LON"), 64)
	if err != nil {
		initialLon = 106.8730
	}

	// Init options MQTT connection
	options := mqtt.NewClientOptions().AddBroker(cfg.MQTTHost).SetClientID("mqtt_publisher_" + vehicleID)
	options.SetKeepAlive(60 * time.Second)
	options.SetPingTimeout(10 * time.Second)
	options.SetCleanSession(true)
	options.SetOnConnectHandler(func(client mqtt.Client) {
		log.Printf("MQTT Publisher for %s connected to broker!", vehicleID)
	})
	options.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("MQTT Publisher for %s connection lost: %v", vehicleID, err)
	})

	// MQTT connection
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect MQTT Publisher for %s: %v", vehicleID, token.Error())
	}
	defer client.Disconnect(250)

	log.Printf("Starting mock location data publishing for vehicle %s...", vehicleID)

	currentLat := initialLat
	currentLon := initialLon

	directionLat := 0.00001
	directionLon := 0.00001

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentLat -= directionLat
		currentLon += directionLon

		payload := models.VehicleLocation{
			VehicleID: vehicleID,
			Latitude:  currentLat,
			Longitude: currentLon,
			Timestamp: time.Now().Unix(),
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			continue
		}

		topic := fmt.Sprintf("/armada/vehicle/%s/location", vehicleID)
		token := client.Publish(topic, 1, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Error publishing message to %s: %v", topic, token.Error())
		} else {
			log.Printf("Published to %s: %s", topic, string(jsonData))
		}
	}

}
