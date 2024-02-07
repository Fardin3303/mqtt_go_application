package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pg/pg/v10"
	"github.com/mqtt_go_application/pkg/models"
)

func InitMQTT(db *pg.DB) error {
	opts := mqtt.NewClientOptions().AddBroker("host.docker.internal:1883")
	opts.SetClientID("go-mqtt-example")

	client := mqtt.NewClient(opts)

	// Define a connection lost handler to handle connection losses
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
		// Attempt to reconnect here if necessary
	})

	// Attempt to connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	// Subscribe to the MQTT topic and handle incoming messages
	topic := "charger/1/connector/1/session/1"
	if token := client.Subscribe(topic, 0, handleMessage); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to MQTT topic %s: %w", topic, token.Error())
	}
	defer client.Unsubscribe(topic)

	// Start a goroutine to publish new sessions every 1 minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Check if the client is connected before publishing messages
				if !client.IsConnected() {
					log.Println("MQTT client is not connected")
					continue
				}

				// Create a new session message
				sessionMsg := models.MQTTMessage{
					SessionID:            generateSessionID(),
					EnergyDeliveredInKWh: 30,
					DurationInSeconds:    45,
					SessionCostInCents:   70,
					Timestamp:            time.Now(),
				}

				// Convert the message to JSON
				msgJSON, err := json.Marshal(sessionMsg)
				if err != nil {
					log.Printf("Error marshalling session message: %s", err.Error())
					continue
				}

				// Publish the message to the MQTT broker
				token := client.Publish(topic, 0, false, msgJSON)
				token.Wait()
				if token.Error() != nil {
					log.Printf("Error publishing session message: %s", token.Error())
					continue
				}

				log.Printf("Published new session message: %s", msgJSON)
			}
		}
	}()

	return nil
}

// Function to handle incoming MQTT messages
func handleMessage(client mqtt.Client, msg mqtt.Message) {
	
	log.Printf("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
}

// Function to generate a unique session ID
func generateSessionID() int {
	// Implement your session ID generation logic here
	// For simplicity, you can use a timestamp-based ID or a random number generator
	return int(time.Now().Unix())
}
